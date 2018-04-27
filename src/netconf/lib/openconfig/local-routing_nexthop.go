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
	"net"
	"netconf/lib"
	"netconf/lib/xml"
)

//
// static-routes/static[prefix]/next-hops
//
type StaticRouteNexthops map[string]*StaticRouteNexthop

func NewStaticRouteNexthops() StaticRouteNexthops {
	return StaticRouteNexthops{}
}

func (s StaticRouteNexthops) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	index, ok := nodes[0].Attrs[OC_INDEX_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", STATICROUTE_NEXTHOP_KEY, OC_INDEX_KEY, nodes[0])
	}

	nh, ok := s[index]
	if !ok {
		nh = NewStaticRouteNexthop(index)
		s[index] = nh
	}

	return nh.Put(nodes[1:], value)
}

func ProcessStaticRouteNexthops(p StaticRouteNexthopProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, rtkey *StaticRouteKey, nexthops StaticRouteNexthops) error {
	for index, nexthop := range nexthops {
		if err := ProcessStaticRouteNexthop(p, reverse, name, key, rtkey, index, nexthop); err != nil {
			return err
		}
	}
	return nil
}

func (s StaticRouteNexthops) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = STATICROUTE_NEXTHOPS_KEY
	e.EncodeToken(start)

	for _, nh := range s {
		err := e.EncodeElement(nh, xml.StartElement{Name: xml.Name{Local: STATICROUTE_NEXTHOP_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// static-routes/static[prefix]/next-hops/next-hop[index]
//
type StaticRouteNexthop struct {
	nclib.SrChanges `xml:"-"`

	Index    string                    `xml:"index"`
	Config   *StaticRouteNexthopConfig `xml:"config"`
	IfaceRef *InterfaceRef             `xml:"interface-ref"`
}

type StaticRouteNexthopProcessor interface {
	staticRouteNexthopProcessor
	StaticRouteNexthopConfigProcessor
	StaticRouteNexthopIfaceRefProcessor
}

type staticRouteNexthopProcessor interface {
	StaticRouteNexthop(string, *NetworkInstanceProtocolKey, *StaticRouteKey, string, *StaticRouteNexthop) error
}

func NewStaticRouteNexthop(index string) *StaticRouteNexthop {
	return &StaticRouteNexthop{
		SrChanges: nclib.NewSrChanges(),
		Index:     index,
		Config:    NewStaticRouteNexthopConfig(),
		IfaceRef:  NewInterfaceRef(),
	}
}

func (s *StaticRouteNexthop) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s} %s",
		STATICROUTE_NEXTHOP_KEY,
		OC_INDEX_KEY, s.Index,
		s.Config,
		s.IfaceRef,
		s.SrChanges,
	)
}

func (s *StaticRouteNexthop) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_INDEX_KEY:
		// s.Index = value // set by NewStaticRouteNexthop

	case OC_CONFIG_KEY:
		if err := s.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	case INTERFACE_REF_KEY:
		if err := s.IfaceRef.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	s.SetChanges(nodes[0].Name)
	return nil
}

func ProcessStaticRouteNexthop(p StaticRouteNexthopProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, rtkey *StaticRouteKey, index string, nexthop *StaticRouteNexthop) error {
	nhFunc := func() error {
		if nexthop.GetChange(OC_INDEX_KEY) {
			return p.StaticRouteNexthop(name, key, rtkey, index, nexthop)
		}
		return nil
	}

	configFunc := func() error {
		if nexthop.GetChange(OC_CONFIG_KEY) {
			return ProcessStaticRouteNexthopConfig(
				p.(StaticRouteNexthopConfigProcessor),
				reverse,
				name,
				key,
				rtkey,
				index,
				nexthop.Config,
			)
		}
		return nil
	}

	refFunc := func() error {
		if nexthop.GetChange(INTERFACE_REF_KEY) {
			return ProcessStaticRouteNexthopIfaceRef(
				p.(StaticRouteNexthopIfaceRefProcessor),
				reverse,
				name,
				key,
				rtkey,
				index,
				nexthop.IfaceRef,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, nhFunc, configFunc, refFunc)
}

//
// static-routes/static[prefix]/next-hops/next-hop[index]/config
//
type StaticRouteNexthopConfig struct {
	nclib.SrChanges `xml:"-"`

	Index   string `xml:"index"`
	Nexthop string `xml:"nexthop"`
}

type StaticRouteNexthopConfigProcessor interface {
	StaticRouteNexthopConfig(string, *NetworkInstanceProtocolKey, *StaticRouteKey, string, *StaticRouteNexthopConfig) error
}

func NewStaticRouteNexthopConfig() *StaticRouteNexthopConfig {
	return &StaticRouteNexthopConfig{
		SrChanges: nclib.NewSrChanges(),
		Index:     "",
		Nexthop:   "",
	}
}

func (c *StaticRouteNexthopConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s='%s'} %s",
		OC_CONFIG_KEY,
		OC_INDEX_KEY, c.Index,
		STATICROUTE_NEXTHOP_KEY, c.Nexthop,
		c.SrChanges,
	)
}

func (c *StaticRouteNexthopConfig) GetNexthop() (net.IP, LocalDefinedNexthop, error) {
	return ParseLocalDefinedNexthops(c.Nexthop)
}

func (c *StaticRouteNexthopConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_INDEX_KEY:
		if len(value) == 0 {
			return fmt.Errorf("Invalid %s. '%s'", OC_INDEX_KEY, value)
		}
		c.Index = value

	case STATICROUTE_NEXTHOP_KEY:
		if _, _, err := ParseLocalDefinedNexthops(value); err != nil {
			return err
		}
		c.Nexthop = value
	}

	c.SetChanges(nodes[0].Name)
	return nil
}

func ProcessStaticRouteNexthopConfig(p StaticRouteNexthopConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, rtkey *StaticRouteKey, index string, config *StaticRouteNexthopConfig) error {
	configFunc := func() error {
		return p.StaticRouteNexthopConfig(name, key, rtkey, index, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

type StaticRouteNexthopIfaceRefProcessor interface {
	StaticRouteNexthopIfaceRefConfig(string, *NetworkInstanceProtocolKey, *StaticRouteKey, string, *InterfaceRefConfig) error
}

func ProcessStaticRouteNexthopIfaceRef(p StaticRouteNexthopIfaceRefProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, rtkey *StaticRouteKey, index string, ifaceRef *InterfaceRef) error {
	refFunc := func() error {
		if ifaceRef.GetChange(OC_CONFIG_KEY) {
			return p.StaticRouteNexthopIfaceRefConfig(name, key, rtkey, index, ifaceRef.Config)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, refFunc)
}
