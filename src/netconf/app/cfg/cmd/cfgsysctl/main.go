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

package main

import (
	"flag"
	"fmt"
	prop "netconf/lib/property"
	"netconf/lib/sysctl"
	"os"

	log "github.com/sirupsen/logrus"
)

const SYSCTL_CONF_PATH = "/etc/sysctl.d/30-beluganos.conf"
const VRFCTL_CONF_PATH = "/etc/vrf.conf"

type Args struct {
	Path    string
	Cmd     string
	Verbose bool
	Args    []string
}

func (a *Args) Parse() {
	vrf := false
	flag.StringVar(&a.Path, "path", "", "config filename.")
	flag.StringVar(&a.Cmd, "cmd", "", "'set' or 'del'")
	flag.BoolVar(&vrf, "vrf", false, "vrf configuration mode.")
	flag.BoolVar(&a.Verbose, "v", false, "show detail message.")
	flag.Parse()
	a.Args = flag.Args()

	if len(a.Path) == 0 {
		a.Path = func() string {
			if vrf {
				return VRFCTL_CONF_PATH
			} else {
				return SYSCTL_CONF_PATH
			}
		}()
	}
}

func main() {
	args := Args{}
	args.Parse()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	conf := ncsclib.NewConfig()
	if err := prop.ReadFile(args.Path, conf); err != nil {
		log.Infof("new config created. %s", args.Path)
	}

	err := func() error {
		switch args.Cmd {
		case "set":
			return setSysctlConfig(conf, &args)
		case "del":
			return delSysctlConfig(conf, &args)
		default:
			return fmt.Errorf("Invalid cmd. %s", args.Cmd)
		}
	}()

	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	if err := prop.WriteFile(args.Path, conf); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func setSysctlConfig(cfg ncsclib.Config, args *Args) error {
	for _, line := range args.Args {
		k, v, err := prop.ParseLine(line)
		if err != nil {
			log.Errorf("%s", err)
			return err
		}

		cfg[k] = v
		log.Debugf("set '%s' = '%s'", k, v)
	}

	log.Debugf("setSysctlConfig OK")
	return nil
}

func delSysctlConfig(cfg ncsclib.Config, args *Args) error {
	for _, line := range args.Args {
		k, v, err := prop.ParseLine(line)
		if err != nil {
			return err
		}

		delete(cfg, k)
		log.Debugf("del '%s' = '%s'", k, v)
	}

	log.Debugf("delSysctlConfig OK")
	return nil
}
