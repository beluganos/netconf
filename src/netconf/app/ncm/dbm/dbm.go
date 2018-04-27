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
	"netconf/lib/sysrepo"
)

//
// Tables
//
type Tables struct {
	session *srlib.SrSession
	ifaces  *InterfaceTable
	subifs  *SubinterfaceTable
	defs    *PolicyDefinitionTable
	stmts   *PolicyStatementTable
}

func NewTables(session *srlib.SrSession) *Tables {
	return &Tables{
		session: session,
		ifaces:  NewInterfaceTable(session),
		subifs:  NewSubinterfaceTable(session),
		defs:    NewPolicyDefinitionTable(session),
		stmts:   NewPolicyStatementTable(session),
	}
}

func (t *Tables) Interface() *InterfaceTable {
	return t.ifaces
}

func (t *Tables) Subinterface() *SubinterfaceTable {
	return t.subifs
}

func (t *Tables) PolicyDefinitions() *PolicyDefinitionTable {
	return t.defs
}

func (t *Tables) PolicyStatements() *PolicyStatementTable {
	return t.stmts
}

func (t *Tables) Refresh() error {
	return t.session.Refresh()
}

var tables *Tables = nil

func Create(session *srlib.SrSession) {
	tables = NewTables(session)
}

func Refresh() error {
	return tables.Refresh()
}

func Interfaces() *InterfaceTable {
	return tables.Interface()
}

func Subinterfaces() *SubinterfaceTable {
	return tables.Subinterface()
}

func PolicyDefinitions() *PolicyDefinitionTable {
	return tables.PolicyDefinitions()
}

func PolicyStatements() *PolicyStatementTable {
	return tables.PolicyStatements()
}
