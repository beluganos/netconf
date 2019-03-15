// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cfglxclib

import (
	"flag"
	"fmt"
	"io/ioutil"
	lxdlib "netconf/lib/lxd"
	"os"
	"path"
	"strings"

	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/api"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	PROFILE_BACKUP_DIR = "/tmp"
	PROFILE_IMAGE_NANE = "base"
	CONTAINER_LOG_DIR  = "/var/log/beluganos"
)

func Usage(code int) {
	flag.PrintDefaults()
	os.Exit(code)
}

type Args struct {
	Oper      Operation
	Name      string
	IFNames   []string
	BackupDir string
	IFPrefix  string
	Image     string
	Verbose   bool
}

func (a *Args) Backup() string {
	return path.Join(a.BackupDir, fmt.Sprintf("backup_lxd_profile_%s.yml", a.Name))
}

func (a *Args) Parse() {
	var opstr string
	flag.StringVar(&opstr, "cmd", "help", fmt.Sprintf("command name [%s]", StrOperations()))
	flag.StringVar(&a.Name, "name", "", "container name")
	flag.StringVar(&a.BackupDir, "backup", PROFILE_BACKUP_DIR, "backup directory")
	flag.StringVar(&a.Image, "image", PROFILE_IMAGE_NANE, "Container Image name")
	flag.StringVar(&a.IFPrefix, "prefix", "v", "Prefix of parent ifname")
	flag.BoolVar(&a.Verbose, "verbose", false, "show detail messages")
	flag.Parse()

	oper, err := ParseOperation(opstr)
	if err != nil {
		Usage(1)
	}

	if len(a.Name) == 0 {
		Usage(1)
	}

	a.IFNames = flag.Args()
	a.Oper = oper
}

func LoadProfile(client *lxdlib.Client, name, path string) error {
	log.Debugf("LoadProfile: name='%s' path='%s'", name, path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("LoadProfile: ReadFile error. %s", err)
		return err
	}

	profile := &api.ProfilePut{}
	if err := yaml.Unmarshal(data, profile); err != nil {
		log.Errorf("LoadProfile: Unmatshal error. %s", err)
		return err
	}

	if err := client.CreateOrUpdateProfile(name, profile); err != nil {
		log.Errorf("LoadProfile: CreateOrUpdateProfile error, %s", err)
		return err
	}

	log.Debugf("LoadProfile: ok. name='%s' path='%s'", name, path)
	return nil
}

func SaveProfile(client *lxdlib.Client, name, path string) error {
	log.Debugf("SaveProfile: name='%s' path='%s'", name, path)

	profile, _, err := client.GetProfile(name)
	if err != nil {
		log.Errorf("SaveProfile: GetProfile error. %s", err)
		return err
	}

	data, err := yaml.Marshal(profile.Writable())
	if err != nil {
		log.Errorf("SaveProfile: Marshal error. %s", err)
		return err
	}

	if err := ioutil.WriteFile(path, data, os.FileMode(0644)); err != nil {
		log.Errorf("SaveProfile: WriteFile error. %s", err)
		return err
	}

	log.Debugf("SaveProfile: ok. name='%s' path='%s'", name, path)
	return nil
}

func ClearBackup(path string) error {
	log.Debugf("ClearBackup: path='%s'", path)

	if err := os.Remove(path); err != nil {
		log.Errorf("ClearBackup: Remove error. %s", err)
		return err
	}

	log.Debugf("ClearBackup: ok. '%s'", path)
	return nil
}

func CreateContainer(client *lxdlib.Client, name string, keep bool, logdir, mngIf, bridgeIf string) error {
	log.Debugf("CreateContainer: name='%s'", name)

	if _, _, err := client.InitializeProfile(name, logdir, mngIf, bridgeIf); err != nil {
		log.Errorf("CreateContainer: InitializeProfile error. %s", err)
		return err
	}

	if _, _, err := client.CreateContainer(name); err != nil {
		log.Errorf("CreateContainer: CreateContainer error. %s", err)
		return err
	}

	if err := client.UpdateContainerState(name, string(shared.Start)); err != nil {
		log.Errorf("CreateContainer: UpdateContainerState error. %s", err)
		if !keep {
			DeleteContainer(client, name)
		}
		return err
	}

	log.Debugf("CreateContainer: ok. name='%s'", name)
	return nil
}

