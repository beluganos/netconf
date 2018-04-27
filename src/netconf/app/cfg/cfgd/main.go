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
	"net"
	"netconf/lib/signal"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	arg := Args{}
	arg.Init()

	if arg.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	lis, err := net.Listen("tcp", arg.ListenAddr())
	if err != nil {
		log.Fatalf("failed to listen: %s %v", arg.ListenAddr(), err)
		os.Exit(1)
	}

	ss := ncsignal.NewServer()
	ss.Register(syscall.SIGPIPE, func(sig os.Signal) {
		log.Debug("SIGNAL: %s", sig)
	}).Start(nil)

	g := grpc.NewServer()
	RegisterRpcApiServer(g, NewRpcApiServer())
	if err := g.Serve(lis); err != nil {
		log.Errorf("grpc server error. %s", err)
		os.Exit(1)
	}
}
