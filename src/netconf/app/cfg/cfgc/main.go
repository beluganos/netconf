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
	api "netconf/app/cfg/api"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Args struct {
	Host    string
	Port    uint
	Cmd     string
	Verbose bool
	Args    []string
}

func (a *Args) Init() {
	flag.StringVar(&a.Host, "host", "localhost", "Host")
	flag.UintVar(&a.Port, "port", api.LISTEN_PORT, "Port")
	flag.StringVar(&a.Cmd, "cmd", "help", "command.")
	flag.BoolVar(&a.Verbose, "verbose", false, "Show detail messages.")
	flag.Parse()
	a.Args = flag.Args()
}

func main() {
	args := Args{}
	args.Init()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	client, conn, err := api.NewInsecureClient(args.Host, args.Port)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Debugf("Command %s %v", args.Cmd, args.Args)

	req := api.NewExecuteRequest(api.NewShell(args.Cmd, args.Args...))
	res, err := client.Execute(context.Background(), req)
	if err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	log.Debugf("Command Success. %s", string(res.Results[0].Output))
	os.Exit(0)
}
