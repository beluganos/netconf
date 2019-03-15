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
	"fmt"
	srlib "netconf/lib/sysrepo"
	"testing"
)

func makeSubinterfaceIPv4(datas [][2]string) *SubinterfaceIPv4 {
	ipv4 := NewSubinterfaceIPv4()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := ipv4.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return ipv4
}

func makeSubinterfaceIPv6(datas [][2]string) *SubinterfaceIPv6 {
	ipv6 := NewSubinterfaceIPv6()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := ipv6.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return ipv6
}

func TestSubinterfaceIPv4(t *testing.T) {
	ipv4 := makeSubinterfaceIPv4([][2]string{
		{"/openconfig-if-ip:ipv4", ""},
	})

	t.Log(ipv4)

	if v := ipv4.Compare(); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. cmp=%t", v)
	}
}

func TestSubinterfaceIPv4_config(t *testing.T) {
	ipv4 := makeSubinterfaceIPv4([][2]string{
		{"/openconfig-if-ip:ipv4", ""},
		{"/openconfig-if-ip:ipv4/config", ""},
	})

	t.Log(ipv4)

	if v := ipv4.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. cmp=%t", v)
	}
	if v := ipv4.Config.Compare(); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. condig.cmp=%t", v)
	}
}

func TestSubinterfaceIPv4_config_mtu(t *testing.T) {
	mtu := uint16(100)
	ipv4 := makeSubinterfaceIPv4([][2]string{
		{"/openconfig-if-ip:ipv4", ""},
		{"/openconfig-if-ip:ipv4/config", ""},
		{"/openconfig-if-ip:ipv4/config/mtu", fmt.Sprintf("%d", mtu)},
	})

	t.Log(ipv4)

	if v := ipv4.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. cmp=%t", v)
	}
	if v := ipv4.Config.Compare(SUBINTERFACE_MTU_KEY); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. config.cmp=%t", v)
	}
	if v := ipv4.Config.Mtu; v != mtu {
		t.Errorf("subaddr.ipv4.Push unmatch. config.mtu=%d", v)
	}
}

func TestSubinterfaceIPv4_addresses_x0(t *testing.T) {
	ipv4 := makeSubinterfaceIPv4([][2]string{
		{"/openconfig-if-ip:ipv4", ""},
		{"/openconfig-if-ip:ipv4/addresses", ""},
	})

	t.Log(ipv4)

	if v := ipv4.Compare(SUBINTERFACE_ADDRS_KEY); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. cmp=%t", v)
	}
	if v := len(ipv4.Addresses); v != 0 {
		t.Errorf("subaddr.ipv4.Push unmatch. #addresses=%d", v)
	}
}

func TestSubinterfaceIPv4_addresses_x1(t *testing.T) {
	ipv4 := makeSubinterfaceIPv4([][2]string{
		{"/openconfig-if-ip:ipv4", ""},
		{"/openconfig-if-ip:ipv4/addresses", ""},
		{"/openconfig-if-ip:ipv4/addresses/address[ip='10.0.0.1']", ""},
	})

	t.Log(ipv4)

	if v := ipv4.Compare(SUBINTERFACE_ADDRS_KEY); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. cmp=%t", v)
	}
	if v := len(ipv4.Addresses); v != 1 {
		t.Errorf("subaddr.ipv4.Push unmatch. #addresses=%d", v)
	}
	if _, ok := ipv4.Addresses["10.0.0.1"]; !ok {
		t.Errorf("subaddr.ipv4.Push unmatch. addresses['10.0.0.1']=%t", ok)
	}
}

func TestSubinterfaceIPv4_addresses_x2(t *testing.T) {
	ipv4 := makeSubinterfaceIPv4([][2]string{
		{"/openconfig-if-ip:ipv4", ""},
		{"/openconfig-if-ip:ipv4/addresses", ""},
		{"/openconfig-if-ip:ipv4/addresses/address[ip='10.0.0.1']", ""},
		{"/openconfig-if-ip:ipv4/addresses/address[ip='10.0.1.1']", ""},
	})

	t.Log(ipv4)

	if v := ipv4.Compare(SUBINTERFACE_ADDRS_KEY); !v {
		t.Errorf("subaddr.ipv4.Push unmatch. cmp=%t", v)
	}
	if v := len(ipv4.Addresses); v != 2 {
		t.Errorf("subaddr.ipv4.Push unmatch. #addresses=%d", v)
	}
	if _, ok := ipv4.Addresses["10.0.0.1"]; !ok {
		t.Errorf("subaddr.ipv4.Push unmatch. addresses['10.0.0.1']=%t", ok)
	}
	if _, ok := ipv4.Addresses["10.0.1.1"]; !ok {
		t.Errorf("subaddr.ipv4.Push unmatch. addresses['10.0.1.1']=%t", ok)
	}
}

