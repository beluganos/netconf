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
	"strconv"
)

const (
	STATICROUTES_KEY          = "static-routes"
	STATICROUTE_KEY           = "static"
	STATICROUTE_IP_KEY        = "ip"
	STATICROUTE_PREFIXLEN_KEY = "prefix-length"
	STATICROUTE_NEXTHOPS_KEY  = "next-hops"
	STATICROUTE_NEXTHOP_KEY   = "next-hop"
)

//
// static-route-key
//
type StaticRouteKey struct {
	IP        string
	PrefixLen uint8
}

func NewStaticRouteKey(ip string, prefixlen uint8) *StaticRouteKey {
	return &StaticRouteKey{
		IP:        ip,
		PrefixLen: prefixlen,
	}
}

func ParseStaticRouteKey(ip string, prefixlen string) (*StaticRouteKey, error) {
	plen, err := strconv.ParseUint(prefixlen, 0, 8)
	if err != nil {
		return nil, err
	}
	return NewStaticRouteKey(ip, uint8(plen)), nil
}

func (s *StaticRouteKey) String() string {
	return fmt.Sprintf("%s/%d", s.IP, s.PrefixLen)
}

//
// static-routes
//
type StaticRoutes map[StaticRouteKey]*StaticRoute

func NewStaticRoutes() StaticRoutes {
	return StaticRoutes{}
}

func (s StaticRoutes) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ip, ok := nodes[0].Attrs[STATICROUTE_IP_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", STATICROUTE_KEY, STATICROUTE_IP_KEY, nodes[0])
	}

	plen, ok := nodes[0].Attrs[STATICROUTE_PREFIXLEN_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", STATICROUTE_KEY, STATICROUTE_PREFIXLEN_KEY, nodes[0])
	}

	key, err := ParseStaticRouteKey(ip, plen)
	if err != nil {
		return err
	}

	route, ok := s[*key]
	if !ok {
		route = NewStaticRoute(key)
		s[*key] = route
	}

	return route.Put(nodes[1:], value)
}

func ProcessStaticRoutes(p StaticRouteProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, routes StaticRoutes) error {
	for rtkey, route := range routes {
		if err := ProcessStaticRoute(p, reverse, name, key, &rtkey, route); err != nil {
			return err
		}
	}
	return nil
}

func (s StaticRoutes) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = STATICROUTES_KEY
	e.EncodeToken(start)

	for _, route := range s {
		err := e.EncodeElement(route, xml.StartElement{Name: xml.Name{Local: STATICROUTE_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// static-routes/static[ip, prefix-length]
//
type StaticRoute struct {
	nclib.SrChanges `xml:"-"`

	IP        string              `xml:"ip"`
	PrefixLen uint8               `xml:"prefix-length"`
	Config    *StaticRouteConfig  `xml:"config"`
	Nexthops  StaticRouteNexthops `xml:"nexthops"`
}

type StaticRouteProcessor interface {
	staticRouteProcessor
	StaticRouteConfigProcessor
	StaticRouteNexthopProcessor
}

type staticRouteProcessor interface {
	StaticRoute(string, *NetworkInstanceProtocolKey, *StaticRouteKey, *StaticRoute) error
}

func NewStaticRoute(key *StaticRouteKey) *StaticRoute {
	return &StaticRoute{
		SrChanges: nclib.NewSrChanges(),
		IP:        key.IP,
		PrefixLen: key.PrefixLen,
		Config:    NewStaticRouteConfig(),
		Nexthops:  NewStaticRouteNexthops(),
	}
}

func (s *StaticRoute) String() string {
	return fmt.Sprintf("%s{%s='%s', %s=%d, %s, %s} %s",
		STATICROUTE_KEY,
		STATICROUTE_IP_KEY, s.IP,
		STATICROUTE_PREFIXLEN_KEY, s.PrefixLen,
		s.Config,
		s.Nexthops,
		s.SrChanges,
	)
}

func (s *StaticRoute) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case STATICROUTE_IP_KEY:
		// s.Ip = value // set by NewStaticRoute

	case STATICROUTE_PREFIXLEN_KEY:
		// s.PrefxLen = value // set by NewStaticRoute

	case OC_CONFIG_KEY:
		if err := s.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case STATICROUTE_NEXTHOPS_KEY:
		if err := s.Nexthops.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	s.SetChange(nodes[0].Name)
	return nil
}

func ProcessStaticRoute(p StaticRouteProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, rtkey *StaticRouteKey, route *StaticRoute) error {
	rtFunc := func() error {
		if route.GetChanges(STATICROUTE_IP_KEY, STATICROUTE_PREFIXLEN_KEY) {
			return p.StaticRoute(name, key, rtkey, route)
		}
		return nil
	}

	configFunc := func() error {
		if route.GetChange(OC_CONFIG_KEY) {
			return ProcessStaticRouteConfig(
				p.(StaticRouteConfigProcessor),
				reverse,
				name,
				key,
				rtkey,
				route.Config,
			)
		}
		return nil
	}

	nhFunc := func() error {
		if route.GetChange(STATICROUTE_NEXTHOPS_KEY) {
			return ProcessStaticRouteNexthops(
				p.(StaticRouteNexthopProcessor),
				reverse,
				name,
				key,
				rtkey,
				route.Nexthops,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, rtFunc, configFunc, nhFunc)
}

//
// static-routes/static/config
//
type StaticRouteConfig struct {
	nclib.SrChanges `xml:"-"`

	Ip        net.IP `xml:"ip"`
	PrefixLen uint8  `xml:"prefix-length"`
}

type StaticRouteConfigProcessor interface {
	StaticRouteConfig(string, *NetworkInstanceProtocolKey, *StaticRouteKey, *StaticRouteConfig) error
}

func NewStaticRouteConfig() *StaticRouteConfig {
	return &StaticRouteConfig{
		SrChanges: nclib.NewSrChanges(),
		Ip:        nil,
		PrefixLen: 0,
	}
}

func (c *StaticRouteConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s=%d} %s",
		OC_CONFIG_KEY,
		STATICROUTE_IP_KEY, c.Ip,
		STATICROUTE_PREFIXLEN_KEY, c.PrefixLen,
		c.SrChanges,
	)
}

func (c *StaticRouteConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case STATICROUTE_IP_KEY:
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("Invalid IP. %s", value)
		}
		c.Ip = ip

	case STATICROUTE_PREFIXLEN_KEY:
		v, err := strconv.ParseUint(value, 0, 8)
		if err != nil {
			return err
		}
		c.PrefixLen = uint8(v)
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessStaticRouteConfig(p StaticRouteConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, rtkey *StaticRouteKey, config *StaticRouteConfig) error {
	configFunc := func() error {
		return p.StaticRouteConfig(name, key, rtkey, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
