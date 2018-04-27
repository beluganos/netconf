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
)

//
// subinterface[index]/ethernet
//
type InterfaceEthernet struct {
	nclib.SrChanges `xml:"-"`

	XMLName xml.Name                 `xml:"https://github.com/beluganos/beluganos/yang/interfaces/ethernet"`
	Config  *InterfaceEthernetConfig `xml:"config"`
}

type InterfaceEthernetProcessor interface {
	InterfaceEthernetConfigProcessor
}

func NewInterfaceEthernet() *InterfaceEthernet {
	return &InterfaceEthernet{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewInterfaceEthernetConfig(),
	}
}

func (i *InterfaceEthernet) String() string {
	return fmt.Sprintf("%s{%s} %s",
		INTERFACE_ETH_KEY,
		i.Config,
		i.SrChanges,
	)
}

func (i *InterfaceEthernet) SetConfig(config *InterfaceEthernetConfig) {
	i.Config = config
	i.SetChange(OC_CONFIG_KEY)
}

func (i *InterfaceEthernet) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := i.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	i.SetChange(nodes[0].Name)
	return nil
}

func ProcessInterfaceEthernet(p InterfaceEthernetProcessor, reverse bool, name string, ethernet *InterfaceEthernet) error {

	configFunc := func() error {
		if ethernet.GetChange(OC_CONFIG_KEY) {
			return ProcessInterfaceEthernetConfig(
				p.(InterfaceEthernetConfigProcessor),
				reverse,
				name,
				ethernet.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// subinterface[index]/ethernet/config
//
type InterfaceEthernetConfig struct {
	nclib.SrChanges `xml:"-"`

	MacAddr net.HardwareAddr `xml:"mac-address,omitempty"`
}

type InterfaceEthernetConfigProcessor interface {
	InterfaceEthernetConfig(string, *InterfaceEthernetConfig) error
}

func NewInterfaceEthernetConfig() *InterfaceEthernetConfig {
	return &InterfaceEthernetConfig{
		SrChanges: nclib.NewSrChanges(),
		MacAddr:   net.HardwareAddr{},
	}
}

func (c *InterfaceEthernetConfig) String() string {
	return fmt.Sprintf("%s{%s='%s'} %s",
		OC_CONFIG_KEY,
		INTERFACE_ETH_MACADDR_KEY, c.MacAddr,
		c.SrChanges,
	)
}

func (c *InterfaceEthernetConfig) SetMacAddr(macAddr net.HardwareAddr) {
	c.MacAddr = macAddr
	c.SetChange(INTERFACE_ETH_MACADDR_KEY)
}

func (c *InterfaceEthernetConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case INTERFACE_ETH_MACADDR_KEY:
		mac, err := net.ParseMAC(value)
		if err != nil {
			return err
		}
		c.MacAddr = mac
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessInterfaceEthernetConfig(p InterfaceEthernetConfigProcessor, reverse bool, name string, config *InterfaceEthernetConfig) error {
	configFunc := func() error {
		return p.InterfaceEthernetConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
