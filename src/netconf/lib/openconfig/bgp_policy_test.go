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
	"netconf/lib/sysrepo"
	"testing"
)

func makeBgpActions(datas [][2]string) *PolicyBgpActions {
	acts := NewPolicyBgpActions()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := acts.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return acts
}

func TestBgpActions(t *testing.T) {
	a := makeBgpActions([][2]string{
		{"/beluganos-bgp-policy:bgp-actions", ""},
	})

	if v := a.Compare(); !v {
		t.Errorf("BgpActions_put unmatch. compare=%t", v)
	}
}

func TestBgpActions_config(t *testing.T) {
	a := makeBgpActions([][2]string{
		{"/beluganos-bgp-policy:bgp-actions/config", ""},
	})

	if v := a.Compare("config"); !v {
		t.Errorf("BgpActions_put unmatch. compare=%t", v)
	}
}

func TestBgpActions_config_set_local_pref(t *testing.T) {
	a := makeBgpActions([][2]string{
		{"/beluganos-bgp-policy:bgp-actions/config/set-local-pref", "200"},
	})

	if v := a.Config.Compare("set-local-pref"); !v {
		t.Errorf("BgpActions_put unmatch. config.compare=%t", v)
	}
	if v := a.Config.SetLocalPref; v != 200 {
		t.Errorf("BgpActions_put unmatch. config.set-local-pref=%d", v)
	}
}

func TestBgpActions_config_set_local_pref_error(t *testing.T) {
	xpath := "/beluganos-bgp-policy:bgp-actions/config/set-local-pref"
	value := "abc"
	a := NewPolicyBgpActions()
	nodes := srlib.ParseXPath(xpath)

	if err := a.Put(nodes[1:], value); err == nil {
		t.Errorf("BgpActions_put must be error. %s", err)
	}
}

func TestBgpActions_config_set_nexthop_ip(t *testing.T) {
	a := makeBgpActions([][2]string{
		{"/beluganos-bgp-policy:bgp-actions/config/set-next-hop", "10.0.1.2"},
	})

	if v := a.Config.Compare("set-next-hop"); !v {
		t.Errorf("BgpActions_put unmatch. config.compare=%t", v)
	}
	if v := a.Config.SetNexthop; v != "10.0.1.2" {
		t.Errorf("BgpActions_put unmatch. config.set-next-hop=%s", v)
	}
}

func TestBgpActions_config_set_nexthop_self(t *testing.T) {
	a := makeBgpActions([][2]string{
		{"/beluganos-bgp-policy:bgp-actions/config/set-next-hop", "SELF"},
	})

	if v := a.Config.Compare("set-next-hop"); !v {
		t.Errorf("BgpActions_put unmatch. config.compare=%t", v)
	}
	if v := a.Config.SetNexthop; v != "SELF" {
		t.Errorf("BgpActions_put unmatch. config.set-next-hop=%s", v)
	}
}

func TestBgpActions_config_set_nexthop_self_error(t *testing.T) {
	xpath := "/beluganos-bgp-policy:bgp-actions/config/set-next-hop"
	value := "self"
	a := NewPolicyBgpActions()
	nodes := srlib.ParseXPath(xpath)

	if err := a.Put(nodes[1:], value); err == nil {
		t.Errorf("BgpActions_put must be error. %s", err)
	}
}
