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

//
// routing-policy/defined-sets
//
type PolicyDefinedSets struct {
	nclib.SrChanges `xml:"-"`

	PrefixSets   PolicyPrefixSets   `xml:"-"`
	NeighborSets PolicyNeighborSets `xml:"-"`
	TagSets      PolicyTagSets      `xml:"-"`
}

type PolicyDefinedSetsProcessor interface {
	PolicyPrefixSetProcessor
	PolicyNeighborSetProcessor
	PolicyTagSetProcessor
}

func NewPolicyDefinedSets() *PolicyDefinedSets {
	return &PolicyDefinedSets{
		SrChanges:    nclib.NewSrChanges(),
		PrefixSets:   NewPolicyPrefixSets(),
		NeighborSets: NewPolicyNeighborSets(),
		TagSets:      NewPolicyTagSets(),
	}
}

func (d *PolicyDefinedSets) String() string {
	return fmt.Sprintf("%s{%s, %s, %s} %s",
		POLICYDEFSETS_KEY,
		d.PrefixSets,
		d.NeighborSets,
		d.TagSets,
		d.SrChanges,
	)
}

func (d *PolicyDefinedSets) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case POLICYPFXSETS_KEY:
		if err := d.PrefixSets.Put(nodes[1:], value); err != nil {
			return err
		}

	case POLICYNEIGHSETS_KEY:
		if err := d.NeighborSets.Put(nodes[1:], value); err != nil {
			return err
		}

	case POLICYTAGSETS_KEY:
		if err := d.TagSets.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	d.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyDefinedSets(p PolicyDefinedSetsProcessor, reverse bool, polsets *PolicyDefinedSets) error {
	pfxFunc := func() error {
		if polsets.GetChange(POLICYPFXSETS_KEY) {
			return ProcessPolicyPrefixSets(
				p.(PolicyPrefixSetProcessor),
				reverse,
				polsets.PrefixSets,
			)
		}
		return nil
	}

	neiFunc := func() error {
		if polsets.GetChange(POLICYNEIGHSETS_KEY) {
			return ProcessPolicyNeighborSets(
				p.(PolicyNeighborSetProcessor),
				reverse,
				polsets.NeighborSets,
			)
		}
		return nil
	}

	tagFunc := func() error {
		if polsets.GetChange(POLICYTAGSETS_KEY) {
			return ProcessPolicyTagSets(
				p.(PolicyTagSetProcessor),
				reverse,
				polsets.TagSets,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, pfxFunc, neiFunc, tagFunc)
}
