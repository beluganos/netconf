// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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

func makeNetworkInstanceLoopbacks(datas [][2]string) (NetworkInstanceLoopbacks, error) {
	los := NewNetworkInstanceLoopbacks()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := los.Put(nodes[1:], value); err != nil {
			return nil, err
		}
	}
	return los, nil
}

func makeNetworkInstanceLoopbackAddrs(datas [][2]string) (NetworkInstanceLoopbackAddrs, error) {
	addrs := NewNetworkInstanceLoopbackAddrs()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := addrs.Put(nodes[1:], value); err != nil {
			return nil, err
		}
	}
	return addrs, nil
}

func TestNetworkInstanceLoopbackAddr(t *testing.T) {
	addrs, err := makeNetworkInstanceLoopbackAddrs([][2]string{
		{"/addresses/address[index='1']", ""},
	})

	if err != nil {
		t.Errorf("lo-address.Put error. %s", err)
	}

	addr := addrs["1"]
	t.Log(addr)

	if v := addr.Compare(); !v {
		t.Errorf("lo-address,Put unmatch. cmp=%t", v)
	}
}

func TestNetworkInstanceLoopbackAddr_config(t *testing.T) {
	addrs, err := makeNetworkInstanceLoopbackAddrs([][2]string{
		{"/addresses/address[index='1']", ""},
		{"/addresses/address[index='1']/index", "1"},
		{"/addresses/address[index='1']/config", ""},
	})

	if err != nil {
		t.Errorf("lo-address.Put error. %s", err)
	}

	addr := addrs["1"]
	t.Log(addr)

	if v := addr.Compare(OC_INDEX_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("lo-address,Put unmatch. cmp=%t", v)
	}
	if v := addr.Config.Compare(); !v {
		t.Errorf("lo-address,config.Put unmatch. cmp=%t", v)
	}
}

func TestNetworkInstanceLoopbackAddr_config_ip(t *testing.T) {
	addrs, err := makeNetworkInstanceLoopbackAddrs([][2]string{
		{"/addresses/address[index='1']", ""},
		{"/addresses/address[index='1']/index", "1"},
		{"/addresses/address[index='1']/config", ""},
		{"/addresses/address[index='1']/config/index", "1"},
		{"/addresses/address[index='1']/config/ip", "10.0.0.1"},
	})

	if err != nil {
		t.Errorf("lo-address.Put error. %s", err)
	}

	addr := addrs["1"]
	t.Log(addr)

	if v := addr.Compare(OC_INDEX_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("lo-address,Put unmatch. cmp=%t", v)
	}
	if v := addr.Config.Compare(OC_INDEX_KEY, NETWORKINSTANCE_LO_IP_KEY); !v {
		t.Errorf("lo-address,config.Put unmatch. cmp=%t", v)
	}
}

func TestNetworkInstanceLoopbackAddr_config_ipv6(t *testing.T) {
	addrs, err := makeNetworkInstanceLoopbackAddrs([][2]string{
		{"/addresses/address[index='1']", ""},
		{"/addresses/address[index='1']/index", "1"},
		{"/addresses/address[index='1']/config", ""},
		{"/addresses/address[index='1']/config/index", "1"},
		{"/addresses/address[index='1']/config/ip", "2001:2001::1"},
	})

	if err != nil {
		t.Errorf("lo-address.Put error. %s", err)
	}

	addr := addrs["1"]
	t.Log(addr)

	if v := addr.Compare(OC_INDEX_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("lo-address,Put unmatch. cmp=%t", v)
	}
	if v := addr.Config.Compare(OC_INDEX_KEY, NETWORKINSTANCE_LO_IP_KEY); !v {
		t.Errorf("lo-address,config.Put unmatch. cmp=%t", v)
	}
}

func TestNetworkInstanceLoopbackAddr_config_plen(t *testing.T) {
	addrs, err := makeNetworkInstanceLoopbackAddrs([][2]string{
		{"/addresses/address[index='1']", ""},
		{"/addresses/address[index='1']/index", "1"},
		{"/addresses/address[index='1']/config", ""},
		{"/addresses/address[index='1']/config/index", "1"},
		{"/addresses/address[index='1']/config/prefix-length", "64"},
	})

	if err != nil {
		t.Errorf("lo-address.Put error. %s", err)
	}

	addr := addrs["1"]
	t.Log(addr)

	if v := addr.Compare(OC_INDEX_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("lo-address,Put unmatch. cmp=%t", v)
	}
	if v := addr.Config.Compare(OC_INDEX_KEY, NETWORKINSTANCE_LO_PLEN_KEY); !v {
		t.Errorf("lo-address,config.Put unmatch. cmp=%t", v)
	}
}

