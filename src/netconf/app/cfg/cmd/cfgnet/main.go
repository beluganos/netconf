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
	"fmt"
	ncnplib "netconf/lib/netplan"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	args := &Args{}
	if err := args.Parse(); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	conf, err := ncnplib.ReadConfigFile(args.Path)
	if err != nil {
		conf = ncnplib.NewConfig()
		log.Warnf("config file created. %s", args.Path)
	}

	err = func() error {
		switch args.Cmd {
		case CMD_SET:
			return setConfig(conf, args)
		case CMD_DEL:
			return delConfig(conf, args)
		default:
			return fmt.Errorf("Invalid commnd. %s", args.Cmd)
		}
	}()

	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	if err := ncnplib.WriteConfigFile(args.Path, conf); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
