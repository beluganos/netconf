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
	"encoding/xml"
	"fmt"
	"netconf/lib"
	"netconf/lib/xml"
	"strconv"
)

//
// mpls/signaling-protocols/ldp/global/discovery/interfaces/interface*
//
type MplsLdpInterfaces map[string]*MplsLdpInterface

func NewMplsLdpInterfaces() MplsLdpInterfaces {
	return MplsLdpInterfaces{}
}

func (m MplsLdpInterfaces) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ifid, ok := nodes[0].Attrs[INTERFACE_ID_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", INTERFACE_KEY, INTERFACE_ID_KEY, nodes[0])
	}

	iface, ok := m[ifid]
	if !ok {
		iface = NewMplsLdpInterface(ifid)
		m[ifid] = iface
	}

	return iface.Put(nodes[1:], value)
}

func ProcessMplsLdpInterfaces(p MplsLdpInterfaceProcessor, reverse bool, name string, ifaces MplsLdpInterfaces) error {
	for ifaceId, iface := range ifaces {
		if err := ProcessMplsLdpInterface(p, reverse, name, ifaceId, iface); err != nil {
			return err
		}
	}
	return nil
}

func (i MplsLdpInterfaces) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = INTERFACES_KEY
	e.EncodeToken(start)

	for _, iface := range i {
		err := e.EncodeElement(iface, xml.StartElement{Name: xml.Name{Local: INTERFACE_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// mpls/signaling-protocols/ldp/global/discovery/interfaces/interface[id]
//
type MplsLdpInterface struct {
	nclib.SrChanges `xml:"-"`

	IfaceId  string                  `xml:"interface-id"`
	IfaceRef *InterfaceRef           `xml:"interface-ref"`
	Config   *MplsLdpInterfaceConfig `xml:"config"`
}

type MplsLdpInterfaceProcessor interface {
	mplsLdpInterfaceProcessor
	MplsLdpInterfaceRefProcessor
	MplsLdpInterfaceConfigProcessor
}

type mplsLdpInterfaceProcessor interface {
	MplsLdpInterface(string, string, *MplsLdpInterface) error
}

func NewMplsLdpInterface(ifid string) *MplsLdpInterface {
	return &MplsLdpInterface{
		SrChanges: nclib.NewSrChanges(),
		IfaceId:   ifid,
		IfaceRef:  NewInterfaceRef(),
		Config:    NewMplsLdpInterfaceConfig(),
	}
}

func (m *MplsLdpInterface) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s} %s",
		INTERFACE_KEY,
		INTERFACE_ID_KEY, m.IfaceId,
		m.IfaceRef,
		m.Config,
		m.SrChanges,
	)
}

func (m *MplsLdpInterface) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case INTERFACE_ID_KEY:
		// m.IfaceId = value // set by NewMplsLdpInterface

	case INTERFACE_REF_KEY:
		if err := m.IfaceRef.Put(nodes[1:], value); err != nil {
			return err
		}

	case OC_CONFIG_KEY:
		if err := m.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpInterface(p MplsLdpInterfaceProcessor, reverse bool, name string, ifaceId string, iface *MplsLdpInterface) error {

	idFunc := func() error {
		if iface.GetChange(INTERFACE_ID_KEY) {
			return p.MplsLdpInterface(name, ifaceId, iface)
		}
		return nil
	}

	configFunc := func() error {
		if iface.GetChange(OC_CONFIG_KEY) {
			return ProcessMplsLdpInterfaceConfig(
				p.(MplsLdpInterfaceConfigProcessor),
				reverse,
				name,
				ifaceId,
				iface.Config,
			)
		}
		return nil
	}

	refFunc := func() error {
		if iface.GetChange(INTERFACE_REF_KEY) {
			return ProcessMplsLdpInterfaceRef(
				p.(MplsLdpInterfaceRefProcessor),
				reverse,
				name,
				ifaceId,
				iface.IfaceRef,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, idFunc, configFunc, refFunc)
}

//
// mpls/signaling-protocols/ldp/address-family/ipvX/interfaces/interface[interface-id]/config
//
type MplsLdpInterfaceConfig struct {
	*InterfaceRefConfig `xml:"-"`

	IfaceId       string `xml:"interface-id"`
	HelloHoldTime uint16 `xml:"hello-holdtime"`
	HelloInterval uint16 `xml:"hello-interval"`
}

type MplsLdpInterfaceConfigProcessor interface {
	MplsLdpInterfaceConfig(string, string, *MplsLdpInterfaceConfig) error
}

func NewMplsLdpInterfaceConfig() *MplsLdpInterfaceConfig {
	return &MplsLdpInterfaceConfig{
		InterfaceRefConfig: NewInterfaceRefConfig(),
		IfaceId:            "",
		HelloHoldTime:      0,
		HelloInterval:      0,
	}
}

func (m *MplsLdpInterfaceConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s}",
		OC_CONFIG_KEY,
		INTERFACE_ID_KEY, m.IfaceId,
		m.InterfaceRefConfig,
	)
}

func (m *MplsLdpInterfaceConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case INTERFACE_ID_KEY:
		m.IfaceId = value

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

func ProcessMplsLdpInterfaceConfig(p MplsLdpInterfaceConfigProcessor, reverse bool, name string, ifaceId string, config *MplsLdpInterfaceConfig) error {
	configFunc := func() error {
		return p.MplsLdpInterfaceConfig(name, ifaceId, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// mpls/signaling-protocols/ldp/address-family/ipvX/interfaces/interface[interface-id]/interface-ref
//
type MplsLdpInterfaceRefProcessor interface {
	MplsLdpInterfaceRefConfig(string, string, *InterfaceRefConfig) error
}

func ProcessMplsLdpInterfaceRef(p MplsLdpInterfaceRefProcessor, reverse bool, name string, ifaceId string, ifaceRef *InterfaceRef) error {
	refFunc := func() error {
		if ifaceRef.GetChange(OC_CONFIG_KEY) {
			return p.MplsLdpInterfaceRefConfig(name, ifaceId, ifaceRef.Config)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, refFunc)
}
