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
	ncianalib "netconf/lib/iana"
	ncxml "netconf/lib/xml"
	"strconv"
)

const (
	INTERFACES_XMLNS          = "https://github.com/beluganos/beluganos/yang/interfaces"
	INTERFACES_MODULE         = "beluganos-interfaces"
	INTERFACES_KEY            = "interfaces"
	INTERFACE_KEY             = "interface"
	INTERFACE_TYPE_KEY        = "type"
	INTERFACE_MTU_KEY         = "mtu"
	INTERFACE_ID_KEY          = "interface-id"
	INTERFACE_REF_KEY         = "interface-ref"
	INTERFACE_ETH_MODULE      = "beluganos-if-ethernet"
	INTERFACE_ETH_KEY         = "ethernet"
	INTERFACE_ETH_MACADDR_KEY = "mac-address"
)

//
// /interfaces
//
type Interfaces map[string]*Interface

func NewInterfaces() Interfaces {
	return Interfaces{}
}

func (i Interfaces) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ifname, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", INTERFACE_KEY, OC_NAME_KEY, nodes[0])
	}

	iface, ok := i[ifname]
	if !ok {
		iface = NewInterface(ifname)
		i[ifname] = iface
	}

	return iface.Put(nodes[1:], value)
}

func ProcessInterfaces(p InterfaceProcessor, reverse bool, ifaces Interfaces) error {
	for name, iface := range ifaces {
		if err := ProcessInterface(p, reverse, name, iface); err != nil {
			return err
		}
	}
	return nil
}

func (i Interfaces) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = INTERFACES_KEY
	start.Name.Space = INTERFACES_XMLNS
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
// /interfaces/interface[name]
//
type Interface struct {
	nclib.SrChanges `xml:"-"`

	Name          string             `xml:"name"`
	Config        *InterfaceConfig   `xml:"config"`
	Ethernet      *InterfaceEthernet `xml:"ethernet"`
	Subinterfaces Subinterfaces      `xml:"subinterfaces"`
}

type interfaceProcessor interface {
	Interface(string, *Interface) error
}

type InterfaceProcessor interface {
	interfaceProcessor
	InterfaceConfigProcessor
	InterfaceEthernetProcessor
	SubinterfaceProcessor
}

func NewInterface(name string) *Interface {
	return &Interface{
		SrChanges:     nclib.NewSrChanges(),
		Name:          name,
		Config:        NewInterfaceConfig(),
		Ethernet:      NewInterfaceEthernet(),
		Subinterfaces: NewSubinterfaces(),
	}
}

func (i *Interface) String() string {
	return fmt.Sprintf("%s{%s:'%s', %s, %s} %s",
		INTERFACE_KEY,
		OC_NAME_KEY, i.Name,
		i.Config,
		i.Ethernet,
		i.SrChanges,
	)
}

func (i *Interface) SetName(name string) {
	i.Name = name
	i.SetChange(OC_NAME_KEY)
}

func (i *Interface) SetConfig(config *InterfaceConfig) {
	i.SetName(config.Name)
	i.Config = config
	i.SetChange(OC_CONFIG_KEY)
}

func (i *Interface) SetEthernet(ethernet *InterfaceEthernet) {
	i.Ethernet = ethernet
	i.SetChange(INTERFACE_ETH_KEY)
}

func (i *Interface) SetSubinterfaces(subifs Subinterfaces) {
	i.Subinterfaces = subifs
	i.SetChange(SUBINTERFACES_KEY)
}

func (i *Interface) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		// i.Name = value // set by NewInterface.

	case OC_CONFIG_KEY:
		if err := i.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case INTERFACE_ETH_KEY:
		if err := i.Ethernet.Put(nodes[1:], value); err != nil {
			return err
		}

	case SUBINTERFACES_KEY:
		if err := i.Subinterfaces.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	i.SetChange(nodes[0].Name)
	return nil
}

func ProcessInterface(p interfaceProcessor, reverse bool, name string, iface *Interface) error {
	ifName := func() error {
		if iface.GetChange(OC_NAME_KEY) {
			return p.Interface(name, iface)
		}
		return nil
	}

	ifConfig := func() error {
		if iface.GetChange(OC_CONFIG_KEY) {
			return ProcessInterfaceConfig(
				p.(InterfaceConfigProcessor),
				reverse,
				name,
				iface.Config,
			)
		}
		return nil
	}

	ifEthernet := func() error {
		if iface.GetChange(INTERFACE_ETH_KEY) {
			return ProcessInterfaceEthernet(
				p.(InterfaceEthernetConfigProcessor),
				reverse,
				name,
				iface.Ethernet,
			)
		}
		return nil
	}

	subifChange := func() error {
		if iface.GetChange(SUBINTERFACES_KEY) {
			return ProcessSubinterfaces(
				p.(SubinterfaceProcessor),
				reverse,
				name,
				iface.Subinterfaces,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, ifName, ifConfig, ifEthernet, subifChange)
}

//
// /interfaces/interface[name]/config
//
type InterfaceConfig struct {
	nclib.SrChanges `xml:"-"`

	Name    string               `xml:"name"`
	Type    ncianalib.IANAifType `xml:"type"`
	Enabled bool                 `xml:"enabled"`
	Mtu     uint16               `xml:"mtu"`
	Desc    string               `xml:"description"`
}

type InterfaceConfigProcessor interface {
	InterfaceConfig(string, *InterfaceConfig) error
}

func NewInterfaceConfig() *InterfaceConfig {
	return &InterfaceConfig{
		SrChanges: nclib.NewSrChanges(),
		Name:      "",
		Type:      "",
		Enabled:   true,
		Mtu:       0,
		Desc:      "",
	}
}

func (c *InterfaceConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s:'%s', %s:%t, %s:%d, %s:'%s'} %s",
		OC_CONFIG_KEY,
		OC_NAME_KEY, c.Name,
		INTERFACE_TYPE_KEY, c.Type,
		OC_ENABLED_KEY, c.Enabled,
		INTERFACE_MTU_KEY, c.Mtu,
		OC_DESCRIPTION_KEY, c.Desc,
		c.SrChanges,
	)
}

func (c *InterfaceConfig) SetName(name string) {
	c.Name = name
	c.SetChange(OC_NAME_KEY)
}

func (c *InterfaceConfig) SetType(t ncianalib.IANAifType) {
	c.Type = t
	c.SetChange(INTERFACE_TYPE_KEY)
}

func (c *InterfaceConfig) SetEnabled(enabled bool) {
	c.Enabled = enabled
	c.SetChange(OC_ENABLED_KEY)
}

func (c *InterfaceConfig) SetMtu(mtu uint16) {
	c.Mtu = mtu
	c.SetChange(INTERFACE_MTU_KEY)
}

func (c *InterfaceConfig) SetDesc(desc string) {
	c.Desc = desc
	c.SetChange(OC_DESCRIPTION_KEY)
}

func (c *InterfaceConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		c.Name = value

	case INTERFACE_TYPE_KEY:
		t, err := ncianalib.ParseIANAifType(value)
		if err != nil {
			return err
		}
		c.Type = t

	case OC_ENABLED_KEY:
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		c.Enabled = enabled

	case INTERFACE_MTU_KEY:
		mtu, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}
		c.Mtu = uint16(mtu)

	case OC_DESCRIPTION_KEY:
		c.Desc = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessInterfaceConfig(p InterfaceConfigProcessor, reverse bool, name string, config *InterfaceConfig) error {

	ifConfig := func() error {
		return p.InterfaceConfig(name, config)
	}

	return nclib.CallFunctions(reverse, ifConfig)
}
