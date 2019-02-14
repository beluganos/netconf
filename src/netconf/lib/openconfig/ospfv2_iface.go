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

type Ospfv2Interfaces map[string]*Ospfv2Interface

func NewOspfv2Interfaces() Ospfv2Interfaces {
	return Ospfv2Interfaces{}
}

func (p Ospfv2Interfaces) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	id, ok := nodes[0].Attrs[OC_ID_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", INTERFACE_KEY, OC_ID_KEY, nodes[0])
	}

	iface, ok := p[id]
	if !ok {
		iface = NewOspfv2Interface(id)
		p[id] = iface
	}

	return iface.Put(nodes[1:], value)
}

func ProcessOspfv2Interfaces(p Ospfv2InterfaceProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaces Ospfv2Interfaces) error {
	for ifaceId, iface := range ifaces {
		if err := ProcessOspfv2Interface(p, reverse, name, key, areaId, ifaceId, iface); err != nil {
			return err
		}
	}
	return nil
}

func (p Ospfv2Interfaces) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = INTERFACES_KEY
	e.EncodeToken(start)

	for _, iface := range p {
		err := e.EncodeElement(iface, xml.StartElement{Name: xml.Name{Local: INTERFACE_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type Ospfv2Interface struct {
	nclib.SrChanges `xml:"-"`

	Id           string                 `xml:"id"`
	Config       *Ospfv2InterfaceConfig `xml:"config"`
	InterfaceRef *InterfaceRef          `xml:"interface-ref"`
	Timers       *Ospfv2InterfaceTimers `xml:"timers"`
}

type Ospfv2InterfaceProcessor interface {
	ospfv2InterfaceProcessor
	Ospfv2InterfaceConfigProcessor
	Ospfv2InterfaceRefProcessor
	Ospfv2InterfaceTimersProcessor
}

type ospfv2InterfaceProcessor interface {
	Ospfv2Interface(string, *NetworkInstanceProtocolKey, string, string, *Ospfv2Interface) error
}

func NewOspfv2Interface(id string) *Ospfv2Interface {
	return &Ospfv2Interface{
		SrChanges:    nclib.NewSrChanges(),
		Id:           id,
		Config:       NewOspfv2InterfaceConfig(),
		InterfaceRef: NewInterfaceRef(),
		Timers:       NewOspfv2InterfaceTimers(),
	}
}

func (o *Ospfv2Interface) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		// o.Id = value // set by NewOspfv2Interface

	case OC_CONFIG_KEY:
		if err := o.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case INTERFACE_REF_KEY:
		if err := o.InterfaceRef.Put(nodes[1:], value); err != nil {
			return nil
		}

	case OSPFV2_TIMERS_KEY:
		if err := o.Timers.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv2Interface(p Ospfv2InterfaceProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, iface *Ospfv2Interface) error {
	ifaceFunc := func() error {
		if iface.GetChange(OC_ID_KEY) {
			return p.Ospfv2Interface(name, key, areaId, ifaceId, iface)
		}
		return nil
	}

	configFunc := func() error {
		if iface.GetChange(OC_CONFIG_KEY) {
			return ProcessOspfv2InterfaceConfig(
				p.(Ospfv2InterfaceConfigProcessor),
				reverse,
				name,
				key,
				areaId,
				ifaceId,
				iface.Config,
			)
		}
		return nil
	}

	ifrefFunc := func() error {
		if iface.GetChange(INTERFACE_REF_KEY) {
			return ProcessOspfv2InterfaceRef(
				p.(Ospfv2InterfaceRefProcessor),
				reverse,
				name,
				key,
				areaId,
				ifaceId,
				iface.InterfaceRef,
			)
		}
		return nil
	}

	timersFunc := func() error {
		if iface.GetChange(OSPFV2_TIMERS_KEY) {
			return ProcessOspfv2InterfaceTimers(
				p.(Ospfv2InterfaceTimersProcessor),
				reverse,
				name,
				key,
				areaId,
				ifaceId,
				iface.Timers,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, ifaceFunc, configFunc, ifrefFunc, timersFunc)
}

type Ospfv2InterfaceConfig struct {
	nclib.SrChanges `xml:"-"`

	Id          string          `xml:"id"`
	Metric      uint16          `xml:"metric"`
	Passive     bool            `xml:"passive"`
	NetworkType OspfNetworkType `xml:"network-type"`
	Priority    uint8           `xml:"priority"`
}

type Ospfv2InterfaceConfigProcessor interface {
	Ospfv2InterfaceConfig(string, *NetworkInstanceProtocolKey, string, string, *Ospfv2InterfaceConfig) error
}

func NewOspfv2InterfaceConfig() *Ospfv2InterfaceConfig {
	return &Ospfv2InterfaceConfig{
		SrChanges:   nclib.NewSrChanges(),
		Id:          "",
		Metric:      0,
		Passive:     false,
		NetworkType: OSPF_BROADCAST_NETWORK,
		Priority:    1,
	}
}

func (c *Ospfv2InterfaceConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		c.Id = value

	case OSPFV2_METRIC_KEY:
		metric, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}
		c.Metric = uint16(metric)

	case OSPFV2_PASSIVE_KEY:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		c.Passive = b

	case OSPFV2_NETWORK_TYPE_KEY:
		n, err := ParseOspfNetworkType(value)
		if err != nil {
			return err
		}
		c.NetworkType = n

	case OSPFV2_PRIORITY_KEY:
		priority, err := strconv.ParseUint(value, 0, 8)
		if err != nil {
			return err
		}
		c.Priority = uint8(priority)
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv2InterfaceConfig(p Ospfv2InterfaceConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, config *Ospfv2InterfaceConfig) error {
	configFunc := func() error {
		return p.Ospfv2InterfaceConfig(name, key, areaId, ifaceId, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

type Ospfv2InterfaceTimers struct {
	nclib.SrChanges `xml:"-"`

	HelloInterval uint32 `xml:"hello-interval"`
	DeadInterval  uint32 `xml:"dead-interval"`
}

type Ospfv2InterfaceTimersProcessor interface {
	Ospfv2InterfaceTimers(string, *NetworkInstanceProtocolKey, string, string, *Ospfv2InterfaceTimers) error
}

func NewOspfv2InterfaceTimers() *Ospfv2InterfaceTimers {
	return &Ospfv2InterfaceTimers{
		SrChanges:     nclib.NewSrChanges(),
		HelloInterval: 0,
		DeadInterval:  0,
	}
}

func (o *Ospfv2InterfaceTimers) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OSPFV2_DEAD_INTERVAL_KEY:
		t, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		o.DeadInterval = uint32(t)

	case OSPFV2_HELLO_INTERVAL_KEY:
		t, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		o.HelloInterval = uint32(t)
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv2InterfaceTimers(p Ospfv2InterfaceTimersProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, timers *Ospfv2InterfaceTimers) error {
	timersFunc := func() error {
		return p.Ospfv2InterfaceTimers(name, key, areaId, ifaceId, timers)
	}

	return nclib.CallFunctions(reverse, timersFunc)
}

type Ospfv2InterfaceRefProcessor interface {
	Ospfv2InterfaceRefConfig(string, *NetworkInstanceProtocolKey, string, string, *InterfaceRefConfig) error
}

func ProcessOspfv2InterfaceRef(p Ospfv2InterfaceRefProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, ifaceRef *InterfaceRef) error {
	refFunc := func() error {
		if ifaceRef.GetChange(OC_CONFIG_KEY) {
			return p.Ospfv2InterfaceRefConfig(name, key, areaId, ifaceId, ifaceRef.Config)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, refFunc)
}
