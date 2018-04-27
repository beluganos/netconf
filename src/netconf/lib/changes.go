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

package nclib

import (
	"strings"
)

type SrChanges map[string]struct{}

func NewSrChanges() SrChanges {
	return SrChanges{}
}

func (c SrChanges) String() string {
	keys := []string{}
	for key, _ := range c {
		keys = append(keys, key)
	}
	return strings.Join(keys, "|")
}

func (c SrChanges) GetChange(key string) bool {
	_, ok := c[key]
	return ok
}

func (c SrChanges) GetChanges(keys ...string) bool {
	for _, key := range keys {
		if ok := c.GetChange(key); !ok {
			return false
		}
	}
	return true
}

func (c SrChanges) OneOfChange(keys ...string) bool {
	for _, key := range keys {
		if ok := c.GetChange(key); ok {
			return true
		}
	}
	return false
}

func (c SrChanges) SetChange(key string) {
	c[key] = struct{}{}
}

func (c SrChanges) SetChanges(keys ...string) {
	for _, key := range keys {
		c.SetChange(key)
	}
}

func (c SrChanges) Compare(keys ...string) bool {
	if len(keys) != len(c) {
		return false
	}

	for _, key := range keys {
		if !c.GetChange(key) {
			return false
		}
	}

	return true
}
