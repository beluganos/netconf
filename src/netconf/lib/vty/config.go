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

package vtylib

import (
	"fmt"
)

type DaemonConfig map[string]DaemonState

func NewDaemonConfig() DaemonConfig {
	return DaemonConfig{}
}

func (d DaemonConfig) Set(name string, state string) error {
	v, err := ParseDaemonState(state)
	if err != nil {
		return err
	}

	d[name] = v
	return nil
}

func (d DaemonConfig) Get(f func(string) error) error {
	for name, state := range d {
		line := fmt.Sprintf("%s=%s", name, state)
		if err := f(line); err != nil {
			return err
		}
	}
	return nil
}
