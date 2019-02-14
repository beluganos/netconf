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
// routing-policy/policy-definitions/policy-definition[name]/statements
//
type PolicyStatements map[string]*PolicyStatement

func NewPolicyStatements() PolicyStatements {
	return PolicyStatements{}
}

func (p PolicyStatements) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	name, ok := nodes[0].Attrs[OC_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", POLICYDEF_STMT_KEY, OC_NAME_KEY, nodes[0])
	}

	stmt, ok := p[name]
	if !ok {
		stmt = NewPolicyStatement(name)
		p[name] = stmt
	}

	return stmt.Put(nodes[1:], value)
}

func ProcessPolicyStatements(p PolicyStatementProcessor, reverse bool, name string, stmts PolicyStatements) error {
	for stmtName, stmt := range stmts {
		if err := ProcessPolicyStatement(p, reverse, name, stmtName, stmt); err != nil {
			return err
		}
	}
	return nil
}

func (p PolicyStatements) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = POLICYDEF_STMTS_KEY
	e.EncodeToken(start)
	for _, stmt := range p {
		err := e.EncodeElement(stmt, xml.StartElement{Name: xml.Name{Local: POLICYDEF_STMT_KEY}})
		if err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

//
// routing-policy/policy-definitions/policy-definition[name]/statements/statement[name]
//
type PolicyStatement struct {
	nclib.SrChanges `xml:"-"`

	XMLName xml.Name                `xml:"statement"`
	Name    string                  `xml:"name"`
	Config  *PolicyStatementConfig  `xml:"config"`
	Actions *PolicyStatementActions `xml:"actions"`
}

type PolicyStatementProcessor interface {
	PolicyStatement(string, string, *PolicyStatement) error
	PolicyStatementConfigProcessor
	PolicyStatementActionsProcessor
}

func NewPolicyStatement(name string) *PolicyStatement {
	return &PolicyStatement{
		SrChanges: nclib.NewSrChanges(),
		Name:      name,
		Config:    NewPolicyStatementConfig(),
		Actions:   NewPolicyStatementActions(),
	}
}

func (p *PolicyStatement) String() string {
	return fmt.Sprintf("%s{%s='%s', %s, %s} %s",
		POLICYDEF_STMT_KEY,
		OC_NAME_KEY, p.Name,
		p.Config,
		p.Actions,
		p.SrChanges,
	)
}

func (p *PolicyStatement) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_NAME_KEY:
		//p.Name = value // set by NewPolicyStatement

	case OC_CONFIG_KEY:
		if err := p.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case POLICYDEF_ACTS_KEY:
		if err := p.Actions.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyStatement(p PolicyStatementProcessor, reverse bool, name string, stmtName string, stmt *PolicyStatement) error {
	nameFunc := func() error {
		if stmt.GetChange(OC_NAME_KEY) {
			return p.PolicyStatement(name, stmtName, stmt)
		}
		return nil
	}

	configFunc := func() error {
		if stmt.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyStatementConfig(
				p.(PolicyStatementConfigProcessor),
				reverse,
				name,
				stmtName,
				stmt.Config,
			)
		}
		return nil
	}

	actsFunc := func() error {
		if stmt.GetChange(POLICYDEF_ACTS_KEY) {
			return ProcessPolicyStatementActions(
				p.(PolicyStatementActionsProcessor),
				reverse,
				name,
				stmtName,
				stmt.Actions,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, nameFunc, configFunc, actsFunc)
}

//
// routing-policy/policy-definitions/policy-definition[name]/statements/statement[name]/config
//
type PolicyStatementConfig struct {
	nclib.SrChanges `xml:"-"`

	Name string `xml:"name"`
}

type PolicyStatementConfigProcessor interface {
	PolicyStatementConfig(string, string, *PolicyStatementConfig) error
}

func NewPolicyStatementConfig() *PolicyStatementConfig {
	return &PolicyStatementConfig{
		SrChanges: nclib.NewSrChanges(),
		Name:      "",
	}
}

func (p *PolicyStatementConfig) String() string {
	return fmt.Sprintf("%s{%s='%s'} %s",
		OC_CONFIG_KEY,
		OC_NAME_KEY, p.Name,
		p.SrChanges,
	)
}

func (p *PolicyStatementConfig) Put(nodes []*ncxml.XPathNode, value string) error {
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

func ProcessPolicyStatementConfig(p PolicyStatementConfigProcessor, reverse bool, name string, stmtName string, config *PolicyStatementConfig) error {
	configFunc := func() error {
		return p.PolicyStatementConfig(name, stmtName, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
