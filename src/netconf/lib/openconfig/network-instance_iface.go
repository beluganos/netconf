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
	ncnet "netconf/lib/net"
	ncxml "netconf/lib/xml"
	"strconv"
)

//
// network-instances/network-instance[name]/interfaces
//
type NetworkInstanceInterfaces map[string]*NetworkInstanceInterface

func NewNetworkInstanceInterfaces() NetworkInstanceInterfaces {
	return NetworkInstanceInterfaces{}
}

func (n NetworkInstanceInterfaces) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	id, ok := nodes[0].Attrs[OC_ID_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", INTERFACE_KEY, OC_ID_KEY, nodes[0])
	}

	iface, ok := n[id]
	if !ok {
		iface = NewNetworkInstanceInterface(id)
		n[id] = iface
	}

	return iface.Put(nodes[1:], value)
}

func ProcessNetworkInstanceInterfaces(p NetworkInstanceInterfaceProcessor, reverse bool, name string, niIfaces NetworkInstanceInterfaces) error {
	for id, niIface := range niIfaces {
		if err := ProcessNetworkInstanceInterface(p, reverse, name, id, niIface); err != nil {
			return err
		}
	}
	return nil
}

func (i NetworkInstanceInterfaces) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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
// network-instances/network-instance[name]/interfaces/interface[id]
//
type NetworkInstanceInterface struct {
	nclib.SrChanges `xml:"-"`

	Id     string                          `xml:"id"`
	Config *NetworkInstanceInterfaceConfig `xml:"config"`
}

type NetworkInstanceInterfaceProcessor interface {
	networkInstanceInterfaceProcessor
	NetworkInstanceInterfaceConfigProcessor
}

type networkInstanceInterfaceProcessor interface {
	NetworkInstanceInterface(string, string, *NetworkInstanceInterface) error
}

func NewNetworkInstanceInterface(id string) *NetworkInstanceInterface {
	return &NetworkInstanceInterface{
		SrChanges: nclib.NewSrChanges(),
		Id:        id,
		Config:    NewNetworkInstanceInterfaceConfig(),
	}
}

func (n *NetworkInstanceInterface) String() string {
	return fmt.Sprintf("%s{%s='%s', %s} %s",
		INTERFACE_KEY,
		OC_ID_KEY, n.Id,
		n.Config,
		n.SrChanges,
	)
}

func (n *NetworkInstanceInterface) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		// n.Id = value // set by NewNetworkInstanceInterface

	case OC_CONFIG_KEY:
		if err := n.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	n.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceInterface(p NetworkInstanceInterfaceProcessor, reverse bool, name string, ifaceId string, niIface *NetworkInstanceInterface) error {

	idFunc := func() error {
		if niIface.GetChange(OC_ID_KEY) {
			return p.NetworkInstanceInterface(name, ifaceId, niIface)
		}
		return nil
	}

	configFunc := func() error {
		if niIface.GetChange(OC_CONFIG_KEY) {
			return ProcessNetworkInstanceInterfaceConfig(
				p.(NetworkInstanceInterfaceConfigProcessor),
				reverse,
				name,
				ifaceId,
				niIface.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, idFunc, configFunc)
}

//
// network-instances/network-instance[name]/interfaces/interface[id]/config
//
type NetworkInstanceInterfaceConfig struct {
	nclib.SrChanges `xml:"-"`

	Id           string `xml:"id"`
	Interface    string `xml:"interface"`
	Subinterface uint32 `xml:"subinterface"`
}

type NetworkInstanceInterfaceConfigProcessor interface {
	NetworkInstanceInterfaceConfig(string, string, *NetworkInstanceInterfaceConfig) error
}

func NewNetworkInstanceInterfaceConfig() *NetworkInstanceInterfaceConfig {
	return &NetworkInstanceInterfaceConfig{
		SrChanges:    nclib.NewSrChanges(),
		Id:           "",
		Interface:    "",
		Subinterface: 0,
	}
}

func (c *NetworkInstanceInterfaceConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s='%s', %s=%d} %s",
		OC_CONFIG_KEY,
		OC_ID_KEY, c.Id,
		INTERFACE_KEY, c.Interface,
		SUBINTERFACE_KEY, c.Subinterface,
		c.SrChanges,
	)
}

func (c *NetworkInstanceInterfaceConfig) IFName() string {
	return ncnet.NewIFName(c.Interface, c.Subinterface)
}

func (c *NetworkInstanceInterfaceConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		c.Id = value

	case INTERFACE_KEY:
		c.Interface = value

	case SUBINTERFACE_KEY:
		subif, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		c.Subinterface = uint32(subif)
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceInterfaceConfig(p NetworkInstanceInterfaceConfigProcessor, reverse bool, name string, ifaceId string, config *NetworkInstanceInterfaceConfig) error {

	configFunc := func() error {
		return p.NetworkInstanceInterfaceConfig(name, ifaceId, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
