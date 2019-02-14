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
// addresses
//
type IPAddresses map[string]*IPAddress

func NewIPAddresses() IPAddresses {
	return IPAddresses{}
}

func (a IPAddresses) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ip, ok := nodes[0].Attrs[SUBINTERFACE_ADDR_IP_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", SUBINTERFACE_ADDR_KEY, SUBINTERFACE_ADDR_IP_KEY, nodes[0])
	}

	ipaddr, ok := a[ip]
	if !ok {
		ipaddr = NewIPAddress(ip)
		a[ip] = ipaddr
	}

	return ipaddr.Put(nodes[1:], value)
}

func (a IPAddresses) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = SUBINTERFACE_ADDRS_KEY
	e.EncodeToken(start)

	for _, addr := range a {
		err := e.EncodeElement(addr, xml.StartElement{Name: xml.Name{Local: SUBINTERFACE_ADDR_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// addresses/address[ip]
//
type IPAddress struct {
	nclib.SrChanges `xml:"-"`

	IP     string           `xml:"ip"`
	Config *IPAddressConfig `xml:"config"`
}

func NewIPAddress(ip string) *IPAddress {
	return &IPAddress{
		SrChanges: nclib.NewSrChanges(),
		IP:        ip,
		Config:    NewIPAddressConfig(),
	}
}

func (a *IPAddress) String() string {
	return fmt.Sprintf("%s{%s=%s, %s} '%s'",
		SUBINTERFACE_ADDR_KEY,
		SUBINTERFACE_ADDR_IP_KEY, a.IP,
		a.Config,
		a.SrChanges,
	)
}

func (a *IPAddress) SetIP(ip string) {
	a.IP = ip
	a.SetChange(SUBINTERFACE_ADDR_IP_KEY)
}

func (a *IPAddress) SetConfig(config *IPAddressConfig) {
	a.Config = config
	a.SetChange(OC_CONFIG_KEY)
}

func (a *IPAddress) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case SUBINTERFACE_ADDR_IP_KEY:
		// a.IP = value // set by NewIPAddress.

	case OC_CONFIG_KEY:
		if err := a.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	a.SetChange(nodes[0].Name)
	return nil
}

//
// addresses/address[ip]/config
//
type IPAddressConfig struct {
	nclib.SrChanges `xml:"-"`

	IP        net.IP `xml:"ip"`
	PrefixLen uint8  `xml:"prefix-length"`
}

func NewIPAddressConfig() *IPAddressConfig {
	return &IPAddressConfig{
		SrChanges: nclib.NewSrChanges(),
		IP:        nil,
		PrefixLen: 0,
	}
}

func (c *IPAddressConfig) String() string {
	return fmt.Sprintf("%s{%s=%s, %s=%d} %s",
		SUBINTERFACE_ADDR_KEY,
		SUBINTERFACE_ADDR_IP_KEY, c.IP,
		SUBINTERFACE_ADDR_PREFIXLEN_KEY, c.PrefixLen,
		c.SrChanges,
	)
}

func (c *IPAddressConfig) IFAddr() *ncnet.IFAddr {
	return ncnet.NewIFAddrWithPlen(c.IP, c.PrefixLen)
}

func (c *IPAddressConfig) SetIP(ip net.IP) {
	c.IP = ip
	c.SetChange(SUBINTERFACE_ADDR_IP_KEY)
}

func (c *IPAddressConfig) SetPLen(plen uint8) {
	c.PrefixLen = plen
	c.SetChange(SUBINTERFACE_ADDR_PREFIXLEN_KEY)
}

func (c *IPAddressConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case SUBINTERFACE_ADDR_IP_KEY:
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("Invalid IP. %s", value)
		}
		c.IP = ip

	case SUBINTERFACE_ADDR_PREFIXLEN_KEY:
		prefixLen, err := strconv.ParseUint(value, 0, 8)
		if err != nil {
			return err
		}
		c.PrefixLen = uint8(prefixLen)

	}

	c.SetChange(nodes[0].Name)
	return nil
}
