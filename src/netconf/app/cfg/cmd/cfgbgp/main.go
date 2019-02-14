// -*- coding; utf-8 -*-

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
	"encoding/json"
	"flag"
	"io"
	ncgobgp "netconf/lib/gobgp"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	GOBGP_CONF_PATH  = "/etc/frr/gobgpd.conf"
	GOBGP_CONF_TYPE  = "toml"
	GOBGP_CONF_STDIN = "-"
)

type Args struct {
	Path    string
	Input   string
	Output  string
	Type    string
	Cmd     string
	Verbose bool
}

func (a *Args) Parse() {
	flag.StringVar(&a.Path, "c", GOBGP_CONF_PATH, "config file name")
	flag.StringVar(&a.Input, "i", GOBGP_CONF_STDIN, "input file name")
	flag.StringVar(&a.Output, "o", "", "output file name")
	flag.StringVar(&a.Type, "type", GOBGP_CONF_TYPE, "config file type")
	flag.StringVar(&a.Cmd, "cmd", "show", "command.(show/add/del)")
	flag.BoolVar(&a.Verbose, "v", false, "show detail message.")
	flag.Parse()
}

func inputConfig(in, t string) *ncgobgp.Config {
	var r io.Reader = nil

	if len(in) == 0 || in == GOBGP_CONF_STDIN {
		r = os.Stdin

		log.Debugf("Read from stdin.")

	} else {
		f, err := os.Open(in)
		if err != nil {
			log.Errorf("%s", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f

		log.Debugf("Read from %s", in)
	}

	cfg, err := ncgobgp.ReadConfig(r, t)
	if err != nil {
		log.Errorf("ReadConfig error. %s", err)
		os.Exit(1)
	}

	log.Debugf("ReadConfig success.")
	return cfg
}

func setConfig(cfg *ncgobgp.Config, in string, inType string) {
	cfg.Merge(inputConfig(in, inType))

	log.Debugf("config added.")
}

func delConfig(cfg *ncgobgp.Config, in string, inType string) {
	cfg.Delete(inputConfig(in, inType))

	log.Debugf("config deleted.")
}

func showConfig(cfg *ncgobgp.Config) {
	b, err := json.MarshalIndent(cfg.AllSettings(), "", "  ")
	if err != nil {
		log.Errorf("Marshal error. %s", err)
		os.Exit(1)
	}

	log.Infof("%s", b)
}

func writeConfig(cfg *ncgobgp.Config, path string) {
	err := func() error {
		if len(path) == 0 {
			log.Debugf("Write to same config.")

			return cfg.WriteConfig()

		} else {
			log.Debugf("Write to %s", path)

			return cfg.WriteConfigAs(path)

		}
	}()
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	log.Debugf("WriteConig ok.")
}

func showUsage() {
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	args := Args{}
	args.Parse()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	cfg, err := ncgobgp.ReadConfigFile(args.Path, args.Type)
	if err != nil {
		log.Errorf("ReadConfigFile error %s", err)
		os.Exit(1)
	}

	switch args.Cmd {
	case "set":
		setConfig(cfg, args.Input, args.Type)
		writeConfig(cfg, args.Output)
	case "del":
		delConfig(cfg, args.Input, args.Type)
		writeConfig(cfg, args.Output)
	case "show":
		showConfig(cfg)
	default:
		log.Errorf("Unknown cmd. %s", args.Cmd)
		showUsage()
	}
}
