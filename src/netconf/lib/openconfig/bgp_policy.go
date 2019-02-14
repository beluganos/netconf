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
	"strconv"
)

const (
	BGP_ACTIONS_KEY               = "bgp-actions"
	BGP_ACTIONS_SET_LOCALPREF_KEY = "set-local-pref"
	BGP_ACTIONS_SET_NEXTHOP_KEY   = "set-next-hop"
)

//
// bgp-policy-actions
//
type PolicyBgpActions struct {
	nclib.SrChanges `xml:"-"`

	Config *PolicyBgpActionsConfig `xml:"config"`
}

func (b *PolicyBgpActions) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Space = BGP_POLICY_XMLNS
	e.EncodeToken(start)
	if err := e.EncodeElement(b.Config, xml.StartElement{Name: xml.Name{Local: OC_CONFIG_KEY}}); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

type PolicyBgpActionsProcessor interface {
	PolicyBgpActionsConfigProcessor
}

func NewPolicyBgpActions() *PolicyBgpActions {
	return &PolicyBgpActions{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewPolicyBgpActionsConfig(),
	}
}

func (b *PolicyBgpActions) String() string {
	return fmt.Sprintf("%s{%s} %s",
		BGP_ACTIONS_KEY,
		b.Config,
		b.SrChanges,
	)
}

func (b *PolicyBgpActions) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := b.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyBgpActions(p PolicyBgpActionsProcessor, reverse bool, name string, stmtName string, actions *PolicyBgpActions) error {
	configFunc := func() error {
		if actions.GetChange(OC_CONFIG_KEY) {
			return ProcessPolicyBgpActionsConfig(
				p.(PolicyBgpActionsConfigProcessor),
				reverse,
				name,
				stmtName,
				actions.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

//
// bgp-policy-actions/config
//
type PolicyBgpActionsConfig struct {
	nclib.SrChanges `xml:"-"`

	SetLocalPref uint32         `xml:"set-local-pref"`
	SetNexthop   BgpNexthopType `xml:"set-next-hop"`
}

type PolicyBgpActionsConfigProcessor interface {
	PolicyBgpActionsConfig(string, string, *PolicyBgpActionsConfig) error
}

func NewPolicyBgpActionsConfig() *PolicyBgpActionsConfig {
	return &PolicyBgpActionsConfig{
		SrChanges:    nclib.NewSrChanges(),
		SetLocalPref: 0,
		SetNexthop:   BGP_NEXTHOP_TYPE,
	}
}

func (c *PolicyBgpActionsConfig) String() string {
	return fmt.Sprintf("%s{%s=%d, %s='%s'} %s",
		OC_CONFIG_KEY,
		BGP_ACTIONS_SET_LOCALPREF_KEY, c.SetLocalPref,
		BGP_ACTIONS_SET_NEXTHOP_KEY, c.SetNexthop,
		c.SrChanges,
	)
}

func (c *PolicyBgpActionsConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_ACTIONS_SET_LOCALPREF_KEY:
		lpref, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		c.SetLocalPref = uint32(lpref)

	case BGP_ACTIONS_SET_NEXTHOP_KEY:
		nh, err := ParseBgpNexthopType(value)
		if err != nil {
			return err
		}
		c.SetNexthop = nh
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessPolicyBgpActionsConfig(p PolicyBgpActionsConfigProcessor, reverse bool, name string, stmtName string, config *PolicyBgpActionsConfig) error {
	configFunc := func() error {
		return p.PolicyBgpActionsConfig(name, stmtName, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
