// -+- coding: utf-8 -*-

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

func makeRoutingPolicy(datas [][2]string) *RoutingPolicy {
	policy := NewRoutingPolicy()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := policy.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return policy
}

func TestRoutingPolicy_Put(t *testing.T) {
	datas := [][2]string{
		{"/beluganos-routing-policy:routing-policy/policy-definitions/policy-definition[name='policy-next-hop-self']/statements/statement[name='stmt-next-hop-self2']/actions/beluganos-bgp-policy:bgp-actions/config/set-local-pref", "200"},
		{"/beluganos-routing-policy:routing-policy/policy-definitions/policy-definition[name='policy-next-hop-self']/statements/statement[name='stmt-next-hop-self2']/actions/beluganos-bgp-policy:bgp-actions/config/set-next-hop", "SELF"},
	}

	policy := NewRoutingPolicy()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := policy.Put(nodes[1:], value); err != nil {
			t.Errorf("RoutingPolicy.Put error. %s %s = %s", err, xpath, value)
		}
	}

	t.Log(policy)

	if v := len(policy.Definitions); v != 1 {
		t.Errorf("RoutingPolicy.Put unmatch. #defs=%d", v)
	}
}
