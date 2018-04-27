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
)

type PolicyPrefixSetMode int

const (
	POLICY_PREFIXSET_MODE PolicyPrefixSetMode = iota
	POLICY_PREFIXSET_MODE_IPV4
	POLICY_PREFIXSET_MODE_IPV6
	POLICY_PREFIXSET_MODE_MIXED
)

var policyPrefixSetModeNames = map[PolicyPrefixSetMode]string{
	POLICY_PREFIXSET_MODE:       "POLICY_PREFIXSET_MODE",
	POLICY_PREFIXSET_MODE_IPV4:  "IPV4",
	POLICY_PREFIXSET_MODE_IPV6:  "IPV6",
	POLICY_PREFIXSET_MODE_MIXED: "MIXED",
}

var policyPrefixSetModeValues = map[string]PolicyPrefixSetMode{
	"POLICY_PREFIXSET_MODE": POLICY_PREFIXSET_MODE,
	"IPV4":                  POLICY_PREFIXSET_MODE_IPV4,
	"IPV6":                  POLICY_PREFIXSET_MODE_IPV6,
	"MIXED":                 POLICY_PREFIXSET_MODE_MIXED,
}

func (v PolicyPrefixSetMode) String() string {
	if s, ok := policyPrefixSetModeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("PolicyPrefixSetMode(%d)", v)
}

func ParsePolicyPrefixSetMode(s string) (PolicyPrefixSetMode, error) {
	if v, ok := policyPrefixSetModeValues[s]; ok {
		return v, nil
	}
	return POLICY_PREFIXSET_MODE, fmt.Errorf("Invalid PolicyPrefixSetMode. %s", s)
}
