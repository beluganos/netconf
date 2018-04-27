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

package ncmcfg

import (
	"flag"
	"fmt"
)

const (
	DEFAULT_API_ADDR = "localhost:50056"
	DEFAULT_CONFIG   = "/etc/nc_module/nc_module.conf"
)

//
// Options
//
type Opts struct {
	Path    string
	DryRun  bool
	Verbose bool
}

func newOpts() *Opts {
	return &Opts{}
}

func (o *Opts) String() string {
	return fmt.Sprintf("Opts{Path='%s', Verbose=%t, DryRun=%t}", o.Path, o.Verbose, o.DryRun)
}

func (o *Opts) Parse() {
	flag.BoolVar(&o.Verbose, "v", false, "show detail log.")
	flag.BoolVar(&o.DryRun, "d", false, "dry run.")
	flag.StringVar(&o.Path, "c", DEFAULT_CONFIG, "Config file.")
	flag.Parse()
}
