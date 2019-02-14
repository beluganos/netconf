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
	ncxml "netconf/lib/xml"
)

const IANAifTypeNS = "iana-if-type"

const (
	OC_CONFIG_KEY      = "config"
	OC_STATS_KEY       = "stats"
	OC_ID_KEY          = "id"
	OC_NAME_KEY        = "name"
	OC_GLOBAL_KEY      = "global"
	OC_INDEX_KEY       = "index"
	OC_TYPE_KEY        = "type"
	OC_ENABLED_KEY     = "enabled"
	OC_IDENT_KEY       = "identifier"
	OC_DESCRIPTION_KEY = "description"
)

type AddressFamily int

const (
	ADDRESS_FAMILY AddressFamily = iota
	ADDRESS_FAMILY_IPV4
	ADDRESS_FAMILY_IPV6
	ADDRESS_FAMILY_MPLS
	ADDRESS_FAMILY_L2_ETHERNET
)

var addressFamily_names = map[AddressFamily]string{
	ADDRESS_FAMILY:             "ADDRESS_FAMILY",
	ADDRESS_FAMILY_IPV4:        "IPV4",
	ADDRESS_FAMILY_IPV6:        "IPV6",
	ADDRESS_FAMILY_MPLS:        "MPLS",
	ADDRESS_FAMILY_L2_ETHERNET: "L2_ETHERNET",
}

var addressFamily_values = map[string]AddressFamily{
	"ADDRESS_FAMILY": ADDRESS_FAMILY,
	"IPV4":           ADDRESS_FAMILY_IPV4,
	"IPV6":           ADDRESS_FAMILY_IPV6,
	"MPLS":           ADDRESS_FAMILY_MPLS,
	"L2_ETHERNET":    ADDRESS_FAMILY_L2_ETHERNET,
}

func (v AddressFamily) String() string {
	if s, ok := addressFamily_names[v]; ok {
		return s
	}
	return fmt.Sprintf("AddressFamily(%d)", v)
}

func ParseAddressFamily(s string) (AddressFamily, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := addressFamily_values[ss]; ok {
		return v, nil
	}
	return ADDRESS_FAMILY, fmt.Errorf("Invalid Address Family. %s", s)
}
