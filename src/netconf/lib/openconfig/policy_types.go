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
	"netconf/lib/xml"
)

type InstallProtocolType int

const (
	INSTALL_PROTOCOL_TYPE InstallProtocolType = iota
	INSTALL_PROTOCOL_BGP
	INSTALL_PROTOCOL_ISIS
	INSTALL_PROTOCOL_OSPF
	INSTALL_PROTOCOL_OSPF3
	INSTALL_PROTOCOL_STATIC
	INSTALL_PROTOCOL_DIRECTLY_CONNECTED
	INSTALL_PROTOCOL_LOCAL_AGGREGATE
)

var installProtocolTypeNames = map[InstallProtocolType]string{
	INSTALL_PROTOCOL_TYPE:               "INSTALL_PROTOCOL_TYPE",
	INSTALL_PROTOCOL_BGP:                "BGP",
	INSTALL_PROTOCOL_ISIS:               "ISIS",
	INSTALL_PROTOCOL_OSPF:               "OSPF",
	INSTALL_PROTOCOL_OSPF3:              "OSPF3",
	INSTALL_PROTOCOL_STATIC:             "STATIC",
	INSTALL_PROTOCOL_DIRECTLY_CONNECTED: "DIRECTLY_CONNECTED",
	INSTALL_PROTOCOL_LOCAL_AGGREGATE:    "LOCAL_AGGREGATE",
}

var installProtocolTypeValues = map[string]InstallProtocolType{
	"INSTALL_PROTOCOL_TYPE": INSTALL_PROTOCOL_TYPE,
	"BGP":                   INSTALL_PROTOCOL_BGP,
	"ISIS":                  INSTALL_PROTOCOL_ISIS,
	"OSPF":                  INSTALL_PROTOCOL_OSPF,
	"OSPF3":                 INSTALL_PROTOCOL_OSPF3,
	"STATIC":                INSTALL_PROTOCOL_STATIC,
	"DIRECTLY_CONNECTED":    INSTALL_PROTOCOL_DIRECTLY_CONNECTED,
	"LOCAL_AGGREGATE":       INSTALL_PROTOCOL_LOCAL_AGGREGATE,
}

func (v InstallProtocolType) String() string {
	if s, ok := installProtocolTypeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("InstallProtocolType(%d)", v)
}

func ParseInstallProtocolType(s string) (InstallProtocolType, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := installProtocolTypeValues[ss]; ok {
		return v, nil
	}
	return INSTALL_PROTOCOL_TYPE, fmt.Errorf("Invalid InstallProtocolType. %s", s)
}

type PolicyMatchSetOptionsType int

const (
	POLICY_MATCH_SET_OPTIONS_TYPE PolicyMatchSetOptionsType = iota
	POLICY_MATCH_SET_OPTIONS_ANY
	POLICY_MATCH_SET_OPTIONS_ALL
	POLICY_MATCH_SET_OPTIONS_INVERT
)

var policyMatchSetOptionsTypeNames = map[PolicyMatchSetOptionsType]string{
	POLICY_MATCH_SET_OPTIONS_TYPE:   "POLICY_MATCH_SET_OPTIONS_TYPE",
	POLICY_MATCH_SET_OPTIONS_ANY:    "ANY",
	POLICY_MATCH_SET_OPTIONS_ALL:    "ALL",
	POLICY_MATCH_SET_OPTIONS_INVERT: "INVERT",
}

var policyMatchSetOptionsTypeValues = map[string]PolicyMatchSetOptionsType{
	"POLICY_MATCH_SET_OPTIONS_TYPE": POLICY_MATCH_SET_OPTIONS_TYPE,
	"ANY":                           POLICY_MATCH_SET_OPTIONS_ANY,
	"ALL":                           POLICY_MATCH_SET_OPTIONS_ALL,
	"INVERT":                        POLICY_MATCH_SET_OPTIONS_INVERT,
}

func (v PolicyMatchSetOptionsType) String() string {
	if s, ok := policyMatchSetOptionsTypeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("PolicyMatchSetOptionsType(%d)", v)
}

func ParsePolicyMatchSetOptionsType(s string) (PolicyMatchSetOptionsType, error) {
	if v, ok := policyMatchSetOptionsTypeValues[s]; ok {
		return v, nil
	}
	return POLICY_MATCH_SET_OPTIONS_TYPE, fmt.Errorf("Invalid PolicyMatchSetOptionsType. %s", s)
}

type PolicyResultType int

const (
	POLICY_RESULT_TYPE PolicyResultType = iota
	POLICY_RESULT_ACCEPT_ROUTE
	POLICY_RESULT_REJECT_ROUTE
)

var policyResultTypeNames = map[PolicyResultType]string{
	POLICY_RESULT_TYPE:         "POLICY_RESULT_TYPE",
	POLICY_RESULT_ACCEPT_ROUTE: "ACCEPT_ROUTE",
	POLICY_RESULT_REJECT_ROUTE: "REJECT_ROUTE",
}

var policyResultTypeValues = map[string]PolicyResultType{
	"POLICY_RESULT_TYPE": POLICY_RESULT_TYPE,
	"ACCEPT_ROUTE":       POLICY_RESULT_ACCEPT_ROUTE,
	"REJECT_ROUTE":       POLICY_RESULT_REJECT_ROUTE,
}

func (v PolicyResultType) String() string {
	if s, ok := policyResultTypeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("PolicyResultType(%d)", v)
}

func ParsePolicyResultType(s string) (PolicyResultType, error) {
	if v, ok := policyResultTypeValues[s]; ok {
		return v, nil
	}
	return POLICY_RESULT_TYPE, fmt.Errorf("Invalid PolicyResultType. %s", s)
}

//
// default-policy-type
//
type PolicyDefaultType int

const (
	POLICY_DEFAULT_TYPE PolicyDefaultType = iota
	POLICY_DEFAULT_ACCEPT_ROUTE
	POLICY_DEFAULT_REJECT_ROUTE
)

var policyDefaultTypeNames = map[PolicyDefaultType]string{
	POLICY_DEFAULT_TYPE:         "POLICY_DEFAULT_TYPE",
	POLICY_DEFAULT_ACCEPT_ROUTE: "ACCEPT_ROUTE",
	POLICY_DEFAULT_REJECT_ROUTE: "REJECT_ROUTE",
}

var policyDefaultTypeValues = map[string]PolicyDefaultType{
	"POLICY_DEFAULT_TYPE": POLICY_DEFAULT_TYPE,
	"ACCEPT_ROUTE":        POLICY_DEFAULT_ACCEPT_ROUTE,
	"REJECT_ROUTE":        POLICY_DEFAULT_REJECT_ROUTE,
}

func (v PolicyDefaultType) String() string {
	if s, ok := policyDefaultTypeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("PolicyDefaultType(%d)", v)
}

func ParsePolicyDefaultType(s string) (PolicyDefaultType, error) {
	if v, ok := policyDefaultTypeValues[s]; ok {
		return v, nil
	}
	return POLICY_DEFAULT_TYPE, fmt.Errorf("Invalid PolicyDefaultType. %s", s)
}
