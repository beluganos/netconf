// -*- coding: utf-8 -*-

// Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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
	srlib "netconf/lib/sysrepo"
	"testing"
)

func makeBgpNeighbors(datas [][2]string) (BgpNeighbors, error) {
	neigh := NewBgpNeighbors()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := neigh.Put(nodes[1:], value); err != nil {
			return neigh, err
		}
	}

	return neigh, nil
}

func TestBgpNeighbors_v4(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/neighbor-address", "10.0.0.1"},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["10.0.0.1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_NEIGHBOR_ADDR_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}
}

func TestBgpNeighbors_v6(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']", ""},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/neighbor-address", "2001:2001::1"},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["2001:2001::1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_NEIGHBOR_ADDR_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}
}

func TestBgpNeighbors_config(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/neighbor-address", "10.0.0.1"},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/config", ""},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["10.0.0.1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_NEIGHBOR_ADDR_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}
	if v := neigh.Config.Compare(); !v {
		t.Errorf("BgpNeighbors.Put unmatch. config cmp=%t", v)
	}
}

func TestBgpNeighbors_config_v4(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/neighbor-address", "10.0.0.1"},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/config", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/config/neighbor-address", "10.0.0.1"},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["10.0.0.1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_NEIGHBOR_ADDR_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}
	if v := neigh.Config.Compare(BGP_NEIGHBOR_ADDR_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. config cmp=%t", v)
	}
}

func TestBgpNeighbors_config_all(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']", ""},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/neighbor-address", "2001:2001::1"},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/config", ""},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/config/neighbor-address", "2001:2001::1"},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/config/peer-as", "65000"},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/config/local-as", "65001"},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["2001:2001::1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_NEIGHBOR_ADDR_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}
	if v := neigh.Config.Compare(BGP_NEIGHBOR_ADDR_KEY, BGP_PEERAS_KEY, BGP_LOCALAS_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. config cmp=%t", v)
	}
}

func TestBgpNeighbors_afisafi_ipv4(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV4_UNICAST']", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV4_UNICAST']/afi-safi-name", "ox-bgp-types:IPV4_UNICAST"},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["10.0.0.1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_AFISAFIS_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}

	afisafi := neigh.AfiSafis["oc-bgp-types:IPV4_UNICAST"]
	t.Log(afisafi)

	if v := afisafi.Compare(BGP_AFISAFI_NAME_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. afisafi cmp=%t", v)
	}
}

func TestBgpNeighbors_afisafi_ipv6(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']", ""},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/afi-safis", ""},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV6_UNICAST']", ""},
		{"/neighbors/neighbor[neighbor-address='2001:2001::1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV6_UNICAST']/afi-safi-name", "ox-bgp-types:IPV6_UNICAST"},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["2001:2001::1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_AFISAFIS_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}

	afisafi := neigh.AfiSafis["oc-bgp-types:IPV6_UNICAST"]
	t.Log(afisafi)

	if v := afisafi.Compare(BGP_AFISAFI_NAME_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. afisafi cmp=%t", v)
	}
}

func TestBgpNeighbors_afisafi_config(t *testing.T) {
	neighs, err := makeBgpNeighbors([][2]string{
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV4_UNICAST']", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV4_UNICAST']/afi-safi-name", "oc-bgp-types:IPV4_UNICAST"},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV4_UNICAST']/config", ""},
		{"/neighbors/neighbor[neighbor-address='10.0.0.1']/afi-safis/afi-safi[afi-safi-name='oc-bgp-types:IPV4_UNICAST']/config/afi-safi-name", "oc-bgp-types:IPV4_UNICAST"},
	})

	if err != nil {
		t.Errorf("BgpNeighbors.Put error. %s", err)
	}

	neigh := neighs["10.0.0.1"]
	t.Log(neigh)

	if v := neigh.Compare(BGP_AFISAFIS_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. cmp=%t", v)
	}

	afisafi := neigh.AfiSafis["oc-bgp-types:IPV4_UNICAST"]
	t.Log(afisafi)

	if v := afisafi.Compare(BGP_AFISAFI_NAME_KEY, OC_CONFIG_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. afisafi cmp=%t", v)
	}
	if v := afisafi.Config.Compare(BGP_AFISAFI_NAME_KEY); !v {
		t.Errorf("BgpNeighbors.Put unmatch. afisafi.config cmp=%t", v)
	}
}
