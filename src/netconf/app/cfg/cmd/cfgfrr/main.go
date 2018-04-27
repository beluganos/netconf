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
	"netconf/lib/vty"
	"os"

	log "github.com/sirupsen/logrus"
)

const DAEMONS_PATH = "/etc/frr/daemons"

type Args struct {
	Path    string
	Cmd     string
	Verbose bool
	Args    []string
}

func (a *Args) Parse() {
	flag.StringVar(&a.Path, "path", DAEMONS_PATH, "daemons filename")
	flag.StringVar(&a.Cmd, "cmd", "", "set or del")
	flag.BoolVar(&a.Verbose, "v", false, "show detail message")
	flag.Parse()
	a.Args = flag.Args()
}

func main() {
	args := Args{}
	args.Parse()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	cfg := vtylib.NewDaemonConfig()
	if err := prop.ReadFile(args.Path, cfg); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	err := func() error {
		switch args.Cmd {
		case "set":
			return setDaemons(cfg, &args)
		case "del":
			return delDaemons(cfg, &args)
		default:
			flag.PrintDefaults()
			return fmt.Errorf("Invalid command. %s", args.Cmd)
		}
	}()

	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	if err := prop.WriteFile(args.Path, cfg); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func setDaemons(cfg vtylib.DaemonConfig, args *Args) error {
	return setDaemonConfig(cfg, args, vtylib.DAEMON_STATE_YES)
}

func delDaemons(cfg vtylib.DaemonConfig, args *Args) error {
	return setDaemonConfig(cfg, args, vtylib.DAEMON_STATE_NO)
}

func setDaemonConfig(cfg vtylib.DaemonConfig, args *Args, state vtylib.DaemonState) error {
	for _, arg := range args.Args {
		if _, err := vtylib.ParseDaemonType(arg); err != nil {
			return err
		}
		cfg[arg] = state
	}
	return nil
}
