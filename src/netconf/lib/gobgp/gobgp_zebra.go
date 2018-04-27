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

type Zebra Entries

func NewZebra(i interface{}) Zebra {
	return Zebra(NewEntries(i))
}

func (z Zebra) Config() ZebraConfig {
	if i, ok := z["config"]; ok {
		return NewZebraConfig(i)
	}
	return nil
}

func (z Zebra) SetConfig(v ZebraConfig) {
	z["config"] = v
}

type ZebraConfig Entries

func NewZebraConfig(i interface{}) ZebraConfig {
	return ZebraConfig(NewEntries(i))
}

func (z ZebraConfig) Enabled() bool {
	return convBool(z, "enabled")
}

func (z ZebraConfig) SetEnabled(enabled bool) {
	z["enabled"] = enabled
}

func (z ZebraConfig) Url() string {
	return convString(z, "url")
}

func (z ZebraConfig) SetUrl(url string) {
	z["url"] = url
}

func (z ZebraConfig) Version() uint8 {
	return uint8(convUint(z, "version"))
}

func (z ZebraConfig) SetVersion(v uint8) {
	z["version"] = v
}

func (z ZebraConfig) RedistributeRouteTypeList() []string {
	list := convList(z, "redistribute-route-type-list")
	types := []string{}
	if list != nil {
		for _, l := range list {
			types = append(types, l.(string))
		}
	}
	return types
}

func (z ZebraConfig) SetRedistributeRouteTypeList(types []string) {
	z["redistribute-route-type-list"] = types
}
