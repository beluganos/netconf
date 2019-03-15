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
	nclib "netconf/lib/sysrepo"
	"testing"
)

func makeSubinterfaces(datas [][2]string) Subinterfaces {
	siface := NewSubinterfaces()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := nclib.ParseXPath(xpath)
		if err := siface.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return siface
}

type TestSubinterfaceFunc func(string, uint32, *Subinterface) error

type TestSubinterfaceProcessor struct {
	f TestSubinterfaceFunc
}

func NewTestSubinterfaceProcessor(f TestSubinterfaceFunc) *TestSubinterfaceProcessor {
	return &TestSubinterfaceProcessor{
		f: f,
	}
}

func (s *TestSubinterfaceProcessor) Subinterface(name string, index uint32, subif *Subinterface) error {
	return s.f(name, index, subif)
}

type TestSubinterfaceConfigFunc func(string, uint32, *SubinterfaceConfig) error

type TestSubinterfaceConfigProcessor struct {
	f TestSubinterfaceConfigFunc
}

func NewTestSubinterfaceConfigProcessor(f TestSubinterfaceConfigFunc) *TestSubinterfaceConfigProcessor {
	return &TestSubinterfaceConfigProcessor{
		f: f,
	}
}

func (s *TestSubinterfaceConfigProcessor) SubinterfaceConfig(name string, index uint32, config *SubinterfaceConfig) error {
	return s.f(name, index, config)
}

//
// ProcessSubinterfaces - #subinterface=0
//
func TestSubinterfaces_processor_subif_0(t *testing.T) {
	cnt := 0
	p := NewTestSubinterfaceProcessor(func(name string, index uint32, subif *Subinterface) error {
		cnt++
		return nil
	})

	subifs := NewSubinterfaces()
	if err := ProcessSubinterfaces(p, false, "test", subifs); err != nil {
		t.Errorf("ProcessSubinterfaces error. %s", err)
	}
	if cnt != 0 {
		t.Errorf("ProcessSubinterfaces unmatch. cnt=%d", cnt)
	}
}

//
// ProcessSubinterfaces - #subinterface=1
//
func TestSubinterfaces_processor_subif_1(t *testing.T) {
	cnt := 0
	p := NewTestSubinterfaceProcessor(func(name string, index uint32, subif *Subinterface) error {
		cnt++
		return nil
	})

	subifs := NewSubinterfaces()
	subifs[0] = NewSubinterface(0)
	subifs[0].SetIndex(0)
	if err := ProcessSubinterfaces(p, false, "test", subifs); err != nil {
		t.Errorf("ProcessSubinterfaces error. %s", err)
	}
	if cnt != 1 {
		t.Errorf("ProcessSubinterfaces unmatch. cnt=%d", cnt)
	}
}

//
// ProcessSubinterfaces - error
//
func TestSubinterfaces_processor_err(t *testing.T) {
	cnt := 0
	p := NewTestSubinterfaceProcessor(func(name string, index uint32, subif *Subinterface) error {
		cnt++
		return fmt.Errorf("")
	})

	subifs := NewSubinterfaces()
	subifs[0] = NewSubinterface(0)
	subifs[0].SetIndex(0)
	if err := ProcessSubinterfaces(p, false, "test", subifs); err == nil {
		t.Errorf("ProcessSubinterfaces must be error. %s", err)
	}
	if cnt != 1 {
		t.Errorf("ProcessSubinterfaces unmatch. cnt=%d", cnt)
	}
}

//
// ProcessSubinterfaceConfig
//
func TestProcessSubinterfaceConfig(t *testing.T) {
	cnt := 0
	p := NewTestSubinterfaceConfigProcessor(func(name string, index uint32, config *SubinterfaceConfig) error {
		cnt++
		return nil
	})

	config := NewSubinterfaceConfig()
	if err := ProcessSubinterfaceConfig(p, true, "", 0, config); err != nil {
		t.Errorf("ProcessSubinterfaceConfig error. %s", err)
	}
	if cnt != 1 {
		t.Errorf("ProcessSubinterfaceConfig unmatch. %d", cnt)
	}
}

