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
	"fmt"
	nclib "netconf/lib"
	ncnet "netconf/lib/net"
	ncxml "netconf/lib/xml"
	"strconv"
)

//
// interface-ref
//
type InterfaceRef struct {
	nclib.SrChanges `xml:"-"`

	Config *InterfaceRefConfig `xml:"config"`
}

func NewInterfaceRef() *InterfaceRef {
	return &InterfaceRef{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewInterfaceRefConfig(),
	}
}

func (i *InterfaceRef) String() string {
	return fmt.Sprintf("%s{%s} %s",
		INTERFACE_REF_KEY,
		i.Config,
		i.SrChanges,
	)
}

func (i *InterfaceRef) SetConfig(config *InterfaceRefConfig) {
	i.Config = config
	i.SetChange(OC_CONFIG_KEY)
}

func (i *InterfaceRef) Put(nodes []*ncxml.XPathNode, value string) error {
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

//
// interface-ref/config
//
type InterfaceRefConfig struct {
	nclib.SrChanges `xml:"-"`

	Iface    string `xml:"interface"`
	SubIface uint32 `xml:"subinterface"`
}

func NewInterfaceRefConfig() *InterfaceRefConfig {
	return &InterfaceRefConfig{
		SrChanges: nclib.NewSrChanges(),
		Iface:     "",
		SubIface:  0,
	}
}

func (c *InterfaceRefConfig) IFName() string {
	return ncnet.NewIFName(c.Iface, c.SubIface)
}

func (c *InterfaceRefConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s=%d} %s",
		OC_CONFIG_KEY,
		INTERFACE_KEY, c.Iface,
		SUBINTERFACE_KEY, c.SubIface,
		c.SrChanges,
	)
}

func (c *InterfaceRefConfig) SetIface(iface string) {
	c.Iface = iface
	c.SetChange(INTERFACE_KEY)
}

func (c *InterfaceRefConfig) SetSubIface(subiface uint32) {
	c.SubIface = subiface
	c.SetChange(SUBINTERFACE_KEY)
}

func (c *InterfaceRefConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case INTERFACE_KEY:
		c.Iface = value

	case SUBINTERFACE_KEY:
		subif, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		c.SubIface = uint32(subif)
	}

	c.SetChange(nodes[0].Name)
	return nil
}
