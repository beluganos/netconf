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

type OspfNetworkType int

const (
	OSPF_NETWORK_TYPE OspfNetworkType = iota
	OSPF_POINT_TO_POINT_NETWORK
	OSPF_BROADCAST_NETWORK
	OSPF_NON_BROADCAST_NETWORK
)

var ospfNetworkTypeNames = map[OspfNetworkType]string{
	OSPF_NETWORK_TYPE:           "OSPF_NETWORK_TYPE",
	OSPF_POINT_TO_POINT_NETWORK: "POINT_TO_POINT_NETWORK",
	OSPF_BROADCAST_NETWORK:      "BROADCAST_NETWORK",
	OSPF_NON_BROADCAST_NETWORK:  "NON_BROADCAST_NETWORK",
}

var ospfNetworkTypeValues = map[string]OspfNetworkType{
	"OSPF_NETWORK_TYPE":      OSPF_NETWORK_TYPE,
	"POINT_TO_POINT_NETWORK": OSPF_POINT_TO_POINT_NETWORK,
	"BROADCAST_NETWORK":      OSPF_BROADCAST_NETWORK,
	"NON_BROADCAST_NETWORK":  OSPF_NON_BROADCAST_NETWORK,
}

func (v OspfNetworkType) String() string {
	if s, ok := ospfNetworkTypeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("OspfNetworkType(%d)", v)
}

func ParseOspfNetworkType(s string) (OspfNetworkType, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := ospfNetworkTypeValues[ss]; ok {
		return v, nil
	}
	return OSPF_NETWORK_TYPE, fmt.Errorf("Invalid OspfNetworkType. %s", s)
}
