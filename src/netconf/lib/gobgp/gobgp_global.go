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

package ncgobgp

//
// [global]
//
type Global Entries

func NewGlobal(i interface{}) Global {
	return Global(NewEntries(i))
}

func (g Global) Config() GlobalConfig {
	return NewGlobalConfig(getValue(g, "config"))
}

//
// [global.config]
//
type GlobalConfig Entries

func NewGlobalConfig(i interface{}) GlobalConfig {
	return GlobalConfig(NewEntries(i))
}

func (g GlobalConfig) As() uint32 {
	return uint32(convUint(g, "as"))
}

func (g GlobalConfig) SetAs(v uint32) {
	g["as"] = v
}

func (g GlobalConfig) RouterId() string {
	return convString(g, "router-id")
}

func (g GlobalConfig) SetRouterId(v string) {
	g["router-id"] = v
}
