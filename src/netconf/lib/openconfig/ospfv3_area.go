// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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
// ospfv3/areas
//
type Ospfv3Areas map[string]*Ospfv3Area

func NewOspfv3Areas() Ospfv3Areas {
	return Ospfv3Areas{}
}

func (o Ospfv3Areas) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ident, ok := nodes[0].Attrs[OC_IDENT_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", OSPFV3_AREA_KEY, OC_IDENT_KEY, nodes[0])
	}

	area, ok := o[ident]
	if !ok {
		area = NewOspfv3Area(ident)
		o[ident] = area
	}

	return area.Put(nodes[1:], value)
}

func ProcessOspfv3Areas(p Ospfv3AreaProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areas Ospfv3Areas) error {
	for areaId, area := range areas {
		if err := ProcessOspfv3Area(p, reverse, name, key, areaId, area); err != nil {
			return err
		}
	}
	return nil
}

func (o Ospfv3Areas) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = OSPFV3_AREAS_KEY
	e.EncodeToken(start)

	for _, area := range o {
		err := e.EncodeElement(area, xml.StartElement{Name: xml.Name{Local: OSPFV3_AREA_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// ospfv3/areas/area[identifier]
//
type Ospfv3Area struct {
	nclib.SrChanges `xml:"-"`

	Ident      string            `xml:"identifier"`
	Config     *Ospfv3AreaConfig `xml:"config"`
	Interfaces Ospfv3Interfaces  `xml:"interfaces"`
	Ranges     Ospfv3AreaRanges  `xml:"ranges"`
}

type Ospfv3AreaProcessor interface {
	ospfv3AreaProcessor
	Ospfv3AreaConfigProcessor
	Ospfv3InterfaceProcessor
	Ospfv3AreaRangeProcessor
}

type ospfv3AreaProcessor interface {
	Ospfv3Area(string, *NetworkInstanceProtocolKey, string, *Ospfv3Area) error
}

func NewOspfv3Area(ident string) *Ospfv3Area {
	return &Ospfv3Area{
		SrChanges:  nclib.NewSrChanges(),
		Ident:      ident,
		Config:     NewOspfv3AreaConfig(),
		Interfaces: NewOspfv3Interfaces(),
		Ranges:     NewOspfv3AreaRanges(),
	}
}

func (o *Ospfv3Area) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_IDENT_KEY:
		// c.Ident = value // set by NewOspfv3Area

	case OC_CONFIG_KEY:
		if err := o.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case INTERFACES_KEY:
		if err := o.Interfaces.Put(nodes[1:], value); err != nil {
			return err
		}

	case OSPFV3_RANGES_KEY:
		if err := o.Ranges.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3Area(p Ospfv3AreaProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, area *Ospfv3Area) error {

	identFunc := func() error {
		if area.GetChange(OC_IDENT_KEY) {
			return p.Ospfv3Area(name, key, areaId, area)
		}
		return nil
	}

	configFunc := func() error {
		if area.GetChange(OC_CONFIG_KEY) {
			return ProcessOspfv3AreaConfig(
				p.(Ospfv3AreaConfigProcessor),
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
			return ProcessOspfv3Interfaces(
				p.(Ospfv3InterfaceProcessor),
				reverse,
				name,
				key,
				areaId,
				area.Interfaces,
			)
		}
		return nil
	}

	rangesFunc := func() error {
		if area.GetChange(OSPFV3_RANGES_KEY) {
			return ProcessOspfv3AreaRanges(
				p.(Ospfv3AreaRangeProcessor),
				reverse,
				name,
				key,
				areaId,
				area.Ranges,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, identFunc, configFunc, ifaccesFunc, rangesFunc)
}

//
// ospfv3/areas/area[identifier]/config
//
type Ospfv3AreaConfig struct {
	nclib.SrChanges `xml:"-"`

	Ident string `xml:"identifier"`
}

type Ospfv3AreaConfigProcessor interface {
	Ospfv3AreaConfig(string, *NetworkInstanceProtocolKey, string, *Ospfv3AreaConfig) error
}

func NewOspfv3AreaConfig() *Ospfv3AreaConfig {
	return &Ospfv3AreaConfig{
		SrChanges: nclib.NewSrChanges(),
		Ident:     "",
	}
}

func (c *Ospfv3AreaConfig) Put(nodes []*ncxml.XPathNode, value string) error {
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

func ProcessOspfv3AreaConfig(p Ospfv3AreaConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, config *Ospfv3AreaConfig) error {
	configFunc := func() error {
		return p.Ospfv3AreaConfig(name, key, areaId, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
