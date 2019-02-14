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
	srlib "netconf/lib/sysrepo"
	"os"

	log "github.com/sirupsen/logrus"
)

type Args struct {
	NoCopy  bool
	Verbose bool
	Args    []string
}

func (a *Args) Parse() {
	flag.BoolVar(&a.NoCopy, "n", false, "do not copy config.")
	flag.BoolVar(&a.Verbose, "v", false, "show detail message.")
	flag.Parse()
	a.Args = flag.Args()
}

func main() {
	args := Args{}
	args.Parse()

	if args.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	conn := srlib.NewSrConnection()
	if err := conn.Connect("ncmi"); err != nil {
		log.Errorf("%s", err)
		return
	}
	defer conn.Disconnect()

	session := srlib.NewSrSession(nil)
	if err := session.Start(conn, srlib.SR_DS_RUNNING); err != nil {
		log.Errorf("%s", err)
		return
	}
	defer session.Stop()

	done := make(chan struct{})

	for _, arg := range args.Args {
		c := NewSubscribeController(session, args.NoCopy, args.Verbose)
		subscr, err := c.Start(arg, done)
		if err != nil {
			log.Errorf("NewSubscriber error. %s", err)
			os.Exit(1)
		}
		defer subscr.Stop()

		log.Infof("Subscriber(%s) STARTED", arg)
	}

	<-done
}
