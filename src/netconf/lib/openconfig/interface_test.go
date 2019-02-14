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
	ncianalib "netconf/lib/iana"
	srlib "netconf/lib/sysrepo"
	"testing"
)

func makeIfaces(datas [][2]string) Interfaces {
	ifaces := NewInterfaces()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := ifaces.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return ifaces
}

func TestInterfaces(t *testing.T) {
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 0 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}
}

func TestInterface(t *testing.T) {
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
}

func TestInterface_name(t *testing.T) {
	ifname := "eth1"
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/name", ifname},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(OC_NAME_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := iface.Name; v != ifname {
		t.Errorf("ifaces.Put unmatch. iface.iface=%s", v)
	}
}

func TestInterface_config(t *testing.T) {
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := iface.Config.Compare(); !v {
		t.Errorf("ifaces.Put unmatch. iface.config.cmp=%t", v)
	}
}

func TestInterface_config_name(t *testing.T) {
	ifname := "eth1"
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config/name", ifname},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := iface.Config.Compare(OC_NAME_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.config.cmp=%t", v)
	}
	if v := iface.Config.Name; v != ifname {
		t.Errorf("ifaces.Put unmatch. iface.config.ifname=%s", ifname)
	}
}

func TestInterface_config_type(t *testing.T) {
	iftype := ncianalib.IANAifType_ethernetCsmacd
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config/type", iftype.String()},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := iface.Config.Compare(INTERFACE_TYPE_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.config.cmp=%t", v)
	}
	if v := iface.Config.Type; v != iftype {
		t.Errorf("ifaces.Put unmatch. iface.config.type=%s", v)
	}
}

func TestInterface_config_mtu(t *testing.T) {
	mtu := uint16(100)
	mtuStr := fmt.Sprintf("%d", mtu)
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config/mtu", mtuStr},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := iface.Config.Compare(INTERFACE_MTU_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.config.cmp=%t", v)
	}
	if v := iface.Config.Mtu; v != mtu {
		t.Errorf("ifaces.Put unmatch. iface.config.mtu=%d", v)
	}
}

func TestInterface_config_desc(t *testing.T) {
	desc := "TEST-DESC"
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/config/description", desc},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := iface.Config.Compare(OC_DESCRIPTION_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.config.cmp=%t", v)
	}
	if v := iface.Config.Desc; v != desc {
		t.Errorf("ifaces.Put unmatch. iface.config.desc=%s", v)
	}
}

func TestInterface_config_subiface_x0(t *testing.T) {
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/subinterfaces", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(SUBINTERFACES_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := len(iface.Subinterfaces); v != 0 {
		t.Errorf("ifaces.Put unmatch. iface.subifaces=%d", v)
	}
}

func TestInterface_config_subiface_x1(t *testing.T) {
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/subinterfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/subinterfaces/subinterface[index='10']", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(SUBINTERFACES_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := len(iface.Subinterfaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. iface.subifaces=%d", v)
	}
	if _, ok := iface.Subinterfaces[10]; !ok {
		t.Errorf("ifaces.Put unmatch. iface.subiface[10]=%t", ok)
	}
}

func TestInterface_config_subiface_x2(t *testing.T) {
	ifaces := makeIfaces([][2]string{
		{"/openconfig-interfaces:interfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/subinterfaces", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/subinterfaces/subinterface[index='10']", ""},
		{"/openconfig-interfaces:interfaces/interface[name='eth1']/subinterfaces/subinterface[index='11']", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("ifaces.Put unmatch. #ifaces=%d", v)
	}

	iface := ifaces["eth1"]

	if v := iface.Compare(SUBINTERFACES_KEY); !v {
		t.Errorf("ifaces.Put unmatch. iface.cmp=%t", v)
	}
	if v := len(iface.Subinterfaces); v != 2 {
		t.Errorf("ifaces.Put unmatch. iface.subifaces=%d", v)
	}
	if _, ok := iface.Subinterfaces[10]; !ok {
		t.Errorf("ifaces.Put unmatch. iface.subiface[10]=%t", ok)
	}
	if _, ok := iface.Subinterfaces[11]; !ok {
		t.Errorf("ifaces.Put unmatch. iface.subiface[11]=%t", ok)
	}
}
