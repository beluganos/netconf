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

package openconfig

import (
	srlib "netconf/lib/sysrepo"
	"testing"
)

func makeNwInstances(datas [][2]string) NetworkInstances {
	insts := NewNetworkInstances()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := insts.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return insts
}

func TestNwInstance(t *testing.T) {
	datas := [][2]string{
		{"/beluganos-network-instance:network-instances/network-instance[name='PE1']/interfaces/interface[id='eth1.10']/config/subinterface", "10"},
	}

	insts := makeNwInstances(datas)
	t.Log(insts)
}

func TestNwInstanceLDP(t *testing.T) {
	ldp_path := "/beluganos-network-instance:network-instances/network-instance[name='PE1']/mpls/signaling-protocols/ldp"
	datas := [][2]string{
		{ldp_path + "/router-id", "10.0.0.1"},
		{ldp_path + "/address-family/ipv4/discovery/transport-address", "10.0.0.2"},
	}
	insts := makeNwInstances(datas)
	t.Log(insts)
}
