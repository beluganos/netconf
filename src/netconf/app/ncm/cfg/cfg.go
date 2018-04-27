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
	"fmt"
)

//
// Cfg
//
type Cfg struct {
	Config *Config
	Opts   *Opts
}

func newCfg() *Cfg {
	return &Cfg{
		Config: newConfig(),
		Opts:   newOpts(),
	}
}

func (c *Cfg) String() string {
	return fmt.Sprintf("Cfg{%s, %s}", c.Config, c.Opts)
}

func (c *Cfg) Init() error {
	c.Opts.Parse()
	return c.Config.Load(c.Opts.Path)
}

var ncmCfg = newCfg()

func GetCfg() *Cfg {
	return ncmCfg
}

func GetConfig() *Config {
	return ncmCfg.Config
}

func GetOpts() *Opts {
	return ncmCfg.Opts
}
