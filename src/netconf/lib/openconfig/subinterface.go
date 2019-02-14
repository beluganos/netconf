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
	"strconv"
)

const (
	SUBINTERFACES_KEY               = "subinterfaces"
	SUBINTERFACE_KEY                = "subinterface"
	SUBINTERFACE_IP_MODULE          = "beluganos-if-ip"
	SUBINTERFACE_IPV4_KEY           = "ipv4"
	SUBINTERFACE_IPV6_KEY           = "ipv6"
	SUBINTERFACE_MTU_KEY            = "mtu"
	SUBINTERFACE_ADDRS_KEY          = "addresses"
	SUBINTERFACE_ADDR_KEY           = "address"
	SUBINTERFACE_ADDR_IP_KEY        = "ip"
	SUBINTERFACE_ADDR_PREFIXLEN_KEY = "prefix-length"
)

//
// Subinterfaces
//
type Subinterfaces map[uint32]*Subinterface

func NewSubinterfaces() Subinterfaces {
	return Subinterfaces{}
}

func (i Subinterfaces) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	indexStr, ok := nodes[0].Attrs[OC_INDEX_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", SUBINTERFACE_KEY, OC_INDEX_KEY, nodes[0])
	}

	index64, err := strconv.ParseUint(indexStr, 0, 32)
	if err != nil {
		return err
	}

	index := uint32(index64)
	subiface, ok := i[index]
	if !ok {
		subiface = NewSubinterface(index)
		i[index] = subiface
	}

	return subiface.Put(nodes[1:], value)
}

func ProcessSubinterfaces(p subinterfaceProcessor, reverse bool, name string, subifaces Subinterfaces) error {
	for index, subif := range subifaces {
		if err := ProcessSubinterface(p, reverse, name, index, subif); err != nil {
			return err
		}
	}
	return nil
}

func (i Subinterfaces) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = SUBINTERFACES_KEY
	e.EncodeToken(start)

	for _, iface := range i {
		err := e.EncodeElement(iface, xml.StartElement{Name: xml.Name{Local: SUBINTERFACE_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// Subinterface
//
type Subinterface struct {
	nclib.SrChanges `xml:"-"`

	Index  uint32              `xml:"index"`
	Config *SubinterfaceConfig `xml:"config"`
	IPv4   *SubinterfaceIPv4   `xml:"ipv4"`
}

type SubinterfaceProcessor interface {
	subinterfaceProcessor
	SubinterfaceConfigProcessor
	SubinterfaceIPv4Processor
}

type subinterfaceProcessor interface {
	Subinterface(string, uint32, *Subinterface) error
}

func NewSubinterface(index uint32) *Subinterface {
	return &Subinterface{
		SrChanges: nclib.NewSrChanges(),
		Index:     index,
		Config:    NewSubinterfaceConfig(),
		IPv4:      NewSubinterfaceIPv4(),
	}
}

func (i *Subinterface) String() string {
	return fmt.Sprintf("%s{%s=%d, %s, %s} %s",
		SUBINTERFACE_KEY,
		OC_INDEX_KEY, i.Index,
		i.Config,
		i.IPv4,
		i.SrChanges,
	)
}

func (i *Subinterface) SetIndex(index uint32) {
	i.Index = index
	i.SetChange(OC_INDEX_KEY)
}

func (i *Subinterface) SetConfig(config *SubinterfaceConfig) {
	i.SetIndex(config.Index)
	i.Config = config
	i.SetChange(OC_CONFIG_KEY)
}

func (i *Subinterface) SetIPv4(ipv4 *SubinterfaceIPv4) {
	i.IPv4 = ipv4
	i.SetChange(SUBINTERFACE_IPV4_KEY)
}

func (i *Subinterface) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_INDEX_KEY:
		// set by NewSubinterface
		//index, err := strconv.ParseUint(value, 0, 32)
		//if err != nil {
		//	return err
		//}
		// i.Index = index

	case OC_CONFIG_KEY:
		if err := i.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case SUBINTERFACE_IPV4_KEY:
		if err := i.IPv4.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	i.SetChange(nodes[0].Name)
	return nil
}

func ProcessSubinterface(p subinterfaceProcessor, reverse bool, name string, index uint32, subif *Subinterface) error {

	subifIndex := func() error {
		if subif.GetChange(OC_INDEX_KEY) {
			return p.Subinterface(name, index, subif)
		}
		return nil
	}

	subifConfig := func() error {
		if subif.GetChange(OC_CONFIG_KEY) {
			return ProcessSubinterfaceConfig(
				p.(SubinterfaceConfigProcessor),
				reverse,
				name,
				index,
				subif.Config,
			)
		}
		return nil
	}

	subifIPv4 := func() error {
		if subif.GetChange(SUBINTERFACE_IPV4_KEY) {
			return ProcessSubinterfaceIPv4(
				p.(SubinterfaceIPv4Processor),
				reverse,
				name,
				index,
				subif.IPv4,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, subifIndex, subifConfig, subifIPv4)
}

//
// subinterface/config
//
type SubinterfaceConfig struct {
	nclib.SrChanges `xml:"-"`

	Index   uint32 `xml:"index"`
	Enabled bool   `xml:"enabled"`
	Desc    string `xml:"description"`
}

type SubinterfaceConfigProcessor interface {
	SubinterfaceConfig(string, uint32, *SubinterfaceConfig) error
}

func NewSubinterfaceConfig() *SubinterfaceConfig {
	return &SubinterfaceConfig{
		SrChanges: nclib.NewSrChanges(),
		Index:     0,
		Enabled:   true,
		Desc:      "",
	}
}

func (c *SubinterfaceConfig) String() string {
	return fmt.Sprintf("%s{%s=%d, %s=%t, %s='%s'} %s",
		OC_CONFIG_KEY,
		OC_INDEX_KEY, c.Index,
		OC_ENABLED_KEY, c.Enabled,
		OC_DESCRIPTION_KEY, c.Desc,
		c.SrChanges,
	)
}

func (c *SubinterfaceConfig) SetIndex(index uint32) {
	c.Index = index
	c.SetChange(OC_INDEX_KEY)
}

func (c *SubinterfaceConfig) SetEnabled(enabled bool) {
	c.Enabled = enabled
	c.SetChange(OC_ENABLED_KEY)
}

func (c *SubinterfaceConfig) SetDesc(desc string) {
	c.Desc = desc
	c.SetChange(OC_DESCRIPTION_KEY)
}

func (c *SubinterfaceConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_INDEX_KEY:
		index, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		c.Index = uint32(index)

	case OC_ENABLED_KEY:
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		c.Enabled = enabled

	case OC_DESCRIPTION_KEY:
		c.Desc = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessSubinterfaceConfig(p SubinterfaceConfigProcessor, reverse bool, name string, index uint32, config *SubinterfaceConfig) error {

	configFunc := func() error {
		return p.SubinterfaceConfig(name, index, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
