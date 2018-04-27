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

func makeAddrs(datas [][2]string) IPAddresses {
	addrs := NewIPAddresses()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := addrs.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return addrs
}

//
// /addresses/address[ip='10.0.0.1']
//
func TestIPAddress_empty(t *testing.T) {
	addrs := makeAddrs([][2]string{
		{"/addresses/address[ip='10.0.0.1']", ""},
	})

	t.Log(addrs)

	// check
	if v := len(addrs); v != 1 {
		t.Errorf("IPAddress.Put unmatch. len=%d", v)
	}

	addr := addrs["10.0.0.1"]

	if v := addr.Compare(); !v {
		t.Errorf("IPAddress.Put unmatch. cmp=%t", v)
	}
}

//
// /addresses/address[ip='NOT_EXIST']
//
func TestIPAddress_ip_not_exist(t *testing.T) {
	addrs := NewIPAddresses()
	xpath := "/addresses/address"
	nodes := srlib.ParseXPath(xpath)

	if err := addrs.Put(nodes[1:], ""); err == nil {
		t.Errorf("addresses.Put must be error.")
	}
}

//
// /addresses/address[ip='10.0.0.1']/unknown
//
func TestIPAddress_unknown(t *testing.T) {
	addrs := NewIPAddresses()
	xpath := "/addresses/address[ip='10.0.0.1']/UNKNOWN"
	nodes := srlib.ParseXPath(xpath)

	if err := addrs.Put(nodes[1:], ""); err != nil {
		t.Errorf("addresses.Put error. %s", err)
	}

	addr := addrs["10.0.0.1"]

	if v := addr.Compare("UNKNOWN"); !v {
		t.Errorf("IPAddress.Put unmatch. cmp=%t", v)
	}
}

//
// /addresses/address[ip='10.0.0.1']/ip
// /addresses/address[ip='10.0.0.1']/config
//
func TestIPAddress(t *testing.T) {
	addrs := makeAddrs([][2]string{
		{"/addresses/address[ip='10.0.0.1']", ""},
		{"/addresses/address[ip='10.0.0.1']/ip", "10.0.0,1"},
		{"/addresses/address[ip='10.0.0.1']/config", ""},
	})

	t.Log(addrs)

	// check
	if v := len(addrs); v != 1 {
		t.Errorf("IPAddress.Put unmatch. len=%d", v)
	}

	addr := addrs["10.0.0.1"]

	if v := addr.Compare(OC_CONFIG_KEY, SUBINTERFACE_ADDR_IP_KEY); !v {
		t.Errorf("IPAddress.Put unmatch. cmp=%t", v)
	}
	if v := addr.IP; v != "10.0.0.1" {
		t.Errorf("IPAddress.Put unmatch. ip=%s", v)
	}
}

//
// /addresses/address[ip='10.0.0.1']/config
//
func TestIPAddress_config_empty(t *testing.T) {
	addrs := makeAddrs([][2]string{
		{"/addresses/address[ip='10.0.0.1']", ""},
		{"/addresses/address[ip='10.0.0.1']/config", ""},
	})

	t.Log(addrs)

	// check
	if v := len(addrs); v != 1 {
		t.Errorf("IPAddress.Put unmatch. len=%d", v)
	}

	addr := addrs["10.0.0.1"]

	if v := addr.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("IPAddress.Put unmatch. addr.cmp=%t", v)
	}
	if v := addr.Config.Compare(); !v {
		t.Errorf("IPAddress.Put unmatch. config.cmp=%t", v)
	}
}

//
// /addresses/address[ip='10.0.0.1']/config/ip = "10.0.0.1"
// /addresses/address[ip='10.0.0.1']/config/prefix-length = 24
//
func TestIPAddress_config(t *testing.T) {
	addrs := makeAddrs([][2]string{
		{"/addresses/address[ip='10.0.0.1']", ""},
		{"/addresses/address[ip='10.0.0.1']/config", ""},
		{"/addresses/address[ip='10.0.0.1']/config/ip", "10.0.0.1"},
		{"/addresses/address[ip='10.0.0.1']/config/prefix-length", "24"},
	})

	t.Log(addrs)

	// check
	if v := len(addrs); v != 1 {
		t.Errorf("IPAddress.Put unmatch. len=%d", v)
	}

	addr := addrs["10.0.0.1"]

	if v := addr.Config.Compare(SUBINTERFACE_ADDR_IP_KEY, SUBINTERFACE_ADDR_PREFIXLEN_KEY); !v {
		t.Errorf("IPAddress.Put unmatch. config.cmp=%t", v)
	}
	if v := addr.Config.IP.String(); v != "10.0.0.1" {
		t.Errorf("IPAddress.Put unmatch. config.ip=%s", v)
	}
	if v := addr.Config.PrefixLen; v != 24 {
		t.Errorf("IPAddress.Put unmatch. config.ip=%d", v)
	}
}

//
// /addresses/address[ip='10.0.0.1']/config/ip = "BAD_IP_ADDR"
// /addresses/address[ip='10.0.0.1']/config/prefix-length = "BAD_PREFIX_LEN"
//
func TestIPAddress_config_err(t *testing.T) {
	addrs := NewIPAddresses()

	nodes := srlib.ParseXPath("/addresses/address[ip='10.0.0.1']/config/ip")
	if err := addrs.Put(nodes[1:], "BAD_IP_ADDR"); err == nil {
		t.Errorf("addresses.Put must be error.")
	}

	nodes = srlib.ParseXPath("/addresses/address[ip='10.0.0.1']/config/prefix-length")
	if err := addrs.Put(nodes[1:], "BAD_PREFIX_LEN"); err == nil {
		t.Errorf("addresses.Put must be error.")
	}
}

//
// /addresses/address[ip='10.0.0.1']/config/ip = "10.0.0.1"
// /addresses/address[ip='10.0.0.1']/config/prefix-length = 24
//
func TestIPAddress_config_conv(t *testing.T) {
	addrs := makeAddrs([][2]string{
		{"/addresses/address[ip='10.0.0.1']", ""},
		{"/addresses/address[ip='10.0.0.1']/config", ""},
		{"/addresses/address[ip='10.0.0.1']/config/ip", "10.0.0.1"},
		{"/addresses/address[ip='10.0.0.1']/config/prefix-length", "24"},
	})

	t.Log(addrs)

	// check
	if v := len(addrs); v != 1 {
		t.Errorf("IPAddress.Put unmatch. len=%d", v)
	}

	addr := addrs["10.0.0.1"]

	if v := addr.Config.IFAddr(); v.String() != "10.0.0.1/24" {
		t.Errorf("IPAddress.Config.IFAddr unmatch. %s", v)
	}
}
