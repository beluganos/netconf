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
	nclib "netconf/lib"
	ncxml "netconf/lib/xml"
)

//
// mpls/global/interface-attributes
//
type MplsInterfaceAttrs map[string]*MplsInterfaceAttr

func NewMplsInterfaceAttrs() MplsInterfaceAttrs {
	return MplsInterfaceAttrs{}
}

func (a MplsInterfaceAttrs) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ifId, ok := nodes[0].Attrs[INTERFACE_ID_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", INTERFACE_KEY, INTERFACE_ID_KEY, nodes[0])
	}

	attr, ok := a[ifId]
	if !ok {
		attr = NewMplsInterfaceAttr(ifId)
		a[ifId] = attr
	}

	return attr.Put(nodes[1:], value)
}

func ProcessMplsInterfaceAttrs(p MplsInterfaceAttrProcessor, reverse bool, name string, attrs MplsInterfaceAttrs) error {
	for ifaceId, attr := range attrs {
		if err := ProcessMplsInterfaceAttr(p, reverse, name, ifaceId, attr); err != nil {
			return err
		}
	}
	return nil
}

func (a MplsInterfaceAttrs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = MPLS_IFATTRS_KEY
	e.EncodeToken(start)

	for _, iface := range a {
		err := e.EncodeElement(iface, xml.StartElement{Name: xml.Name{Local: INTERFACE_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// mpls/global/interface-attributes/interface[interface-id]
//
type MplsInterfaceAttr struct {
	nclib.SrChanges `xml:"-"`

	IfaceId  string                   `xml:"interface-id"`
	Config   *MplsInterfaceAttrConfig `xml:"config"`
	IfaceRef *InterfaceRef            `xml:"interface-ref"`
}

type MplsInterfaceAttrProcessor interface {
	mplsInterfaceAttrProcessor
	MplsInterfaceAttrConfigProcessor
	MplsInterfaceAttrRefProcessor
}

type mplsInterfaceAttrProcessor interface {
	MplsInterfaceAttr(string, string, *MplsInterfaceAttr) error
}

func NewMplsInterfaceAttr(ifaceId string) *MplsInterfaceAttr {
	return &MplsInterfaceAttr{
		SrChanges: nclib.NewSrChanges(),
		IfaceId:   ifaceId,
		Config:    NewMplsInterfaceAttrConfig(),
		IfaceRef:  NewInterfaceRef(),
	}
}

func (a *MplsInterfaceAttr) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s} %s",
		INTERFACE_KEY,
		INTERFACE_ID_KEY, a.IfaceId,
		a.Config,
		a.IfaceRef,
		a.SrChanges,
	)
}

func (a *MplsInterfaceAttr) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case INTERFACE_ID_KEY:
		// a.IfaceId = value // set by NewMplsIfaceAttr

	case OC_CONFIG_KEY:
		if err := a.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case INTERFACE_REF_KEY:
		if err := a.IfaceRef.Put(nodes[1:], value); err != nil {
			return nil
		}
	}

	a.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsInterfaceAttr(p MplsInterfaceAttrProcessor, reverse bool, name string, ifaceId string, attr *MplsInterfaceAttr) error {
	ifaceFunc := func() error {
		if attr.GetChange(INTERFACE_ID_KEY) {
			return p.MplsInterfaceAttr(name, ifaceId, attr)
		}
		return nil
	}

	configFunc := func() error {
		if attr.GetChange(OC_CONFIG_KEY) {
			return ProcessMplsInterfaceAttrConfig(
				p.(MplsInterfaceAttrConfigProcessor),
				reverse,
				name,
				ifaceId,
				attr.Config,
			)
		}
		return nil
	}

	refFunc := func() error {
		if attr.GetChange(INTERFACE_REF_KEY) {
			return ProcessMplsInterfaceAttrRef(
				p.(MplsInterfaceAttrRefProcessor),
				reverse,
				name,
				ifaceId,
				attr.IfaceRef,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, ifaceFunc, configFunc, refFunc)
}

//
// mpls/global/interface-attributes/interface[interface-id]/config
//
type MplsInterfaceAttrConfig struct {
	nclib.SrChanges `xml:"-"`

	IfaceId string `xml:"interface-id"`
}

type MplsInterfaceAttrConfigProcessor interface {
	MplsInterfaceAttrConfig(string, string, *MplsInterfaceAttrConfig) error
}

func NewMplsInterfaceAttrConfig() *MplsInterfaceAttrConfig {
	return &MplsInterfaceAttrConfig{
		SrChanges: nclib.NewSrChanges(),
		IfaceId:   "",
	}
}

func (c *MplsInterfaceAttrConfig) String() string {
	return fmt.Sprintf("%s{%s='%s'} %s",
		OC_CONFIG_KEY,
		INTERFACE_ID_KEY, c.IfaceId,
		c.SrChanges,
	)
}

func (c *MplsInterfaceAttrConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case INTERFACE_ID_KEY:
		c.IfaceId = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsInterfaceAttrConfig(p MplsInterfaceAttrConfigProcessor, reverse bool, name string, ifaceId string, config *MplsInterfaceAttrConfig) error {

	configFunc := func() error {
		return p.MplsInterfaceAttrConfig(name, ifaceId, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// mpls/global/interface-attributes/interface[interface-id]/interface-ref
//
type MplsInterfaceAttrRefProcessor interface {
	MplsInterfaceAttrRefConfig(string, string, *InterfaceRefConfig) error
}

func ProcessMplsInterfaceAttrRef(p MplsInterfaceAttrRefProcessor, reverse bool, name string, ifaceId string, ifaceRef *InterfaceRef) error {

	refFunc := func() error {
		if ifaceRef.GetChange(OC_CONFIG_KEY) {
			return p.MplsInterfaceAttrRefConfig(name, ifaceId, ifaceRef.Config)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, refFunc)
}
