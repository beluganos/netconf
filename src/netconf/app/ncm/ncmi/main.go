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
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_CONFIG_NAME = "/etc/beluganos/netconf/nc-module-if.yaml"
	DEFAULT_CONFIG_TYPE = "yaml"
)

type Args struct {
	Path    string
	Type    string
	DpIds   map[DatapathId]struct{}
	Verbose bool
}

func NewArgs() *Args {
	return &Args{
		Path:    "",
		Type:    "",
		DpIds:   map[DatapathId]struct{}{},
		Verbose: false,
	}
}

func (a *Args) String() string {
	dpids := []DatapathId{}
	for dpid, _ := range a.DpIds {
		dpids = append(dpids, dpid)
	}
	return fmt.Sprintf("config='%s', type='%s', dpids=%v, verbose=%t", a.Path, a.Type, dpids, a.Verbose)
}

func (a *Args) Parse() {
	flag.StringVar(&a.Path, "f", DEFAULT_CONFIG_NAME, "config file name.")
	flag.StringVar(&a.Type, "type", DEFAULT_CONFIG_TYPE, "config file type.")
	flag.BoolVar(&a.Verbose, "v", false, "show detail message.")
	flag.Parse()

	for _, arg := range flag.Args() {
		dpid, err := strconv.ParseUint(arg, 0, 64)
		if err != nil {
			continue
		}
		a.DpIds[DatapathId(dpid)] = struct{}{}
	}
}

func main() {
	args := NewArgs()
	args.Parse()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Debugf("%s", args)

	s := NewConfigServer(args.DpIds)
	if err := s.Serve(args.Path, args.Type); err != nil {
		log.Errorf("ConfigServer.Serve error. %s", err)
		os.Exit(1)
	}
}
