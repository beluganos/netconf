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
	nclib "netconf/lib"
	ncnet "netconf/lib/net"
	ncxml "netconf/lib/xml"
	"strconv"
)

//
// network-instances/network-instance[name]/loopbacks
//
type NetworkInstanceLoopbacks map[string]*NetworkInstanceLoopback

func NewNetworkInstanceLoopbacks() NetworkInstanceLoopbacks {
	return NetworkInstanceLoopbacks{}
}

func (n NetworkInstanceLoopbacks) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	id, ok := nodes[0].Attrs[OC_ID_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", NETWORKINSTANCE_LO_KEY, OC_ID_KEY, nodes[0])
	}

	lo, ok := n[id]
	if !ok {
		lo = NewNetworkInstanceLoopback(id)
		n[id] = lo
	}

	return lo.Put(nodes[1:], value)
}

func ProcessNetworkInstanceLoopbacks(p NetworkInstanceLoopbackProcessor, reverse bool, name string, los NetworkInstanceLoopbacks) error {
	for id, lo := range los {
		if err := ProcessNetworkInstanceLoopback(p, reverse, name, id, lo); err != nil {
			return err
		}
	}
	return nil
}

func (i NetworkInstanceLoopbacks) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = NETWORKINSTANCE_LOS_KEY
	e.EncodeToken(start)

	for _, iface := range i {
		err := e.EncodeElement(iface, xml.StartElement{Name: xml.Name{Local: NETWORKINSTANCE_LO_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// network-instances/network-instance[name]/loopbacks/loopback[id]
//
type NetworkInstanceLoopback struct {
	nclib.SrChanges `xml:"-"`

	Id     string                         `xml:"id"`
	Config *NetworkInstanceLoopbackConfig `xml:"config"`
	Addrs  NetworkInstanceLoopbackAddrs   `xml:"addresses"`
}

type NetworkInstanceLoopbackProcessor interface {
	networkInstanceLoopbackProcessor
	NetworkInstanceLoopbackConfigProcessor
	NetworkInstanceLoopbackAddrProcessor
}

type networkInstanceLoopbackProcessor interface {
	NetworkInstanceLoopback(string, string, *NetworkInstanceLoopback) error
}

func NewNetworkInstanceLoopback(id string) *NetworkInstanceLoopback {
	return &NetworkInstanceLoopback{
		SrChanges: nclib.NewSrChanges(),
		Id:        id,
		Config:    NewNetworkInstanceLoopbackConfig(),
		Addrs:     NewNetworkInstanceLoopbackAddrs(),
	}
}

func (n *NetworkInstanceLoopback) String() string {
	return fmt.Sprintf("%s{%s='%s', %s} %s",
		NETWORKINSTANCE_LO_KEY,
		OC_ID_KEY, n.Id,
		n.Config,
		n.SrChanges,
	)
}

func (n *NetworkInstanceLoopback) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		// n.Id = value // NewNetworkInstanceLoopback

	case OC_CONFIG_KEY:
		if err := n.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case NETWORKINSTANCE_LO_ADDRS_KEY:
		if err := n.Addrs.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	n.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceLoopback(p NetworkInstanceLoopbackProcessor, reverse bool, name string, id string, lo *NetworkInstanceLoopback) error {

	idFunc := func() error {
		if lo.GetChange(OC_ID_KEY) {
			return p.NetworkInstanceLoopback(name, id, lo)
		}
		return nil
	}

	configFunc := func() error {
		if lo.GetChange(OC_CONFIG_KEY) {
			return ProcessNetworkInstanceLoopbackConfig(
				p.(NetworkInstanceLoopbackConfigProcessor),
				reverse,
				name,
				id,
				lo.Config,
			)
		}
		return nil
	}

	addrsFunc := func() error {
		if lo.GetChange(NETWORKINSTANCE_LO_ADDRS_KEY) {
			return ProcessNetworkInstanceLoopbackAddrs(
				p.(NetworkInstanceLoopbackAddrProcessor),
				reverse,
				name,
				id,
				lo.Addrs,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, idFunc, configFunc, addrsFunc)
}

//
// network-instances/network-instance[name]/loopbacks/loopback[id] /config
//
type NetworkInstanceLoopbackConfig struct {
	nclib.SrChanges `xml:"-"`

	Id string `xml:"id"`
}

type NetworkInstanceLoopbackConfigProcessor interface {
	NetworkInstanceLoopbackConfig(string, string, *NetworkInstanceLoopbackConfig) error
}

func NewNetworkInstanceLoopbackConfig() *NetworkInstanceLoopbackConfig {
	return &NetworkInstanceLoopbackConfig{
		SrChanges: nclib.NewSrChanges(),
		Id:        "",
	}
}

func (c *NetworkInstanceLoopbackConfig) String() string {
	return fmt.Sprintf("%s{%s='%s'} %s",
		OC_CONFIG_KEY,
		OC_ID_KEY, c.Id,
		c.SrChanges,
	)
}

func (c *NetworkInstanceLoopbackConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ID_KEY:
		c.Id = value
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceLoopbackConfig(p NetworkInstanceLoopbackConfigProcessor, reverse bool, name string, id string, config *NetworkInstanceLoopbackConfig) error {

	configFunc := func() error {
		return p.NetworkInstanceLoopbackConfig(name, id, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// network-instances/network-instance[name]/loopbacks/loopback[id]/addresses
//
type NetworkInstanceLoopbackAddrs map[string]*NetworkInstanceLoopbackAddr

func NewNetworkInstanceLoopbackAddrs() NetworkInstanceLoopbackAddrs {
	return NetworkInstanceLoopbackAddrs{}
}

func (n NetworkInstanceLoopbackAddrs) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	index, ok := nodes[0].Attrs[OC_INDEX_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", NETWORKINSTANCE_LO_ADDRS_KEY, OC_INDEX_KEY, nodes[0])
	}

	addr, ok := n[index]
	if !ok {
		addr = NewNetworkInstanceLoopbackAddr(index)
		n[index] = addr
	}

	return addr.Put(nodes[1:], value)
}

func ProcessNetworkInstanceLoopbackAddrs(p NetworkInstanceLoopbackAddrProcessor, reverse bool, name string, id string, addrs NetworkInstanceLoopbackAddrs) error {
	for index, addr := range addrs {
		if err := ProcessNetworkInstanceLoopbackAddr(p, reverse, name, id, index, addr); err != nil {
			return err
		}
	}
	return nil
}

func (i NetworkInstanceLoopbackAddrs) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = NETWORKINSTANCE_LO_ADDRS_KEY
	e.EncodeToken(start)

	for _, iface := range i {
		err := e.EncodeElement(iface, xml.StartElement{Name: xml.Name{Local: NETWORKINSTANCE_LO_ADDR_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// network-instances/network-instance[name]/loopbacks/loopback[id]/addresses/address[index]
//
type NetworkInstanceLoopbackAddr struct {
	nclib.SrChanges `xml:"-"`

	Index  string                             `xml:"index"`
	Config *NetworkInstanceLoopbackAddrConfig `xml:"config"`
}

type NetworkInstanceLoopbackAddrProcessor interface {
	networkInstanceLoopbackAddrProcessor
	NetworkInstanceLoopbackAddrConfigProcessor
}

type networkInstanceLoopbackAddrProcessor interface {
	NetworkInstanceLoopbackAddr(string, string, string, *NetworkInstanceLoopbackAddr) error
}

func NewNetworkInstanceLoopbackAddr(index string) *NetworkInstanceLoopbackAddr {
	return &NetworkInstanceLoopbackAddr{
		SrChanges: nclib.NewSrChanges(),
		Index:     index,
		Config:    NewNetworkInstanceLoopbackAddrConfig(),
	}
}

func (n *NetworkInstanceLoopbackAddr) String() string {
	return fmt.Sprintf("%s{%s='%s', %s} %s",
		NETWORKINSTANCE_LO_ADDR_KEY,
		OC_INDEX_KEY, n.Index,
		n.Config,
		n.SrChanges,
	)
}

func (n *NetworkInstanceLoopbackAddr) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_INDEX_KEY:
		// n.Id = value // NewNetworkInstanceLoopbackAddr

	case OC_CONFIG_KEY:
		if err := n.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	n.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceLoopbackAddr(p NetworkInstanceLoopbackAddrProcessor, reverse bool, name string, id string, index string, addr *NetworkInstanceLoopbackAddr) error {

	indexFunc := func() error {
		if addr.GetChange(OC_INDEX_KEY) {
			return p.NetworkInstanceLoopbackAddr(name, id, index, addr)
		}
		return nil
	}

	configFunc := func() error {
		if addr.GetChange(OC_CONFIG_KEY) {
			return ProcessNetworkInstanceLoopbackAddrConfig(
				p.(NetworkInstanceLoopbackAddrConfigProcessor),
				reverse,
				name,
				id,
				index,
				addr.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, indexFunc, configFunc)
}

//
// network-instances/network-instance[name]/loopbacks/loopback[id]/addresses/address[index]/config
//
type NetworkInstanceLoopbackAddrConfig struct {
	nclib.SrChanges `xml:"-"`

	Index     string `xml:"id"`
	Ip        net.IP `xml:"ip"`
	PrefixLen uint8  `xml:"prefix-length"`
}

type NetworkInstanceLoopbackAddrConfigProcessor interface {
	NetworkInstanceLoopbackAddrConfig(string, string, string, *NetworkInstanceLoopbackAddrConfig) error
}

func NewNetworkInstanceLoopbackAddrConfig() *NetworkInstanceLoopbackAddrConfig {
	return &NetworkInstanceLoopbackAddrConfig{
		SrChanges: nclib.NewSrChanges(),
		Index:     "",
		Ip:        nil,
		PrefixLen: 0,
	}
}

func (c *NetworkInstanceLoopbackAddrConfig) IFAddr() *ncnet.IFAddr {
	return ncnet.NewIFAddrWithPlen(c.Ip, c.PrefixLen)
}

func (c *NetworkInstanceLoopbackAddrConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s/%d} %s",
		OC_CONFIG_KEY,
		OC_INDEX_KEY, c.Index,
		c.Ip, c.PrefixLen,
		c.SrChanges,
	)
}

func (c *NetworkInstanceLoopbackAddrConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_INDEX_KEY:
		c.Index = value

	case NETWORKINSTANCE_LO_IP_KEY:
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("Invalid ip address. '%s'", value)
		}
		c.Ip = ip

	case NETWORKINSTANCE_LO_PLEN_KEY:
		plen, err := strconv.ParseUint(value, 0, 8)
		if err != nil {
			return err
		}
		c.PrefixLen = uint8(plen)
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessNetworkInstanceLoopbackAddrConfig(p NetworkInstanceLoopbackAddrConfigProcessor, reverse bool, name string, id string, index string, config *NetworkInstanceLoopbackAddrConfig) error {

	configFunc := func() error {
		return p.NetworkInstanceLoopbackAddrConfig(name, id, index, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
