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
	ncxml "netconf/lib/xml"
	"strconv"
)

type BgpNeighbors map[string]*BgpNeighbor

func NewBgpNeighbors() BgpNeighbors {
	return BgpNeighbors{}
}

func (b BgpNeighbors) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	addr, ok := nodes[0].Attrs[BGP_NEIGHBOR_ADDR_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", BGP_NEIGHBOR_KEY, BGP_NEIGHBOR_ADDR_KEY, nodes[0])
	}

	neigh, ok := b[addr]
	if !ok {
		neigh = NewBgpNeighbor(addr)
		b[addr] = neigh
	}

	return neigh.Put(nodes[1:], value)
}

func ProcessBgpNeighbors(p BgpNeighborProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, neighs BgpNeighbors) error {
	for addr, neigh := range neighs {
		if err := ProcessBgpNeighbor(p, reverse, name, key, addr, neigh); err != nil {
			return err
		}
	}
	return nil
}

func (b BgpNeighbors) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = BGP_NEIGHBORS_KEY
	e.EncodeToken(start)

	for _, neigh := range b {
		err := e.EncodeElement(neigh, xml.StartElement{Name: xml.Name{Local: BGP_NEIGHBOR_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

type BgpNeighbor struct {
	nclib.SrChanges `xml:"-"`

	Address     string                `xml:"address"`
	Config      *BgpNeighborConfig    `xml:"config"`
	Timers      *BgpNeighborTimers    `xml:"timers"`
	Transport   *BgpNeighborTransport `xml:"transport"`
	AfiSafis    BgpAfiSafis           `xml:"afi-safis"`
	ApplyPolicy *PolicyApply          `xml:"apply-policy"`
}

type BgpNeighborProcessor interface {
	bgpNeighborProcessor
	BgpNeighborConfigProcessor
	BgpNeighborTimersProcessor
	BgpNeighborTransportProcessor
	BgpNeighborAfiSafiProcessor
	BgpNeighborApplyPolicyProcessor
}

type bgpNeighborProcessor interface {
	BgpNeighbor(string, *NetworkInstanceProtocolKey, string, *BgpNeighbor) error
}

func NewBgpNeighbor(addr string) *BgpNeighbor {
	return &BgpNeighbor{
		SrChanges:   nclib.NewSrChanges(),
		Address:     addr,
		Config:      NewBgpNeighborConfig(),
		Timers:      NewBgpNeighborTimers(),
		Transport:   NewBgpNeighborTransport(),
		AfiSafis:    NewBgpAfiSafis(),
		ApplyPolicy: NewPolicyApply(),
	}
}

func (b *BgpNeighbor) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s, %s, %s, %s} %s",
		BGP_NEIGHBOR_KEY,
		BGP_NEIGHBOR_ADDR_KEY, b.Address,
		b.Config,
		b.Timers,
		b.Transport,
		b.AfiSafis,
		b.ApplyPolicy,
		b.SrChanges,
	)
}

func (b *BgpNeighbor) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_NEIGHBOR_ADDR_KEY:
		// b.Address = value // set by NewBgpNeighbor:

	case OC_CONFIG_KEY:
		if err := b.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case BGP_TIMERS_KEY:
		if err := b.Timers.Put(nodes[1:], value); err != nil {
			return err
		}

	case BGP_TRANSPORT_KEY:
		if err := b.Transport.Put(nodes[1:], value); err != nil {
			return err
		}

	case POLICYAPPLY_KEY:
		if err := b.ApplyPolicy.Put(nodes[1:], value); err != nil {
			return err
		}

	case BGP_AFISAFIS_KEY:
		if err := b.AfiSafis.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpNeighbor(p BgpNeighborProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, neigh *BgpNeighbor) error {
	neighFunc := func() error {
		if neigh.GetChange(BGP_NEIGHBOR_ADDR_KEY) {
			return p.BgpNeighbor(name, key, addr, neigh)
		}
		return nil
	}

	configFunc := func() error {
		if neigh.GetChange(OC_CONFIG_KEY) {
			return ProcessBgpNeighborConfig(
				p.(BgpNeighborConfigProcessor),
				reverse,
				name,
				key,
				addr,
				neigh.Config,
			)
		}
		return nil
	}

	timersFunc := func() error {
		if neigh.GetChange(BGP_TIMERS_KEY) {
			return ProcessBgpNeighborTimers(
				p.(BgpNeighborTimersProcessor),
				reverse,
				name,
				key,
				addr,
				neigh.Timers,
			)
		}
		return nil
	}

	transFunc := func() error {
		if neigh.GetChange(BGP_TRANSPORT_KEY) {
			return ProcessBgpNeighborTransport(
				p.(BgpNeighborTransportProcessor),
				reverse,
				name,
				key,
				addr,
				neigh.Transport,
			)
		}
		return nil
	}

	applyPolFunc := func() error {
		if neigh.GetChange(POLICYAPPLY_KEY) {
			return ProcessBgpNeighborApplyPolicy(
				p.(BgpNeighborApplyPolicyProcessor),
				reverse,
				name,
				key,
				addr,
				neigh.ApplyPolicy,
			)
		}
		return nil
	}

	afiSafisFunc := func() error {
		if neigh.GetChange(BGP_AFISAFIS_KEY) {
			return ProcessBgpNeighborAfiSafis(
				p.(BgpNeighborAfiSafiProcessor),
				reverse,
				name,
				key,
				addr,
				neigh.AfiSafis,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, neighFunc, configFunc, timersFunc, transFunc, applyPolFunc, afiSafisFunc)
}

type BgpNeighborConfig struct {
	nclib.SrChanges `xml:"-"`

	Address net.IP `xml:"address"`
	PeerAs  uint32 `xml:"peer-as"`
	LocalAs uint32 `xml:"local-as"`
	Desc    string `xml:"description"`
}

type BgpNeighborConfigProcessor interface {
	BgpNeighborConfig(string, *NetworkInstanceProtocolKey, string, *BgpNeighborConfig) error
}

func NewBgpNeighborConfig() *BgpNeighborConfig {
	return &BgpNeighborConfig{
		SrChanges: nclib.NewSrChanges(),
		Address:   nil,
		PeerAs:    0,
		LocalAs:   0,
		Desc:      "",
	}
}

func (b *BgpNeighborConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s=%d, %s=%d, %s='%s'} %s",
		OC_CONFIG_KEY,
		BGP_NEIGHBOR_ADDR_KEY, b.Address,
		BGP_PEERAS_KEY, b.PeerAs,
		BGP_LOCALAS_KEY, b.LocalAs,
		OC_DESCRIPTION_KEY, b.Desc,
		b.SrChanges,
	)
}

func (b *BgpNeighborConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_NEIGHBOR_ADDR_KEY:
		addr := net.ParseIP(value)
		if addr == nil {
			return fmt.Errorf("Invalid %s. %s", BGP_NEIGHBOR_ADDR_KEY, value)
		}
		b.Address = addr

	case BGP_PEERAS_KEY:
		as, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		b.PeerAs = uint32(as)

	case BGP_LOCALAS_KEY:
		as, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		b.LocalAs = uint32(as)

	case OC_DESCRIPTION_KEY:
		b.Desc = value
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpNeighborConfig(p BgpNeighborConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, config *BgpNeighborConfig) error {
	configFunc := func() error {
		return p.BgpNeighborConfig(name, key, addr, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// bgp/neighbors/neighbor[addr]/timers
//
type BgpNeighborTimers struct {
	nclib.SrChanges `xml:"-"`

	Config *BgpNeighborTimersConfig `xml:"config"`
}

type BgpNeighborTimersProcessor interface {
	BgpNeighborTimersConfigProcessor
}

func NewBgpNeighborTimers() *BgpNeighborTimers {
	return &BgpNeighborTimers{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewBgpNeighborTimersConfig(),
	}
}

func (b *BgpNeighborTimers) String() string {
	return fmt.Sprintf("%s{%s} %s",
		BGP_TIMERS_KEY,
		b.Config,
		b.SrChanges,
	)
}

func (b *BgpNeighborTimers) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := b.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpNeighborTimers(p BgpNeighborTimersProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, timers *BgpNeighborTimers) error {
	timersFunc := func() error {
		if timers.GetChange(OC_CONFIG_KEY) {
			return ProcessBgpNeighborTimersConfig(
				p.(BgpNeighborTimersConfigProcessor),
				reverse,
				name,
				key,
				addr,
				timers.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, timersFunc)
}

type BgpNeighborTimersConfig struct {
	nclib.SrChanges `xml:"-"`

	HoldTime  uint64 `xml:"hold-time"`
	KeepAlive uint64 `xml:"keep-alive"`
}

type BgpNeighborTimersConfigProcessor interface {
	BgpNeighborTimersConfig(string, *NetworkInstanceProtocolKey, string, *BgpNeighborTimersConfig) error
}

func NewBgpNeighborTimersConfig() *BgpNeighborTimersConfig {
	return &BgpNeighborTimersConfig{
		SrChanges: nclib.NewSrChanges(),
		HoldTime:  0,
		KeepAlive: 0,
	}
}

func (b *BgpNeighborTimersConfig) String() string {
	return fmt.Sprintf("%s{%s=%d, %s=%d} %s",
		OC_CONFIG_KEY,
		BGP_HOLDTIME_KEY, b.HoldTime,
		BGP_KEEPALIVE_INTERVAL_KEY, b.KeepAlive,
		b.SrChanges,
	)
}

func (b *BgpNeighborTimersConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_HOLDTIME_KEY:
		t, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return err
		}
		b.HoldTime = t

	case BGP_KEEPALIVE_INTERVAL_KEY:
		t, err := strconv.ParseUint(value, 0, 64)
		if err != nil {
			return err
		}
		b.KeepAlive = t
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpNeighborTimersConfig(p BgpNeighborTimersConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, config *BgpNeighborTimersConfig) error {
	configFunc := func() error {
		return p.BgpNeighborTimersConfig(name, key, addr, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

type BgpNeighborTransport struct {
	nclib.SrChanges `xml:"-"`

	Config *BgpNeighborTransportConfig `xml:"config"`
}

type BgpNeighborTransportProcessor interface {
	BgpNeighborTransportConfigProcessor
}

func NewBgpNeighborTransport() *BgpNeighborTransport {
	return &BgpNeighborTransport{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewBgpNeighborTransportConfig(),
	}
}

func (b *BgpNeighborTransport) String() string {
	return fmt.Sprintf("%s{%s} %s",
		BGP_TRANSPORT_KEY,
		b.Config,
		b.SrChanges,
	)
}

func (b *BgpNeighborTransport) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := b.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpNeighborTransport(p BgpNeighborTransportProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, trans *BgpNeighborTransport) error {
	configFunc := func() error {
		if trans.GetChange(OC_CONFIG_KEY) {
			return ProcessBgpNeighborTransportConfig(
				p.(BgpNeighborTransportConfigProcessor),
				reverse,
				name,
				key,
				addr,
				trans.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

type BgpNeighborTransportConfig struct {
	nclib.SrChanges `xml:"-"`

	LocalAddr net.IP `xml:"local-address"`
}

type BgpNeighborTransportConfigProcessor interface {
	BgpNeighborTransportConfig(string, *NetworkInstanceProtocolKey, string, *BgpNeighborTransportConfig) error
}

func NewBgpNeighborTransportConfig() *BgpNeighborTransportConfig {
	return &BgpNeighborTransportConfig{
		SrChanges: nclib.NewSrChanges(),
		LocalAddr: nil,
	}
}

func (b *BgpNeighborTransportConfig) String() string {
	return fmt.Sprintf("%s{%s=%s} %s",
		OC_CONFIG_KEY,
		BGP_LOCAL_ADDR_KEY, b.LocalAddr,
		b.SrChanges,
	)
}

func (b *BgpNeighborTransportConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_LOCAL_ADDR_KEY:
		addr := net.ParseIP(value)
		if addr == nil {
			return fmt.Errorf("Invaid %s. %s", BGP_LOCAL_ADDR_KEY, value)
		}
		b.LocalAddr = addr
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpNeighborTransportConfig(p BgpNeighborTransportConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, config *BgpNeighborTransportConfig) error {
	configFunc := func() error {
		return p.BgpNeighborTransportConfig(name, key, addr, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// bgp/neighbors/neighbor[addr]/apply-policy
//
type BgpNeighborApplyPolicyProcessor interface {
	BgpNeighborApplyPolicyConfig(string, *NetworkInstanceProtocolKey, string, *PolicyApplyConfig) error
}

func ProcessBgpNeighborApplyPolicy(p BgpNeighborApplyPolicyProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, policy *PolicyApply) error {
	configFunc := func() error {
		if policy.GetChange(OC_CONFIG_KEY) {
			return p.BgpNeighborApplyPolicyConfig(name, key, addr, policy.Config)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// bgp/neighbors/neighbor[addr]/afi-safis
//
func ProcessBgpNeighborAfiSafis(p BgpNeighborAfiSafiProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, afisafis BgpAfiSafis) error {
	for afiSafiName, afisafi := range afisafis {
		if err := ProcessBgpNeighborAfiSafi(p, reverse, name, key, addr, afiSafiName, afisafi); err != nil {
			return err
		}
	}
	return nil
}

//
// bgp/neighbors/neighbor[addr]/afi-safis/afi-safi[afi-safi-name]
//
type BgpNeighborAfiSafiProcessor interface {
	BgpNeighborAfiSafi(string, *NetworkInstanceProtocolKey, string, string, *BgpAfiSafi) error
	BgpNeighborAfiSafiConfig(string, *NetworkInstanceProtocolKey, string, string, *BgpAfiSafiConfig) error
}

func ProcessBgpNeighborAfiSafi(p BgpNeighborAfiSafiProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, addr string, afiSafiName string, afiSafi *BgpAfiSafi) error {
	afisafiFunc := func() error {
		if afiSafi.GetChange(BGP_AFISAFI_NAME_KEY) {
			return p.BgpNeighborAfiSafi(name, key, addr, afiSafiName, afiSafi)
		}
		return nil
	}

	configFunc := func() error {
		if afiSafi.GetChange(OC_CONFIG_KEY) {
			return p.BgpNeighborAfiSafiConfig(name, key, addr, afiSafiName, afiSafi.Config)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, afisafiFunc, configFunc)
}
