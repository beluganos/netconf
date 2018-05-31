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

package lxdlib

import (
	"fmt"

	lxd "github.com/lxc/lxd/client"
	api "github.com/lxc/lxd/shared/api"
)

const (
	DEFAULT_BRIDGE_NAME = "lxdbr0"
	DEFAULT_DEVICE_MTU  = "1500"
)

func Connect() (lxd.ContainerServer, error) {
	return lxd.ConnectLXDUnix("", nil)
}

func GetImage(c lxd.ContainerServer, name string) (*api.Image, string, error) {
	alias, _, err := c.GetImageAlias(name)
	if err != nil {
		return nil, "", err
	}

	return c.GetImage(alias.Target)
}

func GetState(c lxd.ContainerServer, name string) (*api.ContainerState, string, error) {
	return c.GetContainerState(name)
}

func CreateContainer(c lxd.ContainerServer, image string, container *api.ContainersPost) (lxd.RemoteOperation, error) {
	im, _, err := GetImage(c, image)
	if err != nil {
		return nil, err
	}

	return c.CreateContainerFromImage(c, *im, *container)
}

func NewDefaultContainer(name string) *api.ContainersPost {
	c := &api.ContainersPost{}
	c.Name = name
	c.Profiles = []string{name}
	return c
}

func NewProfileDeviceNIC(ifname, hostname, hwaddr string) map[string]string {
	return map[string]string{
		"type":      "nic",
		"name":      ifname,
		"host_name": hostname,
		"nictype":   "p2p",
		"mtu":       DEFAULT_DEVICE_MTU,
		"hwaddr":    hwaddr,
	}
}

func AddProfileDeviceNIC(p *api.Profile, ifname, hostname, hwaddr string) error {
	devices := p.Devices
	if _, ok := devices[ifname]; ok {
		return fmt.Errorf("Device already exists. %s", ifname)
	}

	devices[ifname] = NewProfileDeviceNIC(ifname, hostname, hwaddr)

	return nil
}

func DelProfileDeviceNIC(p *api.Profile, ifname string) error {
	devices := p.Devices
	if _, ok := devices[ifname]; !ok {
		return fmt.Errorf("Device not found. %s", ifname)
	}

	delete(devices, ifname)

	return nil
}

func SetProfileDeviceNIC(p *api.Profile, ifname string, key string, value string, negate bool) error {
	devices := p.Devices
	if _, ok := devices[ifname]; !ok {
		return fmt.Errorf("Device not found. %s", ifname)
	}

	if negate {
		delete(devices[ifname], key)
	} else {
		devices[ifname][key] = value
	}

	return nil
}

func NewDefaultProfile() *api.ProfilePut {
	return &api.ProfilePut{
		Config: map[string]string{
			"security.privileged": "true",
		},
		Devices: map[string]map[string]string{
			"eth0": {
				"name":    "eth0",
				"nictype": "bridged",
				"parent":  DEFAULT_BRIDGE_NAME,
				"type":    "nic",
			},
			"root": {
				"path": "/",
				"pool": "default",
				"type": "disk",
			},
		},
	}
}
