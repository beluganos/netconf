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

package ncgobgp

type PolicyDefinition Entries

func NewPolicyDefinition(i interface{}) PolicyDefinition {
	return PolicyDefinition(NewEntries(i))
}

func NewPolicyDefinitions(i interface{}) []PolicyDefinition {
	defs := []PolicyDefinition{}

	switch i.(type) {
	case nil:
	default:
		for _, d := range i.([]interface{}) {
			defs = append(defs, NewPolicyDefinition(d))
		}
	}

	return defs
}

func RawPolicyDefinitions(polDefs []PolicyDefinition) interface{} {
	list := make([]interface{}, len(polDefs))
	for index, polDef := range polDefs {
		list[index] = Entries(polDef).Raw()
	}
	return list
}

func SelectPolicyDefinition(i interface{}, name string) (PolicyDefinition, int) {
	switch i.(type) {
	case nil:
	default:
		for index, d := range i.([]interface{}) {
			if pdef := NewPolicyDefinition(d); pdef.Name() == name {
				return pdef, index
			}
		}
	}

	return nil, -1
}

func (p PolicyDefinition) Name() string {
	return convString(p, "name")
}

func (p PolicyDefinition) SetName(name string) {
	p["name"] = name
}

func (p PolicyDefinition) Statements() []Statement {
	return NewStatements(getValue(p, "statements"))
}

func (p PolicyDefinition) SetStatements(stmts []Statement) {
	p["statements"] = RawStatements(stmts)
}
