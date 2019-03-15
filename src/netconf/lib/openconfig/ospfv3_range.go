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
	"net"
	nclib "netconf/lib"
	ncnet "netconf/lib/net"
	ncxml "netconf/lib/xml"
	"strconv"
)

type Ospfv3AreaRangeKey struct {
	Ip        string
	PrefixLen uint8
}

func NewOspfv3AreaRangeKey(ip string, plen uint8) *Ospfv3AreaRangeKey {
	return &Ospfv3AreaRangeKey{
		Ip:        ip,
		PrefixLen: plen,
	}
}

func ParseOspfv3AreaRangeKey(ip, plen string) (*Ospfv3AreaRangeKey, error) {
	n, err := strconv.ParseUint(plen, 0, 8)
	if err != nil {
		return nil, err
	}

	return NewOspfv3AreaRangeKey(ip, uint8(n)), nil
}

func (o *Ospfv3AreaRangeKey) String() string {
	return fmt.Sprintf("%s/%d", o.Ip, o.PrefixLen)
}

//
// ospfv3/areas/area[id]/ranges
//
type Ospfv3AreaRanges map[Ospfv3AreaRangeKey]*Ospfv3AreaRange

func NewOspfv3AreaRanges() Ospfv3AreaRanges {
	return Ospfv3AreaRanges{}
}

func (o Ospfv3AreaRanges) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ip, ok := nodes[0].Attrs[OSPFV3_RANGE_IP_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", OSPFV3_RANGE_KEY, OSPFV3_RANGE_IP_KEY, nodes[0])
	}
	plen, ok := nodes[0].Attrs[OSPFV3_RANGE_PREFIXLEN_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", OSPFV3_RANGE_KEY, OSPFV3_RANGE_PREFIXLEN_KEY, nodes[0])
	}
	key, err := ParseOspfv3AreaRangeKey(ip, plen)
	if err != nil {
		return err
	}

	rng, ok := o[*key]
	if !ok {
		rng = NewOspfv3AreaRange(key)
		o[*key] = rng
	}

	return rng.Put(nodes[1:], value)
}

func ProcessOspfv3AreaRanges(p Ospfv3AreaRangeProcessor, reverse bool, name string, nikey *NetworkInstanceProtocolKey, areaId string, ranges Ospfv3AreaRanges) error {
	for rngkey, rng := range ranges {
		if err := ProcessOspfv3AreaRange(p, reverse, name, nikey, areaId, &rngkey, rng); err != nil {
			return err
		}
	}

	return nil
}

func (o Ospfv3AreaRanges) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = OSPFV3_RANGES_KEY
	e.EncodeToken(start)

	for _, rng := range o {
		err := e.EncodeElement(rng, xml.StartElement{Name: xml.Name{Local: OSPFV3_RANGE_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type Ospfv3AreaRange struct {
	nclib.SrChanges `xml:"-"`

	Ip        string                 `xml:"ip"`
	PrefixLen uint8                  `xml:"prefix-length"`
	Config    *Ospfv3AreaRangeConfig `xml:"config"`
}

type Ospfv3AreaRangeProcessor interface {
	ospfv3AreaRangeProcessor
	Ospfv3AreaRangeConfigProcessor
}

type ospfv3AreaRangeProcessor interface {
	Ospfv3AreaRange(string, *NetworkInstanceProtocolKey, string, *Ospfv3AreaRangeKey, *Ospfv3AreaRange) error
}

func NewOspfv3AreaRange(key *Ospfv3AreaRangeKey) *Ospfv3AreaRange {
	return &Ospfv3AreaRange{
		SrChanges: nclib.NewSrChanges(),
		Ip:        key.Ip,
		PrefixLen: key.PrefixLen,
		Config:    NewOspfv3AreaRangeConfig(),
	}
}

func (o *Ospfv3AreaRange) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OSPFV3_RANGE_IP_KEY:
		// o.Ip = value // set by NewOspfv3AreaRange

	case OSPFV3_RANGE_PREFIXLEN_KEY:
		// o.PrefixLen = strings.ParseUint(value, 0,8) // set by NewOspfv3AreaRange

	case OC_CONFIG_KEY:
		if err := o.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3AreaRange(p Ospfv3AreaRangeProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, areaId string, rngkey *Ospfv3AreaRangeKey, rng *Ospfv3AreaRange) error {
	rangeFunc := func() error {
		if rng.GetChanges(OSPFV3_RANGE_IP_KEY, OSPFV3_RANGE_PREFIXLEN_KEY) {
			return p.Ospfv3AreaRange(name, key, areaId, rngkey, rng)
		}
		return nil
	}

	configFunc := func() error {
		if rng.GetChange(OC_CONFIG_KEY) {
			return ProcessOspfv3AreaRangeConfig(
				p.(Ospfv3AreaRangeConfigProcessor),
				reverse,
				name,
				key,
				areaId,
				rngkey,
				rng.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, rangeFunc, configFunc)
}

type Ospfv3AreaRangeConfig struct {
	nclib.SrChanges `xml:"-"`

	Ip        net.IP `xml:"ip"`
	PrefixLen uint8  `xml:"prefix-length"`
}

type Ospfv3AreaRangeConfigProcessor interface {
	Ospfv3AreaRangeConfig(string, *NetworkInstanceProtocolKey, string, *Ospfv3AreaRangeKey, *Ospfv3AreaRangeConfig) error
}

func NewOspfv3AreaRangeConfig() *Ospfv3AreaRangeConfig {
	return &Ospfv3AreaRangeConfig{
		SrChanges: nclib.NewSrChanges(),
		Ip:        nil,
		PrefixLen: 0,
	}
}

func (o *Ospfv3AreaRangeConfig) IPNet() *net.IPNet {
	return ncnet.IPToIPNet(o.Ip, int(o.PrefixLen))
}

func (o *Ospfv3AreaRangeConfig) SetIP(ip net.IP) {
	o.Ip = ip
	o.SetChange(OSPFV3_RANGE_IP_KEY)
}

func (o *Ospfv3AreaRangeConfig) SetPLen(plen uint8) {
	o.PrefixLen = plen
	o.SetChange(OSPFV3_RANGE_PREFIXLEN_KEY)
}

func (o *Ospfv3AreaRangeConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OSPFV3_RANGE_IP_KEY:
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("Invalid IP. %s", value)
		}
		o.Ip = ip

	case OSPFV3_RANGE_PREFIXLEN_KEY:
		prefixLen, err := strconv.ParseUint(value, 0, 8)
		if err != nil {
			return err
		}
		o.PrefixLen = uint8(prefixLen)
	}

	o.SetChange(nodes[0].Name)
	return nil
}

func ProcessOspfv3AreaRangeConfig(p Ospfv3AreaRangeConfigProcessor, reverse bool, name string, nikey *NetworkInstanceProtocolKey, areaId string, rngkey *Ospfv3AreaRangeKey, config *Ospfv3AreaRangeConfig) error {
	configFunc := func() error {
		return p.Ospfv3AreaRangeConfig(name, nikey, areaId, rngkey, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
