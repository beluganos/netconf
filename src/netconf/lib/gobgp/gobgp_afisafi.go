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

type AfiSafi Entries

func NewAfiSafis(i interface{}) []AfiSafi {
	safis := []AfiSafi{}

	switch i.(type) {
	case nil:
	default:
		for _, safi := range i.([]interface{}) {
			safis = append(safis, NewAfiSafi(safi))
		}
	}

	return safis
}

func NewAfiSafi(i interface{}) AfiSafi {
	return AfiSafi(NewEntries(i))
}

func RawAfiSafis(afisafis []AfiSafi) interface{} {
	list := make([]interface{}, len(afisafis))
	for index, afisafi := range afisafis {
		list[index] = Entries(afisafi).Raw()
	}
	return list
}

func (a AfiSafi) Config() AfiSafiConfig {
	return NewAfiSafiConfig(getValue(a, "config"))
}

type AfiSafiConfig Entries

func NewAfiSafiConfig(i interface{}) AfiSafiConfig {
	return AfiSafiConfig(NewEntries(i))
}

func (a AfiSafiConfig) AfiSafiName() string {
	return convString(a, "afi-safi-name")
}

func (a AfiSafiConfig) Enabled() bool {
	return convBool(a, "enabled")
}
