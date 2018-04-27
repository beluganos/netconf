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

package ncsclib

import (
	"fmt"
	"strings"
)

type Config map[string]string

func NewConfig() Config {
	return Config{}
}

func (c Config) Set(k string, v string) error {
	c[k] = v
	return nil
}

func (c Config) Get(f func(string) error) error {
	for k, v := range c {
		line := fmt.Sprintf("%s = %s", k, v)
		if err := f(line); err != nil {
			return err
		}
	}

	return nil
}

func FixString(s string) string {
	return strings.Replace(s, ".", "/", -1)
}
