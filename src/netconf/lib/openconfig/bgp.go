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
	BGP_KEY                    = "bgp"
	BGP_AS_KEY                 = "as"
	BGP_PEERAS_KEY             = "peer-as"
	BGP_LOCALAS_KEY            = "local-as"
	BGP_ROUTERID_KEY           = "router-id"
	BGP_NEIGHBORS_KEY          = "neighbors"
	BGP_NEIGHBOR_KEY           = "neighbor"
	BGP_NEIGHBOR_ADDR_KEY      = "neighbor-address"
	BGP_TIMERS_KEY             = "timers"
	BGP_TRANSPORT_KEY          = "transport"
	BGP_HOLDTIME_KEY           = "hold-time"
	BGP_KEEPALIVE_INTERVAL_KEY = "keepalive-interval"
	BGP_LOCAL_ADDR_KEY         = "local-address"
	BGP_ZEBRA_KEY              = "zebra"
	BGP_ZEBRA_VERSION_KEY      = "version"
	BGP_ZEBRA_URL_KEY          = "url"
	BGP_ZEBRA_REDISTROUTES_KEY = "redistribute-routes"
	BGP_AFISAFIS_KEY           = "afi-safis"
	BGP_AFISAFI_KEY            = "afi-safi"
	BGP_AFISAFI_NAME_KEY       = "afi-safi-name"
	BGP_APPLYPOLICY_KEY        = "apply-policy"
	BGP_TYPES_XMLNS            = "http://openconfig.net/yang/bgp-types"
	BGP_TYPES_MODULE           = "oc-bgp-types"
	BGP_POLICY_XMLNS           = "https://github.com/beluganos/beluganos/yang/bgp-policy"
	BGP_POLICY_MODULE          = "boc-bgp-pol"
)

//
// bgp
//
type Bgp struct {
	nclib.SrChanges `xml:"-"`

	Global    *BgpGlobal   `xml:"global"`
	Zebra     *BgpZebra    `xml:"zebra"`
	Neighbors BgpNeighbors `xml:"neighbors"`
}

type BgpProcessor interface {
	bgpProcessor
	BgpGlobalProcessor
	BgpZebraProcessor
	BgpNeighborProcessor
}

type bgpProcessor interface {
	Bgp(string, *NetworkInstanceProtocolKey, *Bgp) error
}

func NewBgp() *Bgp {
	return &Bgp{
		SrChanges: nclib.NewSrChanges(),
		Global:    NewBgpGlobal(),
		Zebra:     NewBgpZebra(),
		Neighbors: NewBgpNeighbors(),
	}
}

func (b *Bgp) String() string {
	return fmt.Sprintf("%s{%s, %s} %s",
		BGP_KEY,
		b.Global,
		b.Zebra,
		b.SrChanges,
	)
}

func (b *Bgp) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_GLOBAL_KEY:
		if err := b.Global.Put(nodes[1:], value); err != nil {
			return err
		}

	case BGP_ZEBRA_KEY:
		if err := b.Zebra.Put(nodes[1:], value); err != nil {
			return err
		}

	case BGP_NEIGHBORS_KEY:
		if err := b.Neighbors.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgp(p BgpProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, bgp *Bgp) error {
	bgpFunc := func() error {
		return p.Bgp(name, key, bgp)
	}

	globalFunc := func() error {
		if bgp.GetChange(OC_GLOBAL_KEY) {
			return ProcessBgpGlobal(
				p.(BgpGlobalProcessor),
				reverse,
				name,
				key,
				bgp.Global,
			)
		}
		return nil
	}

	zebraFunc := func() error {
		if bgp.GetChange(BGP_ZEBRA_KEY) {
			return ProcessBgpZebra(
				p.(BgpZebraProcessor),
				reverse,
				name,
				key,
				bgp.Zebra,
			)
		}
		return nil
	}

	neighFunc := func() error {
		if bgp.GetChange(BGP_NEIGHBORS_KEY) {
			return ProcessBgpNeighbors(
				p.(BgpNeighborProcessor),
				reverse,
				name,
				key,
				bgp.Neighbors,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, bgpFunc, globalFunc, zebraFunc, neighFunc)
}
