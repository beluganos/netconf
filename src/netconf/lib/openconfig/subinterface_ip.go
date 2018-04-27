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
	"netconf/lib"
	"netconf/lib/xml"
	"strconv"
)

//
// subinterface[index]/ipv4
//
type SubinterfaceIPv4 struct {
	nclib.SrChanges `xml:"-"`

	XMLName   xml.Name                `xml:"https://github.com/beluganos/beluganos/yang/interfaces/ip ipv4"`
	Addresses IPAddresses             `xml:"addresses"`
	Config    *SubinterfaceIPv4Config `xml:"config"`
}

type SubinterfaceIPv4Processor interface {
	SubinterfaceIPv4AddrProcessor
	SubinterfaceIPv4ConfigProcessor
}

func NewSubinterfaceIPv4() *SubinterfaceIPv4 {
	return &SubinterfaceIPv4{
		SrChanges: nclib.NewSrChanges(),
		Addresses: NewIPAddresses(),
		Config:    NewSubinterfaceIPv4Config(),
	}
}

func (i *SubinterfaceIPv4) String() string {
	return fmt.Sprintf("%s{%s, %s} %s",
		SUBINTERFACE_IPV4_KEY,
		i.Addresses,
		i.Config,
		i.SrChanges,
	)
}

func (i *SubinterfaceIPv4) SetConfig(config *SubinterfaceIPv4Config) {
	i.Config = config
	i.SetChange(OC_CONFIG_KEY)
}

func (i *SubinterfaceIPv4) SetAddresses(addrs IPAddresses) {
	i.Addresses = addrs
	i.SetChange(SUBINTERFACE_ADDRS_KEY)
}

func (i *SubinterfaceIPv4) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := i.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case SUBINTERFACE_ADDRS_KEY:
		if err := i.Addresses.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	i.SetChange(nodes[0].Name)
	return nil
}

func ProcessSubinterfaceIPv4(p SubinterfaceIPv4Processor, reverse bool, name string, index uint32, ipv4 *SubinterfaceIPv4) error {

	configFunc := func() error {
		if ipv4.GetChange(OC_CONFIG_KEY) {
			return ProcessSubinterfaceIPv4Config(
				p.(SubinterfaceIPv4ConfigProcessor),
				reverse,
				name,
				index,
				ipv4.Config,
			)
		}
		return nil
	}

	addrsFunc := func() error {
		if ipv4.GetChange(SUBINTERFACE_ADDRS_KEY) {
			return ProcessSubinterfaceIPv4Addrs(
				p.(SubinterfaceIPv4AddrProcessor),
				reverse,
				name,
				index,
				ipv4.Addresses,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc, addrsFunc)
}

//
//  subinterface[index]/ipv4/config
//
type SubinterfaceIPv4Config struct {
	nclib.SrChanges `xml:"-"`

	Mtu uint16 `xml:"mtu,omitempty"`
}

type SubinterfaceIPv4ConfigProcessor interface {
	SubinterfaceIPv4Config(string, uint32, *SubinterfaceIPv4Config) error
}

func NewSubinterfaceIPv4Config() *SubinterfaceIPv4Config {
	return &SubinterfaceIPv4Config{
		SrChanges: nclib.NewSrChanges(),
		Mtu:       0,
	}
}

func (c *SubinterfaceIPv4Config) String() string {
	return fmt.Sprintf("%s{%s=%d} %s",
		OC_CONFIG_KEY,
		SUBINTERFACE_MTU_KEY, c.Mtu,
		c.SrChanges,
	)
}

func (c *SubinterfaceIPv4Config) SetMtu(mtu uint16) {
	c.Mtu = mtu
	c.SetChange(SUBINTERFACE_MTU_KEY)
}

func (c *SubinterfaceIPv4Config) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case SUBINTERFACE_MTU_KEY:
		mtu, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}
		c.Mtu = uint16(mtu)
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessSubinterfaceIPv4Config(p SubinterfaceIPv4ConfigProcessor, reverse bool, name string, index uint32, config *SubinterfaceIPv4Config) error {
	configFunc := func() error {
		return p.SubinterfaceIPv4Config(name, index, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// subinterface[index]/ipv4/addresses
//
func ProcessSubinterfaceIPv4Addrs(p SubinterfaceIPv4AddrProcessor, reverse bool, name string, index uint32, addrs IPAddresses) error {
	for ip, addr := range addrs {
		if err := ProcessSubinterfaceIPv4Addr(p, reverse, name, index, ip, addr); err != nil {
			return err
		}
	}
	return nil
}

//
// subinterface[index]/ipv4/addresses/address[ip]
//
type SubinterfaceIPv4AddrProcessor interface {
	subinterfaceIPv4AddrProcessor
	SubinterfaceIPv4AddressConfigProcessor
}

type subinterfaceIPv4AddrProcessor interface {
	SubinterfaceIPv4Address(string, uint32, string, *IPAddress) error
}

func ProcessSubinterfaceIPv4Addr(p subinterfaceIPv4AddrProcessor, reverse bool, name string, index uint32, ip string, addr *IPAddress) error {
	addrFunc := func() error {
		if addr.GetChange(SUBINTERFACE_ADDR_IP_KEY) {
			return p.SubinterfaceIPv4Address(name, index, ip, addr)
		}
		return nil
	}

	configFunc := func() error {
		if addr.GetChange(OC_CONFIG_KEY) {
			return ProcessSubinterfaceIPv4AddressConfig(
				p.(SubinterfaceIPv4AddressConfigProcessor),
				reverse,
				name,
				index,
				ip,
				addr.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, addrFunc, configFunc)
}

//
// subinterface[index]/ipv4/addresses/address[ip]/config
//
type SubinterfaceIPv4AddressConfigProcessor interface {
	SubinterfaceIPv4AddressConfig(string, uint32, string, *IPAddressConfig) error
}

func ProcessSubinterfaceIPv4AddressConfig(p SubinterfaceIPv4AddressConfigProcessor, reverse bool, name string, index uint32, ip string, config *IPAddressConfig) error {
	configFunc := func() error {
		return p.SubinterfaceIPv4AddressConfig(name, index, ip, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
