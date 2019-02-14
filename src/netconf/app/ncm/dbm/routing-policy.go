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

package ncmdbm

import (
	"fmt"
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"
)

type PolicyDefinitionTable struct {
	session *srlib.SrSession
}

func NewPolicyDefinitionTable(session *srlib.SrSession) *PolicyDefinitionTable {
	return &PolicyDefinitionTable{
		session: session,
	}
}

func (t *PolicyDefinitionTable) Select(name string) (*openconfig.PolicyDefinition, error) {

	xpath := fmt.Sprintf("/%s:%s/%s/%s[%s='%s']//*",
		openconfig.ROUTINGPOLICY_MODULE, openconfig.ROUTINGPOLICY_KEY,
		openconfig.POLICYDEFS_KEY,
		openconfig.POLICYDEF_KEY, openconfig.OC_NAME_KEY, name,
	)

	rpol := openconfig.NewRoutingPolicy()
	for cv := range t.session.GetItems(xpath) {
		cv.Dispatch(rpol, nil, nil)
	}

	poldef, ok := rpol.Definitions[name]
	if !ok {
		return nil, fmt.Errorf("PolicyDefinition not found. %s", name)
	}

	return poldef, nil
}

func (t *PolicyDefinitionTable) Walk(f func(string, *openconfig.PolicyDefinition)) {

	xpath := fmt.Sprintf("/%s:%s/%s//*",
		openconfig.ROUTINGPOLICY_MODULE, openconfig.ROUTINGPOLICY_KEY,
		openconfig.POLICYDEFS_KEY,
	)

	rpol := openconfig.NewRoutingPolicy()
	for cv := range t.session.GetItems(xpath) {
		cv.Dispatch(rpol, nil, nil)
	}

	for name, rpdef := range rpol.Definitions {
		f(name, rpdef)
	}
}

type PolicyStatementTable struct {
	session *srlib.SrSession
}

func NewPolicyStatementTable(session *srlib.SrSession) *PolicyStatementTable {
	return &PolicyStatementTable{
		session: session,
	}
}

func (t *PolicyStatementTable) Select(defName string, stmtName string) (*openconfig.PolicyStatement, error) {
	xpath := fmt.Sprintf("/%s:%s/%s/%s[%s='%s']/%s/%s[%s='%s']//*",
		openconfig.ROUTINGPOLICY_MODULE, openconfig.ROUTINGPOLICY_KEY,
		openconfig.POLICYDEFS_KEY,
		openconfig.POLICYDEF_KEY, openconfig.OC_NAME_KEY, defName,
		openconfig.POLICYDEF_STMTS_KEY,
		openconfig.POLICYDEF_STMT_KEY, openconfig.OC_NAME_KEY, stmtName,
	)

	rpol := openconfig.NewRoutingPolicy()
	for cv := range t.session.GetItems(xpath) {
		cv.Dispatch(rpol, nil, nil)
	}

	poldef, ok := rpol.Definitions[defName]
	if !ok {
		return nil, fmt.Errorf("PolicyDefinition not found. %s", defName)
	}

	stmt, ok := poldef.Stmts[stmtName]
	if !ok {
		return nil, fmt.Errorf("PolicyStatement not found. %s", stmtName)
	}

	return stmt, nil
}
