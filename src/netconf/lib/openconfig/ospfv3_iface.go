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
	"strconv"
)

type Ospfv3Interfaces map[string]*Ospfv3Interface

func NewOspfv3Interfaces() Ospfv3Interfaces {
	return Ospfv3Interfaces{}
}

func (p Ospfv3Interfaces) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	id, ok := nodes[0].Attrs[OC_ID_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", INTERFACE_KEY, OC_ID_KEY, nodes[0])
	}

	iface, ok := p[id]
	if !ok {
		iface = NewOspfv3Interface(id)
		p[id] = iface
	}

	return iface.Put(nodes[1:], value)
}

func ProcessOspfv3Interfaces(p Ospfv3InterfaceProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaces Ospfv3Interfaces) error {
	for ifaceId, iface := range ifaces {
		if err := ProcessOspfv3Interface(p, reverse, name, key, areaId, ifaceId, iface); err != nil {
			return err
		}
	}
	return nil
}

func (p Ospfv3Interfaces) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

type Ospfv3Interface struct {
	nclib.SrChanges `xml:"-"`

	Id           string                 `xml:"id"`
	Config       *Ospfv3InterfaceConfig `xml:"config"`
	InterfaceRef *InterfaceRef          `xml:"interface-ref"`
	Timers       *Ospfv3InterfaceTimers `xml:"timers"`
}

type Ospfv3InterfaceProcessor interface {
	ospfv3InterfaceProcessor
	Ospfv3InterfaceConfigProcessor
	Ospfv3InterfaceRefProcessor
	Ospfv3InterfaceTimersProcessor
}

type ospfv3InterfaceProcessor interface {
	Ospfv3Interface(string, *NetworkInstanceProtocolKey, string, string, *Ospfv3Interface) error
}

func NewOspfv3Interface(id string) *Ospfv3Interface {
	return &Ospfv3Interface{
		SrChanges:    nclib.NewSrChanges(),
		Id:           id,
		Config:       NewOspfv3InterfaceConfig(),
		InterfaceRef: NewInterfaceRef(),
		Timers:       NewOspfv3InterfaceTimers(),
	}
}

func (o *Ospfv3Interface) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		// o.Id = value // set by NewOspfv3Interface

	case OC_CONFIG_KEY:
		if err := o.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case INTERFACE_REF_KEY:
		if err := o.InterfaceRef.Put(nodes[1:], value); err != nil {
			return nil
		}

	case OSPFV3_TIMERS_KEY:
		if err := o.Timers.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3Interface(p Ospfv3InterfaceProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, iface *Ospfv3Interface) error {
	ifaceFunc := func() error {
		if iface.GetChange(OC_ID_KEY) {
			return p.Ospfv3Interface(name, key, areaId, ifaceId, iface)
		}
		return nil
	}

	configFunc := func() error {
		if iface.GetChange(OC_CONFIG_KEY) {
			return ProcessOspfv3InterfaceConfig(
				p.(Ospfv3InterfaceConfigProcessor),
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
			return ProcessOspfv3InterfaceRef(
				p.(Ospfv3InterfaceRefProcessor),
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
		if iface.GetChange(OSPFV3_TIMERS_KEY) {
			return ProcessOspfv3InterfaceTimers(
				p.(Ospfv3InterfaceTimersProcessor),
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

type Ospfv3InterfaceConfig struct {
	nclib.SrChanges `xml:"-"`

	Id          string          `xml:"id"`
	Metric      uint16          `xml:"metric"`
	Passive     bool            `xml:"passive"`
	NetworkType OspfNetworkType `xml:"network-type"`
	Priority    uint8           `xml:"priority"`
}

type Ospfv3InterfaceConfigProcessor interface {
	Ospfv3InterfaceConfig(string, *NetworkInstanceProtocolKey, string, string, *Ospfv3InterfaceConfig) error
}

func NewOspfv3InterfaceConfig() *Ospfv3InterfaceConfig {
	return &Ospfv3InterfaceConfig{
		SrChanges:   nclib.NewSrChanges(),
		Id:          "",
		Metric:      0,
		Passive:     false,
		NetworkType: OSPF_BROADCAST_NETWORK,
		Priority:    1,
	}
}

func (c *Ospfv3InterfaceConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		c.Id = value

	case OSPFV3_METRIC_KEY:
		metric, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}
		c.Metric = uint16(metric)

	case OSPFV3_PASSIVE_KEY:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		c.Passive = b

	case OSPFV3_NETWORK_TYPE_KEY:
		n, err := ParseOspfNetworkType(value)
		if err != nil {
			return err
		}
		c.NetworkType = n

	case OSPFV3_PRIORITY_KEY:
		priority, err := strconv.ParseUint(value, 0, 8)
		if err != nil {
			return err
		}
		c.Priority = uint8(priority)
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3InterfaceConfig(p Ospfv3InterfaceConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, config *Ospfv3InterfaceConfig) error {
	configFunc := func() error {
		return p.Ospfv3InterfaceConfig(name, key, areaId, ifaceId, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

type Ospfv3InterfaceTimers struct {
	nclib.SrChanges `xml:"-"`

	HelloInterval uint32 `xml:"hello-interval"`
	DeadInterval  uint32 `xml:"dead-interval"`
}

type Ospfv3InterfaceTimersProcessor interface {
	Ospfv3InterfaceTimers(string, *NetworkInstanceProtocolKey, string, string, *Ospfv3InterfaceTimers) error
}

func NewOspfv3InterfaceTimers() *Ospfv3InterfaceTimers {
	return &Ospfv3InterfaceTimers{
		SrChanges:     nclib.NewSrChanges(),
		HelloInterval: 0,
		DeadInterval:  0,
	}
}

func (o *Ospfv3InterfaceTimers) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OSPFV3_DEAD_INTERVAL_KEY:
		t, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		o.DeadInterval = uint32(t)

	case OSPFV3_HELLO_INTERVAL_KEY:
		t, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		o.HelloInterval = uint32(t)
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3InterfaceTimers(p Ospfv3InterfaceTimersProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, timers *Ospfv3InterfaceTimers) error {
	timersFunc := func() error {
		return p.Ospfv3InterfaceTimers(name, key, areaId, ifaceId, timers)
	}

	return nclib.CallFunctions(reverse, timersFunc)
}

type Ospfv3InterfaceRefProcessor interface {
	Ospfv3InterfaceRefConfig(string, *NetworkInstanceProtocolKey, string, string, *InterfaceRefConfig) error
}

func ProcessOspfv3InterfaceRef(p Ospfv3InterfaceRefProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, ifaceId string, ifaceRef *InterfaceRef) error {
	refFunc := func() error {
		if ifaceRef.GetChange(OC_CONFIG_KEY) {
			return p.Ospfv3InterfaceRefConfig(name, key, areaId, ifaceId, ifaceRef.Config)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, refFunc)
}
