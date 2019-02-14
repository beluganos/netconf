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
)

//
// routing-policy/defined-sets/neighbor-sets
//
type PolicyNeighborSets map[string]*PolicyNeighborSet

func NewPolicyNeighborSets() PolicyNeighborSets {
	return PolicyNeighborSets{}
}

func (r PolicyNeighborSets) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	name, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", POLICYNEIGHSET_KEY, OC_NAME_KEY, nodes[0])
	}

	neigh, ok := r[name]
	if !ok {
		neigh = NewPolicyNeighborSet(name)
		r[name] = neigh
	}

	return neigh.Put(nodes[1:], value)
}

func ProcessPolicyNeighborSets(p PolicyNeighborSetProcessor, reverse bool, neighs PolicyNeighborSets) error {
	for name, neigh := range neighs {
		if err := ProcessPolicyNeighborSet(p, reverse, name, neigh); err != nil {
			return err
		}
	}
	return nil
}

func (p PolicyNeighborSets) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = POLICYNEIGHSETS_KEY
	e.EncodeToken(start)
	for _, nset := range p {
		err := e.EncodeElement(nset, xml.StartElement{Name: xml.Name{Local: POLICYNEIGHSET_KEY}})
		if err != nil {
			return err
		}

	}
	return e.EncodeToken(start.End())
}

//
// routing-policy/defined-sets/neighbor-sets/neighbor-set[name]
//
type PolicyNeighborSet struct {
	nclib.SrChanges `xml:"-"`

	Name   string                   `xml:"name"`
	Config *PolicyNeighborSetConfig `xml:"config"`
}

type PolicyNeighborSetProcessor interface {
	PolicyNeighborSet(string, *PolicyNeighborSet) error
	PolicyNeighborSetConfigProcessor
}

func NewPolicyNeighborSet(name string) *PolicyNeighborSet {
	return &PolicyNeighborSet{
		SrChanges: nclib.NewSrChanges(),
		Name:      name,
		Config:    NewPolicyNeighborSetConfig(),
	}
}

func (r *PolicyNeighborSet) String() string {
	return fmt.Sprintf("%s{%s='%s', %s} %s",
		POLICYNEIGHSET_KEY,
		OC_NAME_KEY, r.Name,
		r.Config,
		r.SrChanges,
	)
}

func (r *PolicyNeighborSet) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		// r.Name = value // set by NewPolicyNeighborSet

	case OC_CONFIG_KEY:
		if err := r.Config.Put(nodes[1:], value); err != nil {
			return nil
		}
	}

	r.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyNeighborSet(p PolicyNeighborSetProcessor, reverse bool, name string, neigh *PolicyNeighborSet) error {
	neighFunc := func() error {
		if neigh.GetChange(OC_NAME_KEY) {
			return p.PolicyNeighborSet(name, neigh)
		}
		return nil
	}

	configFunc := func() error {
		if neigh.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyNeighborSetConfig(
				p.(PolicyNeighborSetConfigProcessor),
				reverse,
				name,
				neigh.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, neighFunc, configFunc)
}

//
// routing-policy/defined-sets/neighbor-sets/neighbor-set[name]/config
//
type PolicyNeighborSetConfig struct {
	nclib.SrChanges `xml:"-"`

	Name  string   `xml:"name"`
	Addrs []net.IP `xml:"address"`
}

type PolicyNeighborSetConfigProcessor interface {
	PolicyNeighborSetConfig(string, *PolicyNeighborSetConfig) error
}

func NewPolicyNeighborSetConfig() *PolicyNeighborSetConfig {
	return &PolicyNeighborSetConfig{
		SrChanges: nclib.NewSrChanges(),
		Name:      "",
		Addrs:     []net.IP{},
	}
}

func (r *PolicyNeighborSetConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s=%v} %s",
		OC_CONFIG_KEY,
		OC_NAME_KEY, r.Name,
		POLICYNEIGHSET_ADDRS_KEY, r.Addrs,
		r.SrChanges,
	)
}

func (r *PolicyNeighborSetConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		r.Name = value

	case POLICYNEIGHSET_ADDRS_KEY:
		ip := net.ParseIP(value)
		if ip == nil {
			return fmt.Errorf("Invalid %s. %s", POLICYNEIGHSET_ADDRS_KEY, value)
		}
		r.Addrs = append(r.Addrs, ip)
	}

	r.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyNeighborSetConfig(p PolicyNeighborSetConfigProcessor, reverse bool, name string, config *PolicyNeighborSetConfig) error {
	configFunc := func() error {
		return p.PolicyNeighborSetConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
