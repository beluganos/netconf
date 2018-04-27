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
	"netconf/lib"
	"netconf/lib/xml"
)

const (
	OSPFV2_KEY                = "ospfv2"
	OSPFV2_ROUTERID_KEY       = "router-id"
	OSPFV2_AREAS_KEY          = "areas"
	OSPFV2_AREA_KEY           = "area"
	OSPFV2_METRIC_KEY         = "metric"
	OSPFV2_PASSIVE_KEY        = "passive"
	OSPFV2_TIMERS_KEY         = "timers"
	OSPFV2_DEAD_INTERVAL_KEY  = "dead-interval"
	OSPFV2_HELLO_INTERVAL_KEY = "hello-interval"
)

//
// ospfv2
//
type Ospfv2 struct {
	nclib.SrChanges `xml:"-"`

	Global *Ospfv2Global `xml:"global"`
	Areas  Ospfv2Areas   `xml:"areas"`
}

type Ospfv2Processor interface {
	ospfv2Processor
	Ospfv2GlobalProcessor
	Ospfv2AreaProcessor
}

type ospfv2Processor interface {
	Ospfv2(string, *NetworkInstanceProtocolKey, *Ospfv2) error
}

func NewOspfv2() *Ospfv2 {
	return &Ospfv2{
		SrChanges: nclib.NewSrChanges(),
		Global:    NewOspfv2Global(),
		Areas:     NewOspfv2Areas(),
	}
}

func (o *Ospfv2) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_GLOBAL_KEY:
		if err := o.Global.Put(nodes[1:], value); err != nil {
			return err
		}

	case OSPFV2_AREAS_KEY:
		if err := o.Areas.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv2(p Ospfv2Processor, reverse bool, name string, key *NetworkInstanceProtocolKey, ospf *Ospfv2) error {
	ospfFunc := func() error {
		return p.Ospfv2(name, key, ospf)
	}

	globalFunc := func() error {
		if ospf.GetChange(OC_GLOBAL_KEY) {
			return ProcessOspfv2Global(
				p.(Ospfv2GlobalProcessor),
				reverse,
				name,
				key,
				ospf.Global,
			)
		}
		return nil
	}

	areasFunc := func() error {
		if ospf.GetChange(OSPFV2_AREAS_KEY) {
			return ProcessOspfv2Areas(
				p.(Ospfv2AreaProcessor),
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
// ospfv2/global
//
type Ospfv2Global struct {
	nclib.SrChanges `xml:"-"`

	Config *Ospfv2GlobalConfig `xml:"config"`
}

type Ospfv2GlobalProcessor interface {
	Ospfv2GlobalConfigProcessor
}

func NewOspfv2Global() *Ospfv2Global {
	return &Ospfv2Global{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewOspfv2GlobalConfig(),
	}
}

func (o *Ospfv2Global) Put(nodes []*ncxml.XPathNode, value string) error {
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

func ProcessOspfv2Global(p Ospfv2GlobalProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, global *Ospfv2Global) error {
	configFunc := func() error {
		if global.GetChange(OC_CONFIG_KEY) {
			return ProcessOspfv2GlobalConfig(
				p.(Ospfv2GlobalConfigProcessor),
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
// ospfv2/global/config
//
type Ospfv2GlobalConfig struct {
	nclib.SrChanges `xml:"-"`

	RouterId string `xml:"router-id"`
}

type Ospfv2GlobalConfigProcessor interface {
	Ospfv2GlobalConfig(string, *NetworkInstanceProtocolKey, *Ospfv2GlobalConfig) error
}

func NewOspfv2GlobalConfig() *Ospfv2GlobalConfig {
	return &Ospfv2GlobalConfig{
		SrChanges: nclib.NewSrChanges(),
		RouterId:  "",
	}
}

func (c *Ospfv2GlobalConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OSPFV2_ROUTERID_KEY:
		c.RouterId = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv2GlobalConfig(p Ospfv2GlobalConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, config *Ospfv2GlobalConfig) error {
	configFunc := func() error {
		return p.Ospfv2GlobalConfig(name, key, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
