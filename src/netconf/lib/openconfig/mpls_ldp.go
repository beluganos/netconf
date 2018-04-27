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
	"netconf/lib"
	"netconf/lib/xml"
)

const (
	MPLS_LDP_ROUTERID_KEY                    = "lsr-id"
	MPLS_LDP_AF_KEY                          = "address-families"
	MPLS_LDP_AF_IPV4_KEY                     = "ipv4"
	MPLS_LDP_AF_IPV6_KEY                     = "ipv6"
	MPLS_LDP_AF_TARNSADDR_KEY                = "transport-address"
	MPLS_LDP_AF_SESSION_HOLDTIME_KEY         = "session-ka-holdtime"
	MPLS_LDP_LABELPOLICY_KEY                 = "label-policy"
	MPLS_LDP_LABELPOLICY_ADOV_KEY            = "advertise"
	MPLS_LDP_LABELPOLICY_ENABLE_KEY          = "enable"
	MPLS_LDP_LABELPOLICY_ENGRESS_EXPNULL_KEY = "egress-explicit-null"
	MPLS_LDP_DISCOVERY_KEY                   = "discovery"
	MPLS_LDP_HELLO_HOLDTIME                  = "hello-holdtime"
	MPLS_LDP_HELLO_INTERVAL                  = "hello-interval"
)

//
// mpls/signaling-protocols/ldp
//
type MplsLdp struct {
	nclib.SrChanges `xml:"-"`

	Global *MplsLdpGlobal `xml:"global"`
}

type MplsLdpProcessor interface {
	MplsLdpGlobalProcessor
}

func NewMplsLdp() *MplsLdp {
	return &MplsLdp{
		SrChanges: nclib.NewSrChanges(),
		Global:    NewMplsLdpGlobal(),
	}
}

func (m *MplsLdp) String() string {
	return fmt.Sprintf("%s{%s} %s",
		MPLS_LDP_KEY,
		m.Global,
		m.SrChanges,
	)
}

func (m *MplsLdp) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_GLOBAL_KEY:
		if err := m.Global.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdp(p MplsLdpProcessor, reverse bool, name string, ldp *MplsLdp) error {
	globalFunc := func() error {
		if ldp.GetChange(OC_GLOBAL_KEY) {
			return ProcessMplsLdpGlobal(
				p.(MplsLdpGlobalProcessor),
				reverse,
				name,
				ldp.Global,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, globalFunc)
}

//
// mpls/signaling-protocols/ldp/global
//
type MplsLdpGlobal struct {
	nclib.SrChanges `xml:"-"`

	Config        *MplsLdpConfig        `xml:"config"`
	AddressFamily *MplsLdpAddressFamily `xml:"address-families"`
	Discovery     *MplsLdpDiscovery     `xml:"discovery"`
}

type MplsLdpGlobalProcessor interface {
	MplsLdpConfigProcessor
	MplsLdpAddressFamilyProcessor
	MplsLdpDiscoveryProcessor
}

func NewMplsLdpGlobal() *MplsLdpGlobal {
	return &MplsLdpGlobal{
		SrChanges:     nclib.NewSrChanges(),
		Config:        NewMplsLdpConfig(),
		AddressFamily: NewMplsLdpAddressFamily(),
		Discovery:     NewMplsLdpDiscovery(),
	}
}

func (m *MplsLdpGlobal) String() string {
	return fmt.Sprintf("%s{%s, %s, %s} %s",
		OC_GLOBAL_KEY,
		m.Config,
		m.AddressFamily,
		m.Discovery,
		m.SrChanges,
	)
}

func (m *MplsLdpGlobal) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := m.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case MPLS_LDP_AF_KEY:
		if err := m.AddressFamily.Put(nodes[1:], value); err != nil {
			return err
		}

	case MPLS_LDP_DISCOVERY_KEY:
		if err := m.Discovery.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil

}

func ProcessMplsLdpGlobal(p MplsLdpGlobalProcessor, reverse bool, name string, global *MplsLdpGlobal) error {

	confgFunc := func() error {
		if global.GetChange(OC_CONFIG_KEY) {
			return ProcessMplsLdpConfig(
				p.(MplsLdpConfigProcessor),
				reverse,
				name,
				global.Config,
			)
		}
		return nil
	}

	afFunc := func() error {
		if global.GetChange(MPLS_LDP_AF_KEY) {
			return ProcessMplsLdpAddressFamily(
				p.(MplsLdpAddressFamilyProcessor),
				reverse,
				name,
				global.AddressFamily,
			)
		}
		return nil
	}

	discFunc := func() error {
		if global.GetChange(MPLS_LDP_DISCOVERY_KEY) {
			return ProcessMplsLdpDiscovery(
				p.(MplsLdpDiscoveryProcessor),
				reverse,
				name,
				global.Discovery,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, confgFunc, afFunc, discFunc)
}

//
// mpls/signaling-protocols/ldp/global/config
//
type MplsLdpConfig struct {
	nclib.SrChanges `xml:"-"`

	LsrId string `xml:"lsr-id"`
}

type MplsLdpConfigProcessor interface {
	MplsLdpConfig(string, *MplsLdpConfig) error
}

func NewMplsLdpConfig() *MplsLdpConfig {
	return &MplsLdpConfig{
		SrChanges: nclib.NewSrChanges(),
		LsrId:     "",
	}
}

func (m *MplsLdpConfig) String() string {
	return fmt.Sprintf("%s{%s='%s'} %s",
		OC_CONFIG_KEY,
		MPLS_LDP_ROUTERID_KEY, m.LsrId,
		m.SrChanges,
	)
}

func (m *MplsLdpConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_ROUTERID_KEY:
		m.LsrId = value
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsLdpConfig(p MplsLdpConfigProcessor, reverse bool, name string, config *MplsLdpConfig) error {
	configFunc := func() error {
		return p.MplsLdpConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
