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

type PolicyPrefixSetPrefixKey struct {
	IpPrefix     string
	MaskLenRange string
}

func NewPolicyPrefixSetPrefixKey(ipPrefix, maskLenRange string) *PolicyPrefixSetPrefixKey {
	return &PolicyPrefixSetPrefixKey{
		IpPrefix:     ipPrefix,
		MaskLenRange: maskLenRange,
	}
}

//
// routing-policy/defined-sets/prefix-sets/prefix-set[name]/prefixes
//
type PolicyPrefixSetPrefixes map[PolicyPrefixSetPrefixKey]*PolicyPrefixSetPrefix

func NewPolicyPrefixSetPrefixes() PolicyPrefixSetPrefixes {
	return PolicyPrefixSetPrefixes{}
}

func (r PolicyPrefixSetPrefixes) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	ip, ok := nodes[0].Attrs[POLICYPFXSET_PREFIX_IP_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", POLICYPFXSET_PREFIX_KEY, POLICYPFXSET_PREFIX_IP_KEY, nodes[0])
	}

	mlrange, ok := nodes[0].Attrs[POLICYPFXSET_PREFIX_MLR_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", POLICYPFXSET_PREFIX_KEY, POLICYPFXSET_PREFIX_MLR_KEY, nodes[0])
	}

	key := NewPolicyPrefixSetPrefixKey(ip, mlrange)
	prefix, ok := r[*key]
	if !ok {
		prefix = NewPolicyPrefixSetPrefix(key)
		r[*key] = prefix
	}

	return prefix.Put(nodes[1:], value)
}

func ProcessPolicyPrefixSetPrefixes(p PolicyPrefixSetPrefixProcessor, reverse bool, name string, prefixes PolicyPrefixSetPrefixes) error {
	for key, prefix := range prefixes {
		if err := ProcessPolicyPrefixSetPrefix(p, reverse, name, &key, prefix); err != nil {
			return err
		}
	}
	return nil
}

func (p PolicyPrefixSetPrefixes) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = POLICYPFXSET_PREFIXES_KEY
	e.EncodeToken(start)

	for _, pset := range p {
		err := e.EncodeElement(pset, xml.StartElement{Name: xml.Name{Local: POLICYPFXSET_PREFIX_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// routing-policy/defined-sets/prefix-sets/prefix-set[name]/prefixes/prefix[prefix, range]
//
type PolicyPrefixSetPrefix struct {
	nclib.SrChanges `xml:"-"`

	IpPrefix     string                       `xml:"ip-prefix"`
	MaskLenRange string                       `xml:"masklength-range"`
	Config       *PolicyPrefixSetPrefixConfig `xml:"config"`
}

type PolicyPrefixSetPrefixProcessor interface {
	PolicyPrefixSetPrefix(string, *PolicyPrefixSetPrefixKey, *PolicyPrefixSetPrefix) error
	PolicyPrefixSetPrefixConfigProcessor
}

func NewPolicyPrefixSetPrefix(key *PolicyPrefixSetPrefixKey) *PolicyPrefixSetPrefix {
	return &PolicyPrefixSetPrefix{
		SrChanges:    nclib.NewSrChanges(),
		IpPrefix:     key.IpPrefix,
		MaskLenRange: key.MaskLenRange,
		Config:       NewPolicyPrefixSetPrefixConfig(),
	}
}

func (p *PolicyPrefixSetPrefix) String() string {
	return fmt.Sprintf("%s{%s='%s', %s='%s', %s} %s",
		POLICYPFXSET_PREFIX_KEY,
		POLICYPFXSET_PREFIX_IP_KEY, p.IpPrefix,
		POLICYPFXSET_PREFIX_MLR_KEY, p.MaskLenRange,
		p.Config,
		p.SrChanges,
	)
}

func (p *PolicyPrefixSetPrefix) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case POLICYPFXSET_PREFIX_IP_KEY:
		// r.IpPrexix = value // set by NewPolicyPrefixSetPrefix

	case POLICYPFXSET_PREFIX_MLR_KEY:
		// r.MaskLenRange = value // set by NewPolicyPrefixSetPrefix

	case OC_CONFIG_KEY:
		if err := p.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyPrefixSetPrefix(p PolicyPrefixSetPrefixProcessor, reverse bool, name string, pfxkey *PolicyPrefixSetPrefixKey, prefix *PolicyPrefixSetPrefix) error {
	pfxFunc := func() error {
		if prefix.GetChanges(POLICYPFXSET_PREFIX_IP_KEY, POLICYPFXSET_PREFIX_MLR_KEY) {
			return p.PolicyPrefixSetPrefix(name, pfxkey, prefix)
		}
		return nil
	}

	configFunc := func() error {
		if prefix.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyPrefixSetPrefixConfig(
				p.(PolicyPrefixSetPrefixConfigProcessor),
				reverse,
				name,
				pfxkey,
				prefix.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, pfxFunc, configFunc)
}

//
// routing-policy/defined-sets/prefix-sets/prefix-set[name]/prefixes/prefix[prefix, range]/config
//
type PolicyPrefixSetPrefixConfig struct {
	nclib.SrChanges `xml:"-"`

	IpPrefix     *net.IPNet `xml:"ip-prefix"`
	MaskLenRange string     `xml:"masklength-range"`
}

type PolicyPrefixSetPrefixConfigProcessor interface {
	PolicyPrefixSetPrefixConfig(string, *PolicyPrefixSetPrefixKey, *PolicyPrefixSetPrefixConfig) error
}

func NewPolicyPrefixSetPrefixConfig() *PolicyPrefixSetPrefixConfig {
	return &PolicyPrefixSetPrefixConfig{
		SrChanges:    nclib.NewSrChanges(),
		IpPrefix:     nil,
		MaskLenRange: "",
	}
}

func (p *PolicyPrefixSetPrefixConfig) String() string {
	return fmt.Sprintf("%s{%s=%s, %s='%s'} %s",
		OC_CONFIG_KEY,
		POLICYPFXSET_PREFIX_IP_KEY, p.IpPrefix,
		POLICYPFXSET_PREFIX_MLR_KEY, p.MaskLenRange,
		p.SrChanges,
	)
}

func (p *PolicyPrefixSetPrefixConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case POLICYPFXSET_PREFIX_IP_KEY:
		_, nw, err := net.ParseCIDR(value)
		if err != nil {
			return err
		}
		p.IpPrefix = nw

	case POLICYPFXSET_PREFIX_MLR_KEY:
		p.MaskLenRange = value
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyPrefixSetPrefixConfig(p PolicyPrefixSetPrefixConfigProcessor, reverse bool, name string, key *PolicyPrefixSetPrefixKey, config *PolicyPrefixSetPrefixConfig) error {
	configFunc := func() error {
		return p.PolicyPrefixSetPrefixConfig(name, key, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
