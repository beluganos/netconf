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
	nclib "netconf/lib"
	ncxml "netconf/lib/xml"
)

//
// apply-policy
//
type PolicyApply struct {
	nclib.SrChanges `xml:"-"`

	Config *PolicyApplyConfig `xml:"config"`
}

func NewPolicyApply() *PolicyApply {
	return &PolicyApply{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewPolicyApplyConfig(),
	}
}

func (p *PolicyApply) String() string {
	return fmt.Sprintf("%s{%s} %s",
		POLICYAPPLY_KEY,
		p.Config,
		p.SrChanges,
	)
}

func (p *PolicyApply) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := p.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	p.SetChange(nodes[0].Name)
	return nil
}

//
// apply-policy/config
//
type PolicyApplyConfig struct {
	nclib.SrChanges `xml:"-"`

	ImportPolicy  []string          `xml:"import-policy"`
	ImportDefault PolicyDefaultType `xml:"default-import-policy"`
	ExportPolicy  []string          `xml:"export-policy"`
	ExportDefault PolicyDefaultType `xml:"default-export-policy"`
}

func NewPolicyApplyConfig() *PolicyApplyConfig {
	return &PolicyApplyConfig{
		SrChanges:     nclib.NewSrChanges(),
		ImportPolicy:  []string{},
		ImportDefault: POLICY_DEFAULT_TYPE,
		ExportPolicy:  []string{},
		ExportDefault: POLICY_DEFAULT_TYPE,
	}
}

func (p *PolicyApplyConfig) String() string {
	return fmt.Sprintf("%s{%s=%v, %s=%s, %s=%v, %s=%s} %s",
		OC_CONFIG_KEY,
		POLICYAPPLY_IMPORT_KEY, p.ImportPolicy,
		POLICYAPPLY_IMPORT_DEF_KEY, p.ImportDefault,
		POLICYAPPLY_EXPORT_KEY, p.ExportPolicy,
		POLICYAPPLY_EXPORT_DEF_KEY, p.ExportDefault,
		p.SrChanges,
	)
}

func (p *PolicyApplyConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case POLICYAPPLY_IMPORT_KEY:
		p.ImportPolicy = append(p.ImportPolicy, value)

	case POLICYAPPLY_IMPORT_DEF_KEY:
		t, err := ParsePolicyDefaultType(value)
		if err != nil {
			return err
		}
		p.ImportDefault = t

	case POLICYAPPLY_EXPORT_KEY:
		p.ExportPolicy = append(p.ExportPolicy, value)

	case POLICYAPPLY_EXPORT_DEF_KEY:
		t, err := ParsePolicyDefaultType(value)
		if err != nil {
			return err
		}
		p.ExportDefault = t
	}

	p.SetChange(nodes[0].Name)
	return nil
}
