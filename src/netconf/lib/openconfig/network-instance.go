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
	"netconf/lib/net"
	"netconf/lib/xml"
)

const (
	NETWORKINSTANCES_XMLNS       = "https://github.com/beluganos/beluganos/yang/network-instanc"
	NETWORKINSTANCES_MODULE      = "beluganos-network-instance"
	NETWORKINSTANCES_KEY         = "network-instances"
	NETWORKINSTANCE_KEY          = "network-instance"
	NETWORKINSTANCE_ROUTERID_KEY = "router-id"
	NETWORKINSTANCE_RD_KEY       = "route-distinguisher"
	NETWORKINSTANCE_RT_KEY       = "route-target"
	NETWORKINSTANCE_PROTOS_KEY   = "protocols"
	NETWORKINSTANCE_PROTO_KEY    = "protocol"
	NETWORKINSTANCE_LOS_KEY      = "loopbacks"
	NETWORKINSTANCE_LO_KEY       = "loopback"
	NETWORKINSTANCE_LO_ADDRS_KEY = "addresses"
	NETWORKINSTANCE_LO_ADDR_KEY  = "address"
	NETWORKINSTANCE_LO_IP_KEY    = "ip"
	NETWORKINSTANCE_LO_PLEN_KEY  = "prefix-length"
)

//
// /network-instances
//
type NetworkInstances map[string]*NetworkInstance

func NewNetworkInstances() NetworkInstances {
	return NetworkInstances{}
}

func (n NetworkInstances) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	name, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", NETWORKINSTANCE_KEY, OC_NAME_KEY, nodes[0])
	}

	inst, ok := n[name]
	if !ok {
		inst = NewNetworkInstance(name)
		n[name] = inst
	}

	return inst.Put(nodes[1:], value)
}

func ProcessNetworkInstances(p NetworkInstanceProcessor, reverse bool, nis NetworkInstances) error {
	for name, ni := range nis {
		if err := ProcessNetworkInstance(p, reverse, name, ni); err != nil {
			return err
		}
	}
	return nil
}

