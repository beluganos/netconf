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

type Transport Entries

func NewTransport(i interface{}) Transport {
	return Transport(NewEntries(i))
}

func (t Transport) Config() TransportConfig {
	return NewTransportConfig(getValue(t, "transport"))
}

func (t Transport) SetConfig(v TransportConfig) {
	t["transport"] = v
}

type TransportConfig Entries

func NewTransportConfig(i interface{}) TransportConfig {
	return TransportConfig(NewEntries(i))
}

func (c TransportConfig) LocalAddress() string {
	return convString(c, "ocal-address")
}

func (c TransportConfig) SetLocalAddress(v string) {
	c["ocal-address"] = v
}
