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

package srocgobgp

import (
	"fmt"
	"netconf/lib/openconfig"
	"strings"
)

func QStringList(names []string) string {
	qnames := make([]string, len(names))
	for index, name := range names {
		qnames[index] = fmt.Sprintf("\"%s\"", name)
	}
	return fmt.Sprintf("[%s]", strings.Join(qnames, ","))
}

func QString(i interface{}) string {
	return fmt.Sprintf("\"%v\"", i)
}

var bgpAfiSafiType = map[openconfig.BgpAfiSafiType]string{
	openconfig.BGP_AFI_SAFI_TYPE:                 "",
	openconfig.BGP_AFI_SAFI_IPV4_UNICAST:         "ipv4-unicast",
	openconfig.BGP_AFI_SAFI_IPV6_UNICAST:         "ipv6-unicast",
	openconfig.BGP_AFI_SAFI_IPV4_LABELED_UNICAST: "ipv4-labelled-unicast",
	openconfig.BGP_AFI_SAFI_IPV6_LABELED_UNICAST: "ipv6-labelled-unicast",
	openconfig.BGP_AFI_SAFI_L3VPN_IPV4_UNICAST:   "l3vpn-ipv4-unicast",
	openconfig.BGP_AFI_SAFI_L3VPN_IPV6_UNICAST:   "l3vpn-ipv6-unicast",
	openconfig.BGP_AFI_SAFI_L3VPN_IPV4_MULTICAST: "l3vpn-ipv4-multicast",
	openconfig.BGP_AFI_SAFI_L3VPN_IPV6_MULTICAST: "l3vpn-ipv6-multicast",
	openconfig.BGP_AFI_SAFI_L2VPN_VPLS:           "l2vpn-vpls",
	openconfig.BGP_AFI_SAFI_L2VPN_EVPN:           "l2vpn-evpn",
}

func BgpAfiSafiType(t openconfig.BgpAfiSafiType) string {
	if s, ok := bgpAfiSafiType[t]; ok {
		return s
	}
	return fmt.Sprintf("Invalid bgpAfiSafiType(%d)", t)
}

var installProtocolTypes = map[openconfig.InstallProtocolType]string{
	openconfig.INSTALL_PROTOCOL_BGP:                "bgp",
	openconfig.INSTALL_PROTOCOL_ISIS:               "isis",
	openconfig.INSTALL_PROTOCOL_OSPF:               "ospf",
	openconfig.INSTALL_PROTOCOL_OSPF3:              "ospfv3",
	openconfig.INSTALL_PROTOCOL_STATIC:             "static",
	openconfig.INSTALL_PROTOCOL_DIRECTLY_CONNECTED: "connected",
	openconfig.INSTALL_PROTOCOL_LOCAL_AGGREGATE:    "aggregate",
}

func InstallProtocolType(t openconfig.InstallProtocolType) string {
	if s, ok := installProtocolTypes[t]; ok {
		return s
	}
	return fmt.Sprintf("InstallProtocolType(%d)", t)
}

func InstallProtocolTypes(types []openconfig.InstallProtocolType) []string {
	ss := make([]string, len(types))
	for index, t := range types {
		ss[index] = InstallProtocolType(t)
	}
	return ss
}

var policyResultTypes = map[openconfig.PolicyResultType]string{
	openconfig.POLICY_RESULT_ACCEPT_ROUTE: "accept-route",
	openconfig.POLICY_RESULT_REJECT_ROUTE: "reject-route",
}

func PolicyResultType(t openconfig.PolicyResultType) string {
	if s, ok := policyResultTypes[t]; ok {
		return s
	}
	return fmt.Sprintf("PolicyResultType(%d)", t)
}

var policyDefaultTypes = map[openconfig.PolicyDefaultType]string{
	openconfig.POLICY_DEFAULT_ACCEPT_ROUTE: "accept-route",
	openconfig.POLICY_DEFAULT_REJECT_ROUTE: "reject-route",
}

func PolicyDefaultType(t openconfig.PolicyDefaultType) string {
	if s, ok := policyDefaultTypes[t]; ok {
		return s
	}
	return fmt.Sprintf("PolicyDefaultType(%d)", t)
}
