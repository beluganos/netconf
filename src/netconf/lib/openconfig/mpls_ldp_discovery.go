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
	"netconf/lib"
	"netconf/lib/xml"
	"strconv"
)

//
// mpls/signaling-protocols/ldp/global/discovery
//
type MplsLdpDiscovery struct {
	nclib.SrChanges `xml:"-"`

	Interfaces *MplsLdpDiscovInterfaces `xml:"interfaces"`
}

type MplsLdpDiscoveryProcessor interface {
	MplsLdpDiscovInterfacesProcessor
}

func NewMplsLdpDiscovery() *MplsLdpDiscovery {
	return &MplsLdpDiscovery{
		SrChanges:  nclib.NewSrChanges(),
		Interfaces: NewMplsLdpDiscovInterfaces(),
	}
}

func (m *MplsLdpDiscovery) String() string {
	return fmt.Sprintf("%s{%v} %s",
		MPLS_LDP_DISCOVERY_KEY,
		m.Interfaces,
		m.SrChanges,
	)
}

func (m *MplsLdpDiscovery) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case INTERFACES_KEY:
		if err := m.Interfaces.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpDiscovery(p MplsLdpDiscoveryProcessor, reverse bool, name string, disc *MplsLdpDiscovery) error {
	ifacesFunc := func() error {
		if disc.GetChange(INTERFACES_KEY) {
			return ProcessMplsLdpDiscovInterfaces(
				p.(MplsLdpDiscovInterfacesProcessor),
				reverse,
				name,
				disc.Interfaces,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, ifacesFunc)
}

//
// mpls/signaling-protocols/ldp/global/discovery/interfaces
//
type MplsLdpDiscovInterfaces struct {
	nclib.SrChanges `xml:"-"`

	Config     *MplsLdpDiscovInterfacesConfig `xml:"config"`
	Interfaces MplsLdpInterfaces              `xml:"interface"`
}

type MplsLdpDiscovInterfacesProcessor interface {
	MplsLdpDiscovInterfacesConfigProcessor
	MplsLdpInterfaceProcessor
}

func NewMplsLdpDiscovInterfaces() *MplsLdpDiscovInterfaces {
	return &MplsLdpDiscovInterfaces{
		SrChanges:  nclib.NewSrChanges(),
		Config:     NewMplsLdpDiscovInterfacesConfig(),
		Interfaces: NewMplsLdpInterfaces(),
	}
}

func (m *MplsLdpDiscovInterfaces) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := m.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case INTERFACE_KEY:
		if err := m.Interfaces.Put(nodes, value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpDiscovInterfaces(p MplsLdpDiscovInterfacesProcessor, reverse bool, name string, ifaces *MplsLdpDiscovInterfaces) error {
	configFunc := func() error {
		if ifaces.GetChange(OC_CONFIG_KEY) {
			return ProcessMplsLdpDiscovInterfacesConfig(
				p.(MplsLdpDiscovInterfacesConfigProcessor),
				reverse,
				name,
				ifaces.Config,
			)
		}
		return nil
	}

	ifacesFunc := func() error {
		return ProcessMplsLdpInterfaces(
			p.(MplsLdpInterfaceProcessor),
			reverse,
			name,
			ifaces.Interfaces,
		)
	}

	return nclib.CallFunctions(reverse, configFunc, ifacesFunc)
}

//
// mpls/signaling-protocols/ldp/global/discovery/interfaces/config
//
type MplsLdpDiscovInterfacesConfig struct {
	nclib.SrChanges `xml:"-"`

	HelloHoldTime uint16 `xml:"hello-holdtime"`
	HelloInterval uint16 `xml:"hello-interval"`
}

type MplsLdpDiscovInterfacesConfigProcessor interface {
	MplsLdpDiscovInterfacesConfig(string, *MplsLdpDiscovInterfacesConfig) error
}

func NewMplsLdpDiscovInterfacesConfig() *MplsLdpDiscovInterfacesConfig {
	return &MplsLdpDiscovInterfacesConfig{
		SrChanges:     nclib.NewSrChanges(),
		HelloHoldTime: 0,
		HelloInterval: 0,
	}
}

func (m *MplsLdpDiscovInterfacesConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_HELLO_HOLDTIME:
		v, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}
		m.HelloHoldTime = uint16(v)

	case MPLS_LDP_HELLO_INTERVAL:
		v, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}
		m.HelloInterval = uint16(v)
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpDiscovInterfacesConfig(p MplsLdpDiscovInterfacesConfigProcessor, reverse bool, name string, config *MplsLdpDiscovInterfacesConfig) error {
	configFunc := func() error {
		return p.MplsLdpDiscovInterfacesConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