//
// ProcessSubinterfaceConfig
//
func TestProcessSubinterfaceConfig_err(t *testing.T) {
	cnt := 0
	p := NewTestSubinterfaceConfigProcessor(func(name string, index uint32, config *SubinterfaceConfig) error {
		cnt++
		return fmt.Errorf("")
	})

	config := NewSubinterfaceConfig()
	if err := ProcessSubinterfaceConfig(p, true, "", 0, config); err == nil {
		t.Errorf("ProcessSubinterfaceConfig must be error. %s", err)
	}
	if cnt != 1 {
		t.Errorf("ProcessSubinterfaceConfig unmatch. %d", cnt)
	}
}

//
// /subinterfaces
//
func TestSubinterfaces(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 0 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}
}

//
// /subinterfaces/subinterface[index='10']
//
func TestSubinterface_empty(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces/subinterface[index='10']", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}

	iface := ifaces[10]

	if v := iface.Compare(); !v {
		t.Errorf("subifaces.Put unmatch. subiface.cmp=%t", v)
	}
}

//
// /subinterfaces/subinterface[index='NOT EXIST']
//
func TestSubinterface_index_not_exist(t *testing.T) {
	ifaces := NewSubinterfaces()
	xpath := "/subinterfaces/subinterface"
	nodes := nclib.ParseXPath(xpath)

	if err := ifaces.Put(nodes[1:], ""); err == nil {
		t.Errorf("subifaces.Put must be error.")
	}
}

//
// /subinterfaces/subinterface[index='TOO BIG UINT']
//
func TestSubinterface_index_not_uint32(t *testing.T) {
	xpath := "/subinterfaces/subinterface[index='12345678901234']"
	nodes := nclib.ParseXPath(xpath)

	ifaces := NewSubinterfaces()
	if err := ifaces.Put(nodes[1:], ""); err == nil {
		t.Errorf("subifaces.Put must be error.")
	}
}

//
// /subinterfaces/subinterface[index='10']/index = 10
// /subinterfaces/subinterface[index='10']/config
// /subinterfaces/subinterface[index='10']/beluganos-if-ip:ipv4
//
func TestSubinterface(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces/subinterface[index='10']", ""},
		{"/subinterfaces/subinterface[index='10']/index", "10"},
		{"/subinterfaces/subinterface[index='10']/config", ""},
		{"/subinterfaces/subinterface[index='10']/beluganos-if-ip:ipv4", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}

	iface := ifaces[10]

	if v := iface.Compare(OC_INDEX_KEY, OC_CONFIG_KEY, SUBINTERFACE_IPV4_KEY); !v {
		t.Errorf("subifaces.Put unmatch. subiface.cmp=%t", v)
	}
	if v := iface.Index; v != 10 {
		t.Errorf("subifaces.Put unmatch. subiface.index=%d", v)
	}
}

//
// /subinterfaces/subinterface[index='10']/index = 10
// /subinterfaces/subinterface[index='10']/config
// /subinterfaces/subinterface[index='10']/beluganos-if-ip:ipv6
//
func TestSubinterfaceV6(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces/subinterface[index='10']", ""},
		{"/subinterfaces/subinterface[index='10']/index", "10"},
		{"/subinterfaces/subinterface[index='10']/config", ""},
		{"/subinterfaces/subinterface[index='10']/beluganos-if-ip:ipv6", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}

	iface := ifaces[10]

	if v := iface.Compare(OC_INDEX_KEY, OC_CONFIG_KEY, SUBINTERFACE_IPV6_KEY); !v {
		t.Errorf("subifaces.Put unmatch. subiface.cmp=%t", v)
	}
	if v := iface.Index; v != 10 {
		t.Errorf("subifaces.Put unmatch. subiface.index=%d", v)
	}
}

