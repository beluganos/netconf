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
// routing-policy/policy-definitions/policy-definition[name]/statements/statement[name]/actions
//
type PolicyStatementActions struct {
	nclib.SrChanges `xml:"-"`

	Config *PolicyStatementActionsConfig `xml:"config"`
	Bgp    *PolicyBgpActions             `xml:"bgp-actions"`
}

type PolicyStatementActionsProcessor interface {
	PolicyStatementActionsConfigProcessor
	PolicyBgpActionsProcessor
}

func NewPolicyStatementActions() *PolicyStatementActions {
	return &PolicyStatementActions{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewPolicyStatementActionsConfig(),
		Bgp:       NewPolicyBgpActions(),
	}
}

func (p *PolicyStatementActions) String() string {
	return fmt.Sprintf("%s{%s, %s} %s",
		POLICYDEF_ACTS_KEY,
		p.Config,
		p.Bgp,
		p.SrChanges,
	)
}

func (p *PolicyStatementActions) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := p.Config.Put(nodes[1:], value); err != nil {
			return err
		}

	case BGP_ACTIONS_KEY:
		if err := p.Bgp.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyStatementActions(p PolicyStatementActionsProcessor, reverse bool, name string, stmtName string, actions *PolicyStatementActions) error {
	configFunc := func() error {
		if actions.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyStatementActionsConfig(
				p.(PolicyStatementActionsConfigProcessor),
				reverse,
				name,
				stmtName,
				actions.Config,
			)
		}
		return nil
	}

	bgpFunc := func() error {
		if actions.GetChange(BGP_ACTIONS_KEY) {
			return ProcessPolicyBgpActions(
				p.(PolicyBgpActionsProcessor),
				reverse,
				name,
				stmtName,
				actions.Bgp,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc, bgpFunc)
}

//
// routing-policy/policy-definitions/policy-definition[name]/statements/statement[name]/actions/config
//
type PolicyStatementActionsConfig struct {
	nclib.SrChanges `xml:"-"`

	PolicyResult PolicyResultType `xml:"policy-result"`
}

type PolicyStatementActionsConfigProcessor interface {
	PolicyStatementActionsConfig(string, string, *PolicyStatementActionsConfig) error
}

func NewPolicyStatementActionsConfig() *PolicyStatementActionsConfig {
	return &PolicyStatementActionsConfig{
		SrChanges:    nclib.NewSrChanges(),
		PolicyResult: 0,
	}
}

func (c *PolicyStatementActionsConfig) String() string {
	return fmt.Sprintf("%s{%s=%s} %s",
		OC_CONFIG_KEY,
		POLICYDEF_ACTS_POLRESULT_KEY, c.PolicyResult,
		c.SrChanges,
	)
}

func (p *PolicyStatementActionsConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case POLICYDEF_ACTS_POLRESULT_KEY:
		pr, err := ParsePolicyResultType(value)
		if err != nil {
			return err
		}
		p.PolicyResult = pr
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyStatementActionsConfig(p PolicyStatementActionsConfigProcessor, reverse bool, name string, stmtName string, config *PolicyStatementActionsConfig) error {
	configFunc := func() error {
		return p.PolicyStatementActionsConfig(name, stmtName, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
