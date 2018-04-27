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

package cfgbelugcmd

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type RibsConfig map[string]interface{}

func (c RibsConfig) Ribs() (map[string]interface{}, error) {
	ribs, ok := c["ribs"]
	if !ok {
		return nil, fmt.Errorf("[ribs] not found.")
	}

	return ribs.(map[string]interface{}), nil
}

func (c RibsConfig) RibsVrf() (map[string]interface{}, error) {
	ribs, err := c.Ribs()
	if err != nil {
		return nil, err
	}

	vrf, ok := ribs["vrf"]
	if !ok {
		return nil, fmt.Errorf("[ribs.vrf] not found.")
	}

	return vrf.(map[string]interface{}), nil
}

func (c RibsConfig) SetRibsVrf(key string, val interface{}) error {
	vrf, err := c.RibsVrf()
	if err != nil {
		return err
	}

	vrf[key] = val
	return nil
}

// func ReadConfig(path string, cfg *Config) error {
func ReadConfig(path string, cfg *RibsConfig) error {
	_, err := toml.DecodeFile(path, cfg)
	return err
}

// func WriteConfig(path string, cfg *Config) error {
func WriteConfig(path string, cfg *RibsConfig) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(cfg)
}