func TestNetworkInstanceLoopbackAddr_x2(t *testing.T) {
	addrs, err := makeNetworkInstanceLoopbackAddrs([][2]string{
		{"/addresses/address[index='1']", ""},
		{"/addresses/address[index='1']/index", "1"},
		{"/addresses/address[index='1']/config", ""},
		{"/addresses/address[index='1']/config/index", "1"},
		{"/addresses/address[index='1']/config/ip", "10.0.0.1"},

		{"/addresses/address[index='2']", ""},
		{"/addresses/address[index='2']/index", "2"},
		{"/addresses/address[index='2']/config", ""},
		{"/addresses/address[index='2']/config/index", "2"},
		{"/addresses/address[index='2']/config/ip", "2001:2001::1"},
	})

	if err != nil {
		t.Errorf("lo-address.Put error. %s", err)
	}

	addr := addrs["1"]
	t.Log(addr)

	if v := addr.Compare(OC_INDEX_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("lo-address,Put unmatch. cmp=%t", v)
	}
	if v := addr.Config.Compare(OC_INDEX_KEY, NETWORKINSTANCE_LO_IP_KEY); !v {
		t.Errorf("lo-address,config.Put unmatch. cmp=%t", v)
	}

	addr = addrs["2"]
	t.Log(addr)

	if v := addr.Compare(OC_INDEX_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("lo-address,Put unmatch. cmp=%t", v)
	}
	if v := addr.Config.Compare(OC_INDEX_KEY, NETWORKINSTANCE_LO_IP_KEY); !v {
		t.Errorf("lo-address,config.Put unmatch. cmp=%t", v)
	}
}

func TestNewNetworkInstanceLoopbacks(t *testing.T) {
	los, err := makeNetworkInstanceLoopbacks([][2]string{
		{"/loopbacks/loopback[id='lo']", ""},
	})

	if err != nil {
		t.Errorf("loopbacks.Put error. %s", err)
	}

	lo := los["lo"]
	t.Log(lo)

	if v := lo.Compare(); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}
}

func TestNewNetworkInstanceLoopbacks_config(t *testing.T) {
	los, err := makeNetworkInstanceLoopbacks([][2]string{
		{"/loopbacks/loopback[id='lo']", ""},
		{"/loopbacks/loopback[id='lo']/id", "lo"},
		{"/loopbacks/loopback[id='lo']/config", ""},
	})

	if err != nil {
		t.Errorf("loopbacks.Put error. %s", err)
	}

	lo := los["lo"]
	t.Log(lo)

	if v := lo.Compare(OC_ID_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}
	if v := lo.Config.Compare(); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}
}

func TestNewNetworkInstanceLoopbacks_config_id(t *testing.T) {
	los, err := makeNetworkInstanceLoopbacks([][2]string{
		{"/loopbacks/loopback[id='lo']", ""},
		{"/loopbacks/loopback[id='lo']/id", "lo"},
		{"/loopbacks/loopback[id='lo']/config", ""},
		{"/loopbacks/loopback[id='lo']/config/id", "lo"},
	})

	if err != nil {
		t.Errorf("loopbacks.Put error. %s", err)
	}

	lo := los["lo"]
	t.Log(lo)

	if v := lo.Compare(OC_ID_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}
	if v := lo.Config.Compare(OC_ID_KEY); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}
}

func TestNewNetworkInstanceLoopbacks_addrs(t *testing.T) {
	los, err := makeNetworkInstanceLoopbacks([][2]string{
		{"/loopbacks/loopback[id='lo']", ""},
		{"/loopbacks/loopback[id='lo']/id", "lo"},
		{"/loopbacks/loopback[id='lo']/config", ""},
		{"/loopbacks/loopback[id='lo']/config/id", "lo"},
		{"/loopbacks/loopback[id='lo']/addresses", ""},
		{"/loopbacks/loopback[id='lo']/addresses/address[index='0']", ""},
	})

	if err != nil {
		t.Errorf("loopbacks.Put error. %s", err)
	}

	lo := los["lo"]
	t.Log(lo)

	if v := lo.Compare(OC_ID_KEY, OC_CONFIG_KEY, NETWORKINSTANCE_LO_ADDRS_KEY); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}
	if v := lo.Config.Compare(OC_ID_KEY); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}

	addr := lo.Addrs["0"]
	if v := addr.Compare(); !v {
		t.Errorf("loopbacks.Put unmatch. cmp=%t", v)
	}
}
