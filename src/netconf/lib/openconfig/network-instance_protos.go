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

type NetworkInstanceProtocolKey struct {
	Ident InstallProtocolType
	Name  string
}

func NewNetworkInstanceProtocolKey(ident InstallProtocolType, name string) *NetworkInstanceProtocolKey {
	return &NetworkInstanceProtocolKey{
		Ident: ident,
		Name:  name,
	}
}

func ParseNetworkInstanceProtocolKey(ident string, name string) (*NetworkInstanceProtocolKey, error) {
	id, err := ParseInstallProtocolType(ident)
	if err != nil {
		return nil, err
	}

	return NewNetworkInstanceProtocolKey(id, name), nil
}

func (n *NetworkInstanceProtocolKey) String() string {
	return fmt.Sprintf("%s/%s", n.Ident, n.Name)
}

type NetworkInstanceProtocols map[NetworkInstanceProtocolKey]*NetworkInstanceProtocol

func NewNetworkInstanceProtocols() NetworkInstanceProtocols {
	return NetworkInstanceProtocols{}
}

func (p NetworkInstanceProtocols) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ident, ok := nodes[0].Attrs[OC_IDENT_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", NETWORKINSTANCE_PROTO_KEY, OC_IDENT_KEY, nodes[0])
	}
	name, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", NETWORKINSTANCE_PROTO_KEY, OC_NAME_KEY, nodes[0])
	}
	key, err := ParseNetworkInstanceProtocolKey(ident, name)
	if err != nil {
		return err
	}

	proto, ok := p[*key]
	if !ok {
		proto = NewNetworkInstanceProtocol(key)
		p[*key] = proto
	}

	return proto.Put(nodes[1:], value)
}

func ProcessNetworkInstanceProtocols(p NetworkInstanceProtocolProcessor, reverse bool, name string, protos NetworkInstanceProtocols) error {
	for key, proto := range protos {
		if err := ProcessNetworkInstanceProtocol(p, reverse, name, &key, proto); err != nil {
			return err
		}
	}
	return nil
}

func (p NetworkInstanceProtocols) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = NETWORKINSTANCE_PROTOS_KEY
	e.EncodeToken(start)

	for _, proto := range p {
		err := e.EncodeElement(proto, xml.StartElement{Name: xml.Name{Local: NETWORKINSTANCE_PROTO_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type NetworkInstanceProtocol struct {
	nclib.SrChanges `xml:"-"`

	// Key          *NetworkInstanceProtocolKey
	Ident        InstallProtocolType            `xml:"identifier"`
	Name         string                         `xml:"name"`
	Config       *NetworkInstanceProtocolConfig `xml:"config"`
	StaticRoutes StaticRoutes                   `xml:"-"`
	Ospfv2       *Ospfv2                        `xml:"-"`
	Bgp          *Bgp                           `xml:"-"`
}

type NetworkInstanceProtocolProcessor interface {
	networkInstanceProtocolProcessor
	NetworkInstanceProtocolConfigProcessor
	StaticRouteProcessor
	Ospfv2Processor
	BgpProcessor
}

type networkInstanceProtocolProcessor interface {
	NetworkInstanceProtocol(string, *NetworkInstanceProtocolKey, *NetworkInstanceProtocol) error
}

func NewNetworkInstanceProtocol(key *NetworkInstanceProtocolKey) *NetworkInstanceProtocol {
	return &NetworkInstanceProtocol{
		SrChanges:    nclib.NewSrChanges(),
		Ident:        key.Ident,
		Name:         key.Name,
		Config:       NewNetworkInstanceProtocolConfig(),
		StaticRoutes: NewStaticRoutes(),
		Ospfv2:       NewOspfv2(),
		Bgp:          NewBgp(),
	}
}

func (p *NetworkInstanceProtocol) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_IDENT_KEY:
		// p.Ident = value // set by NewProtocol.

	case OC_NAME_KEY:
		// p.Name = value // NewProtocol

	case OC_CONFIG_KEY:
		if err := p.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case STATICROUTES_KEY:
		if err := p.StaticRoutes.Put(nodes[1:], value); err != nil {
			return err
		}

	case OSPFV2_KEY:
		if err := p.Ospfv2.Put(nodes[1:], value); err != nil {
			return err
		}

	case BGP_KEY:
		if err := p.Bgp.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceProtocol(p NetworkInstanceProtocolProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, proto *NetworkInstanceProtocol) error {
	niProtoFunc := func() error {
		if proto.GetChanges(OC_IDENT_KEY, OC_NAME_KEY) {
			return p.NetworkInstanceProtocol(name, key, proto)
		}
		return nil
	}

	configFunc := func() error {
		if proto.GetChange(OC_CONFIG_KEY) {
			return ProcessNetworkInstanceProtocolConfig(
				p.(NetworkInstanceProtocolConfigProcessor),
				reverse,
				name,
				key,
				proto.Config,
			)
		}
		return nil
	}

	staticFunc := func() error {
		if proto.GetChange(STATICROUTES_KEY) && key.Ident == INSTALL_PROTOCOL_STATIC {
			return ProcessStaticRoutes(
				p.(StaticRouteProcessor),
				reverse,
				name,
				key,
				proto.StaticRoutes,
			)
		}
		return nil
	}

	ospfFunc := func() error {
		if proto.GetChange(OSPFV2_KEY) && key.Ident == INSTALL_PROTOCOL_OSPF {
			return ProcessOspfv2(
				p.(Ospfv2Processor),
				reverse,
				name,
				key,
				proto.Ospfv2,
			)
		}
		return nil
	}

	bgpFunc := func() error {
		if proto.GetChange(BGP_KEY) && key.Ident == INSTALL_PROTOCOL_BGP {
			return ProcessBgp(
				p.(BgpProcessor),
				reverse,
				name,
				key,
				proto.Bgp,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, niProtoFunc, configFunc, staticFunc, ospfFunc, bgpFunc)
}

type NetworkInstanceProtocolConfig struct {
	nclib.SrChanges `xml:"-"`

	Ident InstallProtocolType `xml:"identifier"`
	Name  string              `xml:"name"`
}

type NetworkInstanceProtocolConfigProcessor interface {
	NetworkInstanceProtocolConfig(string, *NetworkInstanceProtocolKey, *NetworkInstanceProtocolConfig) error
}

func NewNetworkInstanceProtocolConfig() *NetworkInstanceProtocolConfig {
	return &NetworkInstanceProtocolConfig{
		SrChanges: nclib.NewSrChanges(),
		Ident:     INSTALL_PROTOCOL_TYPE,
		Name:      "",
	}
}

func (c *NetworkInstanceProtocolConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_IDENT_KEY:
		ident, err := ParseInstallProtocolType(value)
		if err != nil {
			return err
		}
		c.Ident = ident

	case OC_NAME_KEY:
		c.Name = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceProtocolConfig(p NetworkInstanceProtocolConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, config *NetworkInstanceProtocolConfig) error {
	configFunc := func() error {
		return p.NetworkInstanceProtocolConfig(name, key, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