func (i NetworkInstances) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = NETWORKINSTANCES_KEY
	start.Name.Space = NETWORKINSTANCES_XMLNS
	e.EncodeToken(start)

	for _, ni := range i {
		err := e.EncodeElement(ni, xml.StartElement{Name: xml.Name{Local: NETWORKINSTANCE_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// /network-instances/network-instance[name]
//
type NetworkInstance struct {
	nclib.SrChanges `xml:"-"`

	Name       string                    `xml:"name"`
	Config     *NetworkInstanceConfig    `xml:"config"`
	Loopbacks  NetworkInstanceLoopbacks  `xml:"loopbacks"`
	Interfaces NetworkInstanceInterfaces `xml:"interfaces"`
	Mpls       *Mpls                     `xml:"mpls"`
	Protocols  NetworkInstanceProtocols  `xml:"protocols"`
}

type NetworkInstanceProcessor interface {
	networkInstanceProcessor
	NetworkInstanceConfigProcessor
	NetworkInstanceLoopbackProcessor
	NetworkInstanceInterfaceProcessor
	MplsProcessor
	NetworkInstanceProtocolProcessor
}

type networkInstanceProcessor interface {
	NetworkInstance(string, *NetworkInstance) error
}

func NewNetworkInstance(name string) *NetworkInstance {
	return &NetworkInstance{
		SrChanges:  nclib.NewSrChanges(),
		Name:       name,
		Config:     NewNetworkInstanceConfig(),
		Loopbacks:  NewNetworkInstanceLoopbacks(),
		Interfaces: NewNetworkInstanceInterfaces(),
		Mpls:       NewMpls(),
		Protocols:  NewNetworkInstanceProtocols(),
	}
}

func (n *NetworkInstance) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s, %s, %s, %v} %s",
		NETWORKINSTANCE_KEY,
		OC_NAME_KEY, n.Name,
		n.Config,
		n.Loopbacks,
		n.Interfaces,
		n.Mpls,
		n.Protocols,
		n.SrChanges,
	)
}

func (n *NetworkInstance) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		// n.Name = value // set by NewNetworkInstance

	case OC_CONFIG_KEY:
		if err := n.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	case NETWORKINSTANCE_LOS_KEY:
		if err := n.Loopbacks.Put(nodes[1:], value); err != nil {
			return err
		}
	case INTERFACES_KEY:
		if err := n.Interfaces.Put(nodes[1:], value); err != nil {
			return err
		}

	case MPLS_KEY:
		if err := n.Mpls.Put(nodes[1:], value); err != nil {
			return err
		}

	case NETWORKINSTANCE_PROTOS_KEY:
		if err := n.Protocols.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	n.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstance(p NetworkInstanceProcessor, reverse bool, name string, ni *NetworkInstance) error {

	nameFunc := func() error {
		if ni.GetChange(OC_NAME_KEY) {
			return p.NetworkInstance(name, ni)
		}
		return nil
	}

	configFunc := func() error {
		if ni.GetChange(OC_CONFIG_KEY) {
			return ProcessNetworkInstanceConfig(
				p.(NetworkInstanceConfigProcessor),
				reverse,
				name,
				ni.Config,
			)
		}
		return nil
	}

	losFunc := func() error {
		if ni.GetChange(NETWORKINSTANCE_LOS_KEY) {
			return ProcessNetworkInstanceLoopbacks(
				p.(NetworkInstanceLoopbackProcessor),
				reverse,
				name,
				ni.Loopbacks,
			)
		}
		return nil
	}

	ifsFunc := func() error {
		if ni.GetChange(INTERFACES_KEY) {
			return ProcessNetworkInstanceInterfaces(
				p.(NetworkInstanceInterfaceProcessor),
				reverse,
				name,
				ni.Interfaces,
			)
		}
		return nil
	}

	mplsFunc := func() error {
		if ni.GetChange(MPLS_KEY) {
			return ProcessMpls(
				p.(MplsProcessor),
				reverse,
				name,
				ni.Mpls,
			)
		}
		return nil
	}

	protosFunc := func() error {
		if ni.GetChange(NETWORKINSTANCE_PROTOS_KEY) {
			return ProcessNetworkInstanceProtocols(
				p.(NetworkInstanceProtocolProcessor),
				reverse,
				name,
				ni.Protocols,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, nameFunc, configFunc, losFunc, ifsFunc, mplsFunc, protosFunc)
}

//
// /network-instances/network-instance[name]/config
//
type NetworkInstanceConfig struct {
	nclib.SrChanges `xml:"-"`

	Name     string                   `xml:"name"`
	Type     NetworkInstanceType      `xml:"type"`
	Desc     string                   `xml:"description"`
	RouterId *ncnet.RouterId          `xml:"router-id"`
	RD       ncnet.RouteDistinguisher `xml:"route-distinguisher"`
	RT       ncnet.RouteDistinguisher `xml:"route-target"`
}

type NetworkInstanceConfigProcessor interface {
	NetworkInstanceConfig(string, *NetworkInstanceConfig) error
}

func NewNetworkInstanceConfig() *NetworkInstanceConfig {
	return &NetworkInstanceConfig{
		SrChanges: nclib.NewSrChanges(),
		Name:      "",
		Type:      NETWORK_INSTANCE_DEFAULT,
		Desc:      "",
		RouterId:  nil,
		RD:        ncnet.RouteDistinguisherNone{},
		RT:        ncnet.RouteDistinguisherNone{},
	}
}

func (c *NetworkInstanceConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s=%s, %s=%s, %s='%s', %s='%s', %s='%s'} %s",
		OC_CONFIG_KEY,
		OC_NAME_KEY, c.Name,
		OC_TYPE_KEY, c.Type,
		OC_DESCRIPTION_KEY, c.Desc,
		NETWORKINSTANCE_ROUTERID_KEY, c.RouterId,
		NETWORKINSTANCE_RD_KEY, c.RD,
		NETWORKINSTANCE_RT_KEY, c.RT,
		c.SrChanges,
	)
}

func (c *NetworkInstanceConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		c.Name = value

	case OC_TYPE_KEY:
		t, err := ParseNetworkInstanceType(value)
		if err != nil {
			return err
		}
		c.Type = t

	case OC_DESCRIPTION_KEY:
		c.Desc = value

	case NETWORKINSTANCE_ROUTERID_KEY:
		id, err := ncnet.ParseRouterId(value)
		if err != nil {
			return err
		}
		c.RouterId = id

	case NETWORKINSTANCE_RD_KEY:
		rd, err := ncnet.ParseRouteDistinguisher(value)
		if err != nil {
			return err
		}
		c.RD = rd

	case NETWORKINSTANCE_RT_KEY:
		rt, err := ncnet.ParseRouteDistinguisher(value)
		if err != nil {
			return err
		}
		c.RT = rt
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceConfig(p NetworkInstanceConfigProcessor, reverse bool, name string, config *NetworkInstanceConfig) error {
	configFunc := func() error {
		return p.NetworkInstanceConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
