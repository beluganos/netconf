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
	nclib "netconf/lib"
	ncxml "netconf/lib/xml"
)

const (
	OSPFV3_KEY                 = "ospfv3"
	OSPFV3_ROUTERID_KEY        = "router-id"
	OSPFV3_AREAS_KEY           = "areas"
	OSPFV3_AREA_KEY            = "area"
	OSPFV3_METRIC_KEY          = "metric"
	OSPFV3_PASSIVE_KEY         = "passive"
	OSPFV3_NETWORK_TYPE_KEY    = "network-type"
	OSPFV3_PRIORITY_KEY        = "priority"
	OSPFV3_TIMERS_KEY          = "timers"
	OSPFV3_DEAD_INTERVAL_KEY   = "dead-interval"
	OSPFV3_HELLO_INTERVAL_KEY  = "hello-interval"
	OSPFV3_RANGES_KEY          = "ranges"
	OSPFV3_RANGE_KEY           = "range"
	OSPFV3_RANGE_IP_KEY        = "ip"
	OSPFV3_RANGE_PREFIXLEN_KEY = "prefix-length"
)

//
// ospfv3
//
type Ospfv3 struct {
	nclib.SrChanges `xml:"-"`

	Global *Ospfv3Global `xml:"global"`
	Areas  Ospfv3Areas   `xml:"areas"`
}

type Ospfv3Processor interface {
	ospfv3Processor
	Ospfv3GlobalProcessor
	Ospfv3AreaProcessor
}

type ospfv3Processor interface {
	Ospfv3(string, *NetworkInstanceProtocolKey, *Ospfv3) error
}

func NewOspfv3() *Ospfv3 {
	return &Ospfv3{
		SrChanges: nclib.NewSrChanges(),
		Global:    NewOspfv3Global(),
		Areas:     NewOspfv3Areas(),
	}
}

func (o *Ospfv3) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_GLOBAL_KEY:
		if err := o.Global.Put(nodes[1:], value); err != nil {
			return err
		}

	case OSPFV3_AREAS_KEY:
		if err := o.Areas.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3(p Ospfv3Processor, reverse bool, name string, key *NetworkInstanceProtocolKey, ospf *Ospfv3) error {
	ospfFunc := func() error {
		return p.Ospfv3(name, key, ospf)
	}

	globalFunc := func() error {
		if ospf.GetChange(OC_GLOBAL_KEY) {
			return ProcessOspfv3Global(
				p.(Ospfv3GlobalProcessor),
				reverse,
				name,
				key,
				ospf.Global,
			)
		}
		return nil
	}

	areasFunc := func() error {
		if ospf.GetChange(OSPFV3_AREAS_KEY) {
			return ProcessOspfv3Areas(
				p.(Ospfv3AreaProcessor),
				reverse,
				name,
				key,
				ospf.Areas,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, ospfFunc, globalFunc, areasFunc)
}

//
// ospfv3/global
//
type Ospfv3Global struct {
	nclib.SrChanges `xml:"-"`

	Config *Ospfv3GlobalConfig `xml:"config"`
}

type Ospfv3GlobalProcessor interface {
	Ospfv3GlobalConfigProcessor
}

func NewOspfv3Global() *Ospfv3Global {
	return &Ospfv3Global{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewOspfv3GlobalConfig(),
	}
}

func (o *Ospfv3Global) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := o.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3Global(p Ospfv3GlobalProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, global *Ospfv3Global) error {
	configFunc := func() error {
		if global.GetChange(OC_CONFIG_KEY) {
			return ProcessOspfv3GlobalConfig(
				p.(Ospfv3GlobalConfigProcessor),
				reverse,
				name,
				key,
				global.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// ospfv3/global/config
//
type Ospfv3GlobalConfig struct {
	nclib.SrChanges `xml:"-"`

	RouterId string `xml:"router-id"`
}

type Ospfv3GlobalConfigProcessor interface {
	Ospfv3GlobalConfig(string, *NetworkInstanceProtocolKey, *Ospfv3GlobalConfig) error
}

func NewOspfv3GlobalConfig() *Ospfv3GlobalConfig {
	return &Ospfv3GlobalConfig{
		SrChanges: nclib.NewSrChanges(),
		RouterId:  "",
	}
}

func (c *Ospfv3GlobalConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OSPFV3_ROUTERID_KEY:
		c.RouterId = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3GlobalConfig(p Ospfv3GlobalConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, config *Ospfv3GlobalConfig) error {
	configFunc := func() error {
		return p.Ospfv3GlobalConfig(name, key, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