func DeleteContainer(client *lxdlib.Client, name string) {
	log.Debugf("DelteContainer: name='%s'", name)

	if err := client.UpdateContainerState(name, string(shared.Stop)); err != nil {
		log.Errorf("DeleteContainer: UpdateContainerState error. %s", err)
	}

	if err := client.DeleteContainer(name); err != nil {
		log.Warnf("DelteContainer: DeleteContainer error. %s", err)
	}

	if err := client.DeleteProfile(name); err != nil {
		log.Warnf("DelteContainer: DeleteProfile error. %s", err)
	}

	log.Debugf("DelteContainer: ok. name='%s'", name)
}

func AddInterface(client *lxdlib.Client, name string, prefix string, ifname string, hwaddr string, mtu uint16) error {
	log.Debugf("AddInterface: name='%s' iface='%s' mac='%s' mru=%d", name, ifname, hwaddr, mtu)

	profile, _, err := client.GetProfile(name)
	if err != nil {
		log.Errorf("AddInterface: GetProfile errror. %s", err)
		return err
	}

	hostname := fmt.Sprintf("%s%s", prefix, ifname)
	if err := lxdlib.AddProfileDeviceNIC(profile, ifname, hostname, hwaddr, mtu); err != nil {
		log.Warnf("AddInterface: AddProfileDeviceNIC error. %s", err)
	}

	if err := client.UpdateProfile(name, profile.Writable()); err != nil {
		log.Errorf("AddInterface: UpdateProfile errror. %s", err)
		DeleteInterface(client, name, ifname)
		return err
	}

	log.Debugf("AddInterface: ok name='%s' iface='%s'", name, ifname)
	return nil
}

func ParseInterfaceArgs(args []string) map[string]string {
	m := make(map[string]string)

	for _, arg := range args {
		items := strings.SplitN(arg, "=", 2)
		k, v := func() (string, string) {
			switch len(items) {
			case 1:
				return items[0], ""
			case 2:
				return items[0], items[1]
			default:
				return "", ""
			}
		}()
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if len(k) > 0 {
			m[k] = v
		}
	}

	return m
}

func SetInterface(client *lxdlib.Client, name string, ifname string, negate bool, args ...string) error {
	log.Debugf("SetInterface: name='%s' iface='%s' negate=%t %s", name, ifname, negate, args)

	profile, _, err := client.GetProfile(name)
	if err != nil {
		log.Errorf("SetInterface: GetProfile errror. %s", err)
		return err
	}

	for key, val := range ParseInterfaceArgs(args) {
		if err := lxdlib.SetProfileDeviceNIC(profile, ifname, key, val, negate); err != nil {
			return err
		}
	}

	if err := client.UpdateProfile(name, profile.Writable()); err != nil {
		log.Errorf("DeleteInterface: UpdateProfile errror. %s", err)
		return err
	}

	log.Debugf("SetInterface: ok name='%s' iface='%s'", name, ifname)
	return nil
}

func DeleteInterface(client *lxdlib.Client, name string, ifname string) error {
	log.Debugf("DeleteInterface: name='%s' iface='%s'", name, ifname)

	profile, _, err := client.GetProfile(name)
	if err != nil {
		log.Errorf("DeleteInterface: GetProfile errror. %s", err)
		return err
	}

	if err := lxdlib.DelProfileDeviceNIC(profile, ifname); err != nil {
		log.Warnf("DeleteInterface: DeleteProfileDeviceNIC error. %s", err)
		return nil
	}

	if err := client.UpdateProfile(name, profile.Writable()); err != nil {
		log.Errorf("DeleteInterface: UpdateProfile errror. %s", err)
		return err
	}

	log.Debugf("DeleteInterface: ok name='%s' iface='%s'", name, ifname)
	return nil
}

func logDirPath(name string) string {
	return fmt.Sprintf("%s/%s", CONTAINER_LOG_DIR, name)
}

func MakeLogDir(name string) string {
	path := logDirPath(name)
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Warnf("MakeLogDir: mkdir error. %s", err)
		return ""
	}

	log.Debugf("MakeLogDir: ok. %s", path)
	return path
}

func RmLogDir(name string) {
	path := logDirPath(name)
	if err := os.RemoveAll(path); err != nil {
		log.Warnf("MakeLogDir: rmdir error. %s", err)
	} else {
		log.Debugf("MakeLogDir: rmdir ok. %s", path)
	}
}
