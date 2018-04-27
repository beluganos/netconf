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
	api "netconf/app/cfg/api"
)

type Args struct {
	Host    string
	Port    uint
	Verbose bool
}

func (a *Args) Init() {
	flag.StringVar(&a.Host, "listen", "0.0.0.0", "listen address.")
	flag.UintVar(&a.Port, "port", api.LISTEN_PORT, "port number.")
	flag.BoolVar(&a.Verbose, "verbose", false, "show detail message.")
	flag.Parse()
}

func (a *Args) ListenAddr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
