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
// ospfv2/areas
//
type Ospfv2Areas map[string]*Ospfv2Area

func NewOspfv2Areas() Ospfv2Areas {
	return Ospfv2Areas{}
}

func (o Ospfv2Areas) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ident, ok := nodes[0].Attrs[OC_IDENT_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", OSPFV2_AREA_KEY, OC_IDENT_KEY, nodes[0])
	}

	area, ok := o[ident]
	if !ok {
		area = NewOspfv2Area(ident)
		o[ident] = area
	}

	return area.Put(nodes[1:], value)
}

func ProcessOspfv2Areas(p Ospfv2AreaProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areas Ospfv2Areas) error {
	for areaId, area := range areas {
		if err := ProcessOspfv2Area(p, reverse, name, key, areaId, area); err != nil {
			return err
		}
	}
	return nil
}

func (o Ospfv2Areas) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = OSPFV2_AREAS_KEY
	e.EncodeToken(start)

	for _, area := range o {
		err := e.EncodeElement(area, xml.StartElement{Name: xml.Name{Local: OSPFV2_AREA_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// ospfv2/areas/area[identifier]
//
type Ospfv2Area struct {
	nclib.SrChanges `xml:"-"`

	Ident      string            `xml:"identifier"`
	Config     *Ospfv2AreaConfig `xml:"config"`
	Interfaces Ospfv2Interfaces  `xml:"interfaces"`
}

type Ospfv2AreaProcessor interface {
	ospfv2AreaProcessor
	Ospfv2AreaConfigProcessor
	Ospfv2InterfaceProcessor
}

type ospfv2AreaProcessor interface {
	Ospfv2Area(string, *NetworkInstanceProtocolKey, string, *Ospfv2Area) error
}

func NewOspfv2Area(ident string) *Ospfv2Area {
	return &Ospfv2Area{
		SrChanges:  nclib.NewSrChanges(),
		Ident:      ident,
		Config:     NewOspfv2AreaConfig(),
		Interfaces: NewOspfv2Interfaces(),
	}
}

func (o *Ospfv2Area) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_IDENT_KEY:
		// c.Ident = value // set by NewOspfv2Area

	case OC_CONFIG_KEY:
		if err := o.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case INTERFACES_KEY:
		if err := o.Interfaces.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv2Area(p Ospfv2AreaProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, area *Ospfv2Area) error {

	identFunc := func() error {
		if area.GetChange(OC_IDENT_KEY) {
			return p.Ospfv2Area(name, key, areaId, area)
		}
		return nil
	}

	configFunc := func() error {
		if area.GetChange(OC_CONFIG_KEY) {
			return ProcessOspfv2AreaConfig(
				p.(Ospfv2AreaConfigProcessor),
				reverse,
				name,
				key,
				areaId,
				area.Config,
			)
		}
		return nil
	}

	ifaccesFunc := func() error {
		if area.GetChange(INTERFACES_KEY) {
			return ProcessOspfv2Interfaces(
				p.(Ospfv2InterfaceProcessor),
				reverse,
				name,
				key,
				areaId,
				area.Interfaces,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, identFunc, configFunc, ifaccesFunc)
}

//
// ospfv2/areas/area[identifier]/config
//
type Ospfv2AreaConfig struct {
	nclib.SrChanges `xml:"-"`

	Ident string `xml:"identifier"`
}

type Ospfv2AreaConfigProcessor interface {
	Ospfv2AreaConfig(string, *NetworkInstanceProtocolKey, string, *Ospfv2AreaConfig) error
}

func NewOspfv2AreaConfig() *Ospfv2AreaConfig {
	return &Ospfv2AreaConfig{
		SrChanges: nclib.NewSrChanges(),
		Ident:     "",
	}
}

func (c *Ospfv2AreaConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_IDENT_KEY:
		c.Ident = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv2AreaConfig(p Ospfv2AreaConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, config *Ospfv2AreaConfig) error {
	configFunc := func() error {
		return p.Ospfv2AreaConfig(name, key, areaId, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