//
// subinterfaces/subinterface[index='10']/UNKNOWN
//
func TestSubinterface_unknown(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces/subinterface[index='10']", ""},
		{"/subinterfaces/subinterface[index='10']/UNKNOWN", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}

	iface := ifaces[10]

	if v := iface.Compare("UNKNOWN"); !v {
		t.Errorf("subifaces.Put unmatch. subiface.cmp=%t", v)
	}
}

//
// subinterfaces/subinterface[index='10']/config
//
func TestSubinterface_config_empty(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces/subinterface[index='10']", ""},
		{"/subinterfaces/subinterface[index='10']/config", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}

	iface := ifaces[10]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("subifaces.Put unmatch. subiface.cmp=%t", v)
	}
	if v := iface.Config.Compare(); !v {
		t.Errorf("subifaces.Put unmatch. subiface.config.cmp=%t", v)
	}
}

//
// /subinterfaces/subinterface[index='10']/config/index = 10
// /subinterfaces/subinterface[index='10']/config/enabled = true
// /subinterfaces/subinterface[index='10']/config/description = 'TEST-DESC'
//
func TestSubinterface_config(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces/subinterface[index='10']", ""},
		{"/subinterfaces/subinterface[index='10']/config", ""},
		{"/subinterfaces/subinterface[index='10']/config/index", "10"},
		{"/subinterfaces/subinterface[index='10']/config/enabled", "true"},
		{"/subinterfaces/subinterface[index='10']/config/description", "TEST-DESC"},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}

	iface := ifaces[10]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("subifaces.Put unmatch. subiface.cmp=%t", v)
	}
	if v := iface.Config.Compare(OC_INDEX_KEY, OC_DESCRIPTION_KEY, OC_ENABLED_KEY); !v {
		t.Errorf("subifaces.Put unmatch. subiface.config.cmp=%t", v)
	}
	if v := iface.Config.Index; v != 10 {
		t.Errorf("subifaces.Put unmatch. subiface.config.index=%d", v)
	}
	if v := iface.Config.Enabled; v != true {
		t.Errorf("subifaces.Put unmatch. subiface.config.enabled=%t", v)
	}
	if v := iface.Config.Desc; v != "TEST-DESC" {
		t.Errorf("subifaces.Put unmatch. subiface.config.desc=%s", v)
	}
}

//
// /subinterfaces/subinterface[index='10']/config/UNKNOWN = ""
//
func TestSubinterface_config_unknown(t *testing.T) {
	ifaces := makeSubinterfaces([][2]string{
		{"/subinterfaces/subinterface[index='10']/config/UNKNOWN", ""},
	})

	t.Log(ifaces)

	if v := len(ifaces); v != 1 {
		t.Errorf("subifaces.Put unmatch. #subifaces=%d", v)
	}

	iface := ifaces[10]

	if v := iface.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("subifaces.Put unmatch. subiface.cmp=%t", v)
	}
	if v := iface.Config.Compare("UNKNOWN"); !v {
		t.Errorf("subifaces.Put unmatch. subiface.config.cmp=%t", v)
	}
}

//
// /subinterfaces/subinterface[index='10']/config/index = "TOO BIG UINT32"
//
func TestSubinterface_config_index_err(t *testing.T) {
	xpath := "/subinterfaces/subinterface[index='10']/config/index"
	nodes := nclib.ParseXPath(xpath)

	ifaces := NewSubinterfaces()
	if err := ifaces.Put(nodes[1:], "1234567890123"); err == nil {
		t.Errorf("subifaces.Put must be error.")
	}
}

//
// /subinterfaces/subinterface[index='10']/config/enabled = "BOOL"
//
func TestSubinterface_config_enabled_err(t *testing.T) {
	xpath := "/subinterfaces/subinterface[index='10']/config/enabled"
	nodes := nclib.ParseXPath(xpath)

	ifaces := NewSubinterfaces()
	if err := ifaces.Put(nodes[1:], "BOOL"); err == nil {
		t.Errorf("subifaces.Put must be error.")
	}
}
