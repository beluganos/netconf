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
	ncmcfg "netconf/app/ncm/cfg"
	ncmdbm "netconf/app/ncm/dbm"
	ncm "netconf/app/ncm/modules"
	ncsignal "netconf/lib/signal"
	srlib "netconf/lib/sysrepo"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func subscribeNetworkInstanceChange(s *srlib.SrSession) *srlib.Subscriber {
	factory := ncm.NewNIChangeFactory(ncmcfg.GetOpts().DryRun)
	ctrl := ncm.NewNIChangeController(s, factory)
	subscr, err := ctrl.Subscribe()
	if err != nil {
		log.Errorf("subscribeNetworkInstanceChange error. %s", err)
		os.Exit(1)
	}

	log.Infof("START: Subscriber(NetworkInstanceChange)")
	return subscr
}

func main() {
	if err := ncmcfg.GetCfg().Init(); err != nil {
		log.Errorf("Init Config error. %s", err)
		os.Exit(1)
	}

	if ncmcfg.GetOpts().Verbose {
		log.SetLevel(log.DebugLevel)
		log.Debugf("%s", ncmcfg.GetCfg())
	}

	conn := srlib.NewSrConnection()
	if err := conn.Connect("ncmd"); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
	defer conn.Disconnect()

	session := srlib.NewSrSession(nil)
	if err := session.Start(conn, srlib.SR_DS_STARTUP); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
	defer session.Stop()

	ncmdbm.Create(session)

	niSubscr := subscribeNetworkInstanceChange(session)
	defer niSubscr.Stop()

	ss := ncsignal.NewServer()
	ss.Register(syscall.SIGPIPE, func(sig os.Signal) {
		log.Infof("SIGNAL %s", sig)
	}).Serve(nil)
}