func TestSubinterfaceIPv6(t *testing.T) {
	ipv6 := makeSubinterfaceIPv6([][2]string{
		{"/openconfig-if-ip:ipv6", ""},
	})

	t.Log(ipv6)

	if v := ipv6.Compare(); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. cmp=%t", v)
	}
}

func TestSubinterfaceIPv6_config(t *testing.T) {
	ipv6 := makeSubinterfaceIPv4([][2]string{
		{"/openconfig-if-ip:ipv6", ""},
		{"/openconfig-if-ip:ipv6/config", ""},
	})

	t.Log(ipv6)

	if v := ipv6.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. cmp=%t", v)
	}
	if v := ipv6.Config.Compare(); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. condig.cmp=%t", v)
	}
}

func TestSubinterfaceIPv6_config_mtu(t *testing.T) {
	mtu := uint16(100)
	ipv6 := makeSubinterfaceIPv6([][2]string{
		{"/openconfig-if-ip:ipv6", ""},
		{"/openconfig-if-ip:ipv6/config", ""},
		{"/openconfig-if-ip:ipv6/config/mtu", fmt.Sprintf("%d", mtu)},
	})

	t.Log(ipv6)

	if v := ipv6.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. cmp=%t", v)
	}
	if v := ipv6.Config.Compare(SUBINTERFACE_MTU_KEY); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. config.cmp=%t", v)
	}
	if v := ipv6.Config.Mtu; v != mtu {
		t.Errorf("subaddr.ipv6.Push unmatch. config.mtu=%d", v)
	}
}

func TestSubinterfaceIPv6_addresses_x0(t *testing.T) {
	ipv6 := makeSubinterfaceIPv6([][2]string{
		{"/openconfig-if-ip:ipv6", ""},
		{"/openconfig-if-ip:ipv6/addresses", ""},
	})

	t.Log(ipv6)

	if v := ipv6.Compare(SUBINTERFACE_ADDRS_KEY); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. cmp=%t", v)
	}
	if v := len(ipv6.Addresses); v != 0 {
		t.Errorf("subaddr.ipv6.Push unmatch. #addresses=%d", v)
	}
}

func TestSubinterfaceIPv6_addresses_x1(t *testing.T) {
	ipv6 := makeSubinterfaceIPv6([][2]string{
		{"/openconfig-if-ip:ipv6", ""},
		{"/openconfig-if-ip:ipv6/addresses", ""},
		{"/openconfig-if-ip:ipv6/addresses/address[ip='2001:2001::']", ""},
	})

	t.Log(ipv6)

	if v := ipv6.Compare(SUBINTERFACE_ADDRS_KEY); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. cmp=%t", v)
	}
	if v := len(ipv6.Addresses); v != 1 {
		t.Errorf("subaddr.ipv6.Push unmatch. #addresses=%d", v)
	}
	if _, ok := ipv6.Addresses["2001:2001::"]; !ok {
		t.Errorf("subaddr.ipv6.Push unmatch. addresses['2001:2001::']=%t", ok)
	}
}

func TestSubinterfaceIPv6_addresses_x2(t *testing.T) {
	ipv6 := makeSubinterfaceIPv6([][2]string{
		{"/openconfig-if-ip:ipv6", ""},
		{"/openconfig-if-ip:ipv6/addresses", ""},
		{"/openconfig-if-ip:ipv6/addresses/address[ip='2001:2001::']", ""},
		{"/openconfig-if-ip:ipv6/addresses/address[ip='2001:2001::1']", ""},
	})

	t.Log(ipv6)

	if v := ipv6.Compare(SUBINTERFACE_ADDRS_KEY); !v {
		t.Errorf("subaddr.ipv6.Push unmatch. cmp=%t", v)
	}
	if v := len(ipv6.Addresses); v != 2 {
		t.Errorf("subaddr.ipv6.Push unmatch. #addresses=%d", v)
	}
	if _, ok := ipv6.Addresses["2001:2001::"]; !ok {
		t.Errorf("subaddr.ipv6.Push unmatch. addresses['2001:2001::']=%t", ok)
	}
	if _, ok := ipv6.Addresses["2001:2001::1"]; !ok {
		t.Errorf("subaddr.ipv6.Push unmatch. addresses['2001:2001::1']=%t", ok)
	}
}
