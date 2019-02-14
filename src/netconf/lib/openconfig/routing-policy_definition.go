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
// routing-policy/policy-definitions
//
type PolicyDefinitions map[string]*PolicyDefinition

func NewPolicyDefinitions() PolicyDefinitions {
	return PolicyDefinitions{}
}

func (p PolicyDefinitions) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	name, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", POLICYDEF_KEY, OC_NAME_KEY, nodes[0])
	}

	pol, ok := p[name]
	if !ok {
		pol = NewPolicyDefinition(name)
		p[name] = pol
	}

	return pol.Put(nodes[1:], value)
}

func ProcessPolicyDefinitions(p PolicyDefinitionProcessor, reverse bool, poldefs PolicyDefinitions) error {
	for name, poldef := range poldefs {
		if err := ProcessPolicyDefinition(p, reverse, name, poldef); err != nil {
			return err
		}
	}
	return nil
}

func (p PolicyDefinitions) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = POLICYDEFS_KEY
	e.EncodeToken(start)
	for _, rpol := range p {
		err := e.EncodeElement(rpol, xml.StartElement{Name: xml.Name{Local: POLICYDEF_KEY}})
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

//
// routing-policy/policy-definitions/policy-definition[name]
//
type PolicyDefinition struct {
	nclib.SrChanges `xml:"-"`

	Name   string                  `xml:"name"`
	Config *PolicyDefinitionConfig `xml:"config"`
	Stmts  PolicyStatements        `xml:"statements"`
}

type PolicyDefinitionProcessor interface {
	PolicyDefinition(string, *PolicyDefinition) error
	PolicyDefinitionConfigProcessor
	PolicyStatementProcessor
}

func NewPolicyDefinition(name string) *PolicyDefinition {
	return &PolicyDefinition{
		SrChanges: nclib.NewSrChanges(),
		Name:      name,
		Config:    NewPolicyDefinitionConfig(),
		Stmts:     NewPolicyStatements(),
	}
}

func (p *PolicyDefinition) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s} %s",
		POLICYDEF_KEY,
		OC_NAME_KEY, p.Name,
		p.Config,
		p.Stmts,
		p.SrChanges,
	)
}

func (p *PolicyDefinition) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		// p.Name = value // set by NewPolicyDefinition

	case OC_CONFIG_KEY:
		if err := p.Config.Put(nodes[1:], value); err != nil {
			return nil
		}

	case POLICYDEF_STMTS_KEY:
		if err := p.Stmts.Put(nodes[1:], value); err != nil {
			return nil
		}
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyDefinition(p PolicyDefinitionProcessor, reverse bool, name string, poldef *PolicyDefinition) error {

	nameFunc := func() error {
		if poldef.GetChange(OC_NAME_KEY) {
			return p.PolicyDefinition(name, poldef)
		}
		return nil
	}

	configFunc := func() error {
		if poldef.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyDefinitionConfig(
				p.(PolicyDefinitionConfigProcessor),
				reverse,
				name,
				poldef.Config,
			)
		}
		return nil
	}

	stmtFunc := func() error {
		if poldef.GetChange(POLICYDEF_STMTS_KEY) {
			return ProcessPolicyStatements(
				p.(PolicyStatementProcessor),
				reverse,
				name,
				poldef.Stmts,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, nameFunc, configFunc, stmtFunc)
}

//
// routing-policy/policy-definitions/policy-definition[name]/config
//
type PolicyDefinitionConfig struct {
	nclib.SrChanges `xml:"-"`

	Name string `xml:"name"`
}

type PolicyDefinitionConfigProcessor interface {
	PolicyDefinitionConfig(string, *PolicyDefinitionConfig) error
}

func NewPolicyDefinitionConfig() *PolicyDefinitionConfig {
	return &PolicyDefinitionConfig{
		SrChanges: nclib.NewSrChanges(),
		Name:      "",
	}
}

func (c *PolicyDefinitionConfig) String() string {
	return fmt.Sprintf("%s{%s='%s'} %s",
		OC_CONFIG_KEY,
		OC_NAME_KEY, c.Name,
		c.SrChanges,
	)
}

func (p *PolicyDefinitionConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		p.Name = value
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyDefinitionConfig(p PolicyDefinitionConfigProcessor, reverse bool, name string, config *PolicyDefinitionConfig) error {

	configFunc := func() error {
		return p.PolicyDefinitionConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
