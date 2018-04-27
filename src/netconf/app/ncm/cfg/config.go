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
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

const (
	DEFAULT_CLI_PATH = "/usr/bin"
	NC_HOME_ENV      = "NC_HOME"
)

//
// Config - Global
//
type GlobalConfig struct {
	Persist bool `toml:"persist"`
}

func (c *GlobalConfig) String() string {
	return fmt.Sprintf("Global{persist=%t}", c.Persist)
}

//
// Config - cli
//
type CliConfig struct {
	Path    string `toml:"path"`
	LxcInit string `toml:"lxcinit"`
	Lxd     string `toml:"lxd"`
	Sys     string `toml:"sys"`
	Vty     string `toml:"vty"`
	GoBgp   string `toml:"gobgp"`
}

func (c *CliConfig) String() string {
	return fmt.Sprintf("CLI{path='%s'}", c.Path)
}

func (c *CliConfig) ToPath(name string) string {
	return path.Join(c.Path, name)
}

func (c *CliConfig) LxcInitPath() string {
	return c.ToPath(c.LxcInit)
}

func (c *CliConfig) LxdPath() string {
	return c.ToPath(c.Lxd)
}

func (c *CliConfig) SysPath() string {
	return c.ToPath(c.Sys)
}

func (c *CliConfig) GoBgpPath() string {
	return c.ToPath(c.GoBgp)
}

func (c *CliConfig) VtyPath() string {
	return c.ToPath(c.Vty)
}

//
// Config - FRR
//
type FrrConfig struct {
	AutoRstart string `toml:"auto_restart"` // 'restart', 'reload' or others(no operation)
}

func (c *FrrConfig) String() string {
	return fmt.Sprintf("Frr{AutoRestart='%s'}", c.AutoRstart)
}

//
// Config
//
type Config struct {
	Global *GlobalConfig `toml:"global"`
	Frr    *FrrConfig    `toml:"frr"`
	Cli    *CliConfig    `toml:"cli"`
}

func newConfig() *Config {
	return &Config{
		Global: &GlobalConfig{},
		Frr:    &FrrConfig{},
		Cli:    &CliConfig{},
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{%s, %s, %s}", c.Global, c.Frr, c.Cli)
}

func GetCliPathFromEnv() string {
	if path, ok := os.LookupEnv(NC_HOME_ENV); ok {
		return fmt.Sprintf("%s/bin", path)
	}

	return DEFAULT_CLI_PATH
}

func (c *Config) Load(path string) error {
	c.Cli.Path = GetCliPathFromEnv() // set default value
	if _, err := toml.DecodeFile(path, c); err != nil {
		return err
	}

	return nil
}
