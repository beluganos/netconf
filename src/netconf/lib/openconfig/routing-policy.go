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

const (
	ROUTINGPOLICY_XMLNS          = "https://github.com/beluganos/beluganos/yang/routing-policy"
	ROUTINGPOLICY_MODULE         = "beluganos-routing-policy"
	ROUTINGPOLICY_KEY            = "routing-policy"
	POLICYDEFS_KEY               = "policy-definitions"
	POLICYDEF_KEY                = "policy-definition"
	POLICYDEF_ACTS_KEY           = "actions"
	POLICYDEF_ACTS_POLRESULT_KEY = "policy-result"
	POLICYDEF_STMTS_KEY          = "statements"
	POLICYDEF_STMT_KEY           = "statement"
	POLICYDEFSETS_KEY            = "defined-sets"
	POLICYPFXSETS_KEY            = "prefix-sets"
	POLICYPFXSET_KEY             = "prefix-set"
	POLICYPFXSET_PREFIXES_KEY    = "prefixes"
	POLICYPFXSET_PREFIX_KEY      = "prefix"
	POLICYPFXSET_PREFIX_IP_KEY   = "ip-prefix"
	POLICYPFXSET_PREFIX_MLR_KEY  = "masklength-range"
	POLICYPFXSET_MODE_KEY        = "mode"
	POLICYNEIGHSETS_KEY          = "neighbor-sets"
	POLICYNEIGHSET_KEY           = "neighbor-set"
	POLICYNEIGHSET_ADDRS_KEY     = "address"
	POLICYTAGSETS_KEY            = "tag-sets"
	POLICYTAGSET_KEY             = "tag-set"
	POLICYTAGSET_TAGVALUE_KEY    = "tag-value"
	POLICYAPPLY_KEY              = "apply-policy"
	POLICYAPPLY_IMPORT_KEY       = "import-policy"
	POLICYAPPLY_IMPORT_DEF_KEY   = "default-import-policy"
	POLICYAPPLY_EXPORT_KEY       = "export-policy"
	POLICYAPPLY_EXPORT_DEF_KEY   = "default-export-policy"
)

//
// routing-policy
//
type RoutingPolicy struct {
	nclib.SrChanges `xml:"-"`

	XMLName     xml.Name           `xml:"https://github.com/beluganos/beluganos/yang/routing-policy routing-policy"`
	DefinedSets *PolicyDefinedSets `xml:"defined-sets"`
	Definitions PolicyDefinitions  `xml:"policy-definitions"`
}

func NewRoutingPolicy() *RoutingPolicy {
	return &RoutingPolicy{
		SrChanges:   nclib.NewSrChanges(),
		DefinedSets: NewPolicyDefinedSets(),
		Definitions: NewPolicyDefinitions(),
	}
}

type RoutingPolicyProcessor interface {
	PolicyDefinedSetsProcessor
	PolicyDefinitionProcessor
}

func (r *RoutingPolicy) String() string {
	return fmt.Sprintf("%s{%s, %s} %s",
		ROUTINGPOLICY_KEY,
		r.DefinedSets,
		r.Definitions,
		r.SrChanges,
	)
}

func (p *RoutingPolicy) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case POLICYDEFSETS_KEY:
		if err := p.DefinedSets.Put(nodes[1:], value); err != nil {
			return err
		}

	case POLICYDEFS_KEY:
		if err := p.Definitions.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	p.SetChange(nodes[0].Name)
	return nil
}

func ProcessRoutingPolicy(p RoutingPolicyProcessor, reverse bool, rpol *RoutingPolicy) error {

	defsetFunc := func() error {
		if rpol.GetChange(POLICYDEFSETS_KEY) {
			return ProcessPolicyDefinedSets(
				p.(PolicyDefinedSetsProcessor),
				reverse,
				rpol.DefinedSets,
			)
		}
		return nil
	}

	defsFunc := func() error {
		if rpol.GetChange(POLICYDEFS_KEY) {
			return ProcessPolicyDefinitions(
				p.(PolicyDefinitionProcessor),
				reverse,
				rpol.Definitions,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, defsetFunc, defsFunc)
}
