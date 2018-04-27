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
	"netconf/lib"
	"netconf/lib/xml"
	"strconv"
)

//
// routing-policy/defined-sets/tag-sets
//
type PolicyTagSets map[string]*PolicyTagSet

func NewPolicyTagSets() PolicyTagSets {
	return PolicyTagSets{}
}

func (r PolicyTagSets) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	name, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", POLICYTAGSET_KEY, OC_NAME_KEY, nodes[0])
	}

	tag, ok := r[name]
	if !ok {
		tag = NewPolicyTagSet(name)
		r[name] = tag
	}

	return tag.Put(nodes[1:], value)
}

func ProcessPolicyTagSets(p PolicyTagSetProcessor, reverse bool, tags PolicyTagSets) error {
	for name, tag := range tags {
		if err := ProcessPolicyTagSet(p, reverse, name, tag); err != nil {
			return err
		}
	}
	return nil
}

func (p PolicyTagSets) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = POLICYTAGSETS_KEY
	e.EncodeToken(start)

	for _, tset := range p {
		err := e.EncodeElement(tset, xml.StartElement{Name: xml.Name{Local: POLICYTAGSET_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// routing-policy/defined-sets/tag-sets/tag-set[name]
//
type PolicyTagSet struct {
	nclib.SrChanges `xml:"-"`

	Name   string              `xml:"name"`
	Config *PolicyTagSetConfig `xml:"config"`
}

type PolicyTagSetProcessor interface {
	PolicyTagSet(string, *PolicyTagSet) error
	PolicyTagSetConfigProcessor
}

func NewPolicyTagSet(name string) *PolicyTagSet {
	return &PolicyTagSet{
		SrChanges: nclib.NewSrChanges(),
		Name:      name,
		Config:    NewPolicyTagSetConfig(),
	}
}

func (r *PolicyTagSet) String() string {
	return fmt.Sprintf("%s{%s='%s', %s} %s",
		POLICYTAGSET_KEY,
		OC_NAME_KEY, r.Name,
		r.Config,
		r.SrChanges,
	)
}

func (r *PolicyTagSet) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		// r.Name = value // set by NewPolicyTagSet

	case OC_CONFIG_KEY:
		if err := r.Config.Put(nodes[1:], value); err != nil {
			return nil
		}
	}

	r.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyTagSet(p PolicyTagSetProcessor, reverse bool, name string, tag *PolicyTagSet) error {
	tagFunc := func() error {
		if tag.GetChange(OC_NAME_KEY) {
			return p.PolicyTagSet(name, tag)
		}
		return nil
	}

	configFunc := func() error {
		if tag.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyTagSetConfig(
				p.(PolicyTagSetConfigProcessor),
				reverse,
				name,
				tag.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, tagFunc, configFunc)
}

//
// routing-policy/defined-sets/tag-sets/tag-set[name]/config
//
type PolicyTagSetConfig struct {
	nclib.SrChanges `xml:"-"`

	Name     string `xml:"name"`
	TagValue uint32 `xml:"tag-value"`
}

type PolicyTagSetConfigProcessor interface {
	PolicyTagSetConfig(string, *PolicyTagSetConfig) error
}

func NewPolicyTagSetConfig() *PolicyTagSetConfig {
	return &PolicyTagSetConfig{
		SrChanges: nclib.NewSrChanges(),
		Name:      "",
		TagValue:  0,
	}
}

func (r *PolicyTagSetConfig) String() string {
	return fmt.Sprintf("%s{%s='%s'} %s",
		OC_CONFIG_KEY,
		OC_NAME_KEY, r.Name,
		r.SrChanges,
	)
}

func (r *PolicyTagSetConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		r.Name = value

	case POLICYTAGSET_TAGVALUE_KEY:
		value, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		r.TagValue = uint32(value)
	}

	r.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyTagSetConfig(p PolicyTagSetConfigProcessor, reverse bool, name string, config *PolicyTagSetConfig) error {
	configFunc := func() error {
		return p.PolicyTagSetConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
