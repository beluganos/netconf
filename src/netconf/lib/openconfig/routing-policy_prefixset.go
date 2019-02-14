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
	nclib "netconf/lib"
	ncxml "netconf/lib/xml"
)

//
// routing-policy/defined-sets/prefix-sets
//
type PolicyPrefixSets map[string]*PolicyPrefixSet

func NewPolicyPrefixSets() PolicyPrefixSets {
	return PolicyPrefixSets{}
}

func (r PolicyPrefixSets) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	name, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", POLICYPFXSET_KEY, OC_NAME_KEY, nodes[0])
	}

	pfx, ok := r[name]
	if !ok {
		pfx = NewPolicyPrefixSet(name)
		r[name] = pfx
	}

	return pfx.Put(nodes[1:], value)
}

func ProcessPolicyPrefixSets(p PolicyPrefixSetProcessor, reverse bool, prefixes PolicyPrefixSets) error {
	for name, prefix := range prefixes {
		if err := ProcessPolicyPrefixSet(p, reverse, name, prefix); err != nil {
			return err
		}
	}
	return nil
}

func (p PolicyPrefixSets) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = POLICYPFXSETS_KEY
	e.EncodeToken(start)

	for _, pset := range p {
		err := e.EncodeElement(pset, xml.StartElement{Name: xml.Name{Local: POLICYPFXSET_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// routing-policy/defined-sets/prefix-sets/prefix-set[name]
//
type PolicyPrefixSet struct {
	nclib.SrChanges `xml:"-"`

	Name     string                  `xml:"name"`
	Config   *PolicyPrefixSetConfig  `xml:"config"`
	Prefixes PolicyPrefixSetPrefixes `xnl:"prefixes"`
}

type PolicyPrefixSetProcessor interface {
	PolicyPrefixSet(string, *PolicyPrefixSet) error
	PolicyPrefixSetConfigProcessor
	PolicyPrefixSetPrefixProcessor
}

func NewPolicyPrefixSet(name string) *PolicyPrefixSet {
	return &PolicyPrefixSet{
		SrChanges: nclib.NewSrChanges(),
		Name:      name,
		Config:    NewPolicyPrefixSetConfig(),
		Prefixes:  NewPolicyPrefixSetPrefixes(),
	}
}

func (r *PolicyPrefixSet) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s} %s",
		POLICYPFXSET_KEY,
		OC_NAME_KEY, r.Name,
		r.Config,
		r.Prefixes,
		r.SrChanges,
	)
}

func (r *PolicyPrefixSet) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		// r.Name = value // set by NewPolicyPrefixSet

	case OC_CONFIG_KEY:
		if err := r.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case POLICYPFXSET_PREFIXES_KEY:
		if err := r.Prefixes.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	r.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyPrefixSet(p PolicyPrefixSetProcessor, reverse bool, name string, prefix *PolicyPrefixSet) error {
	pfxFunc := func() error {
		if prefix.GetChange(OC_NAME_KEY) {
			return p.PolicyPrefixSet(name, prefix)
		}
		return nil
	}

	configFunc := func() error {
		if prefix.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyPrefixSetConfig(
				p.(PolicyPrefixSetConfigProcessor),
				reverse,
				name,
				prefix.Config,
			)
		}
		return nil
	}

	pfxsFunc := func() error {
		if prefix.GetChange(POLICYPFXSET_PREFIXES_KEY) {
			return ProcessPolicyPrefixSetPrefixes(
				p.(PolicyPrefixSetPrefixProcessor),
				reverse,
				name,
				prefix.Prefixes,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, pfxFunc, configFunc, pfxsFunc)
}

//
// routing-policy/defined-sets/prefix-sets/prefix-set[name]/config
//
type PolicyPrefixSetConfig struct {
	nclib.SrChanges `xnl:"-"`

	Name string              `xml:"name"`
	Mode PolicyPrefixSetMode `xml:"mode"`
}

type PolicyPrefixSetConfigProcessor interface {
	PolicyPrefixSetConfig(string, *PolicyPrefixSetConfig) error
}

func NewPolicyPrefixSetConfig() *PolicyPrefixSetConfig {
	return &PolicyPrefixSetConfig{
		SrChanges: nclib.NewSrChanges(),
		Name:      "",
		Mode:      POLICY_PREFIXSET_MODE,
	}
}

func (r *PolicyPrefixSetConfig) String() string {
	return fmt.Sprintf("%s{%s='%s', %s=%s} %s",
		OC_CONFIG_KEY,
		OC_NAME_KEY, r.Name,
		POLICYPFXSET_MODE_KEY, r.Mode,
		r.SrChanges,
	)
}

func (r *PolicyPrefixSetConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		r.Name = value

	case POLICYPFXSET_MODE_KEY:
		mode, err := ParsePolicyPrefixSetMode(value)
		if err != nil {
			return err
		}
		r.Mode = mode
	}

	r.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyPrefixSetConfig(p PolicyPrefixSetConfigProcessor, reverse bool, name string, config *PolicyPrefixSetConfig) error {
	configFunc := func() error {
		return p.PolicyPrefixSetConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
