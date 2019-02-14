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
	"net"
	nclib "netconf/lib"
	ncxml "netconf/lib/xml"
	"strconv"
)

//
// mpls/signaling-protocols/ldp/global/address-family
//
type MplsLdpAddressFamily struct {
	nclib.SrChanges `xml:"-"`

	IPv4 *MplsLdpAddressFamilyV4 `xml:"ipv4"`
	// IPv6 *MplsLdpAddressFamilyV6 `xml:"ipv6"`
}

type MplsLdpAddressFamilyProcessor interface {
	MplsLdpAddressFamilyV4Processor
}

func NewMplsLdpAddressFamily() *MplsLdpAddressFamily {
	return &MplsLdpAddressFamily{
		SrChanges: nclib.NewSrChanges(),
		IPv4:      NewMplsLdpAddressFamilyV4(),
		// IPv6: NewMplsLdpAddressFamilyV6(),
	}
}

func (m *MplsLdpAddressFamily) String() string {
	return fmt.Sprintf("%s{%s} %s",
		MPLS_LDP_AF_KEY,
		m.IPv4,
		m.SrChanges,
	)
}

func (m *MplsLdpAddressFamily) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_AF_IPV4_KEY:
		if err := m.IPv4.Put(nodes[1:], value); err != nil {
			return err
		}

		//case MPLS_LDP_AF_IPV6_KEY:
		//	if err := m.IPv6.Put(nodes[1:], value); err != nil {
		//		return err
		//	}

	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpAddressFamily(p MplsLdpAddressFamilyProcessor, reverse bool, name string, af *MplsLdpAddressFamily) error {
	ipv4Func := func() error {
		if af.GetChange(MPLS_LDP_AF_IPV4_KEY) {
			return ProcessMplsLdpAddressFamilyV4(
				p.(MplsLdpAddressFamilyV4Processor),
				reverse,
				name,
				af.IPv4,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, ipv4Func)
}

//
// mpls/signaling-protocols/ldp/address-family/ipv4
//
type MplsLdpAddressFamilyV4 struct {
	nclib.SrChanges `xml:"-"`

	Config *MplsLdpAddressFamilyV4Config `xml:"config"`
}

type MplsLdpAddressFamilyV4Processor interface {
	MplsLdpAddressFamilyV4ConfigProcessor
}

func NewMplsLdpAddressFamilyV4() *MplsLdpAddressFamilyV4 {
	return &MplsLdpAddressFamilyV4{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewMplsLdpAddressFamilyV4Config(),
	}
}

func (m *MplsLdpAddressFamilyV4) String() string {
	return fmt.Sprintf("%s{%v} %s",
		MPLS_LDP_AF_IPV4_KEY,
		m.Config,
		m.SrChanges,
	)
}

func (m *MplsLdpAddressFamilyV4) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := m.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpAddressFamilyV4(p MplsLdpAddressFamilyV4Processor, reverse bool, name string, ipv4 *MplsLdpAddressFamilyV4) error {
	configFunc := func() error {
		if ipv4.GetChange(OC_CONFIG_KEY) {
			return ProcessMplsLdpAddressFamilyV4Config(
				p.(MplsLdpAddressFamilyV4ConfigProcessor),
				reverse,
				name,
				ipv4.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// mpls/signaling-protocols/ldp/address-family/ipv4/config
//
type MplsLdpAddressFamilyV4Config struct {
	nclib.SrChanges `xml:"-"`

	TransportAddr   net.IP                `xml:"transport-address"`
	SessionHoldTime uint16                `xml:"session-ka-holdtime"`
	LabelPolicy     *MplsLdpV4LabelPolicy `xml:"label-policy"`
}

type MplsLdpAddressFamilyV4ConfigProcessor interface {
	MplsLdpAddressFamilyV4Config(string, *MplsLdpAddressFamilyV4Config) error
}

func NewMplsLdpAddressFamilyV4Config() *MplsLdpAddressFamilyV4Config {
	return &MplsLdpAddressFamilyV4Config{
		SrChanges:       nclib.NewSrChanges(),
		TransportAddr:   nil,
		SessionHoldTime: 0,
		LabelPolicy:     NewMplsLdpV4LabelPolicy(),
	}
}

func (m *MplsLdpAddressFamilyV4Config) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_AF_TARNSADDR_KEY:
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("Invalid %s, %s", MPLS_LDP_AF_TARNSADDR_KEY, value)
		}
		m.TransportAddr = ip

	case MPLS_LDP_AF_SESSION_HOLDTIME_KEY:
		ht, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}
		m.SessionHoldTime = uint16(ht)

	case MPLS_LDP_LABELPOLICY_KEY:
		if err := m.LabelPolicy.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpAddressFamilyV4Config(p MplsLdpAddressFamilyV4ConfigProcessor, reverse bool, name string, config *MplsLdpAddressFamilyV4Config) error {
	configFunc := func() error {
		return p.MplsLdpAddressFamilyV4Config(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// mpls/signaling-protocols/ldp/global/address-families/ipv4/config/label-policy
//
type MplsLdpV4LabelPolicy struct {
	nclib.SrChanges `xml:"-"`

	Advertise *MplsLdpV4Advertise `xml:"advertise"`
}

func NewMplsLdpV4LabelPolicy() *MplsLdpV4LabelPolicy {
	return &MplsLdpV4LabelPolicy{
		SrChanges: nclib.NewSrChanges(),
		Advertise: NewMplsLdpV4Advertise(),
	}
}

func (m *MplsLdpV4LabelPolicy) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_LABELPOLICY_ADOV_KEY:
		if err := m.Advertise.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

//
// mpls/signaling-protocols/ldp/global/address-families/ipv4/config/label-policy/advertise
//
type MplsLdpV4Advertise struct {
	nclib.SrChanges `xml:"-"`

	EngressExplicitNull *MplsLdpEngressExplicitNull `xml:"egress-explicit-null"`
}

func NewMplsLdpV4Advertise() *MplsLdpV4Advertise {
	return &MplsLdpV4Advertise{
		SrChanges:           nclib.NewSrChanges(),
		EngressExplicitNull: NewMplsLdpEngressExplicitNull(),
	}
}

func (m *MplsLdpV4Advertise) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_LABELPOLICY_ENGRESS_EXPNULL_KEY:
		if err := m.EngressExplicitNull.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

//
// mpls/signaling-protocols/ldp/global/address-families/ipv4/config/label-policy/advertise/egress-explicit-null
//
type MplsLdpEngressExplicitNull struct {
	nclib.SrChanges `xml:"-"`

	Enable bool `xml:"enable"`
}

func NewMplsLdpEngressExplicitNull() *MplsLdpEngressExplicitNull {
	return &MplsLdpEngressExplicitNull{
		SrChanges: nclib.NewSrChanges(),
		Enable:    true,
	}
}

func (m *MplsLdpEngressExplicitNull) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_LABELPOLICY_ENABLE_KEY:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return nil
		}
		m.Enable = b
	}

	m.SetChange(nodes[0].Name)
	return nil
}
