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

package main

import (
	"fmt"
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"
)

type InterfaceVals struct {
	Vals []*srlib.SrVal
}

func NewInterfaceVals() *InterfaceVals {
	return &InterfaceVals{
		Vals: []*srlib.SrVal{},
	}
}

func (i *InterfaceVals) append(v *srlib.SrVal) {
	i.Vals = append(i.Vals, v)
}

func (i *InterfaceVals) interfaceXPath(name string) string {
	return fmt.Sprintf("/%s:%s/%s[%s='%s']",
		openconfig.INTERFACES_MODULE,
		openconfig.INTERFACES_KEY,
		openconfig.INTERFACE_KEY,
		openconfig.OC_NAME_KEY, name,
	)
}

func (i *InterfaceVals) interfaceEthXPath(name string) string {
	return fmt.Sprintf("%s/%s:%s",
		i.interfaceXPath(name),
		openconfig.INTERFACE_ETH_MODULE,
		openconfig.INTERFACE_ETH_KEY,
	)
}

/*
func (i *InterfaceVals) subinterfaceXPath(name string, index uint32) string {
	return fmt.Sprintf("%s/%s/%s[%s='%d']",
		i.interfaceXPath(name),
		openconfig.SUBINTERFACES_KEY,
		openconfig.SUBINTERFACE_KEY,
		openconfig.OC_INDEX_KEY, index,
	)
}

func (i *InterfaceVals) subinterfaceIPv4XPath(name string, index uint32) string {
	return fmt.Sprintf("%s/%s:%s",
		i.subinterfaceXPath(name, index),
		openconfig.SUBINTERFACE_IP_MODULE,
		openconfig.SUBINTERFACE_IPV4_KEY,
	)
}

func (i *InterfaceVals) subinterfaceIPv4AddrXPath(name string, index uint32, ip string) string {
	return fmt.Sprintf("%s/%s/%s[%s='%s']",
		i.subinterfaceIPv4XPath(name, index),
		openconfig.SUBINTERFACE_ADDRS_KEY,
		openconfig.SUBINTERFACE_ADDR_KEY,
		openconfig.SUBINTERFACE_ADDR_IP_KEY, ip,
	)
}
*/

func (i *InterfaceVals) Interface(name string, iface *openconfig.Interface) error {
	return nil
}

func (i *InterfaceVals) InterfaceConfig(name string, config *openconfig.InterfaceConfig) error {
	base := fmt.Sprintf("%s/%s", i.interfaceXPath(name), openconfig.OC_CONFIG_KEY)

	if config.GetChange(openconfig.OC_NAME_KEY) {
		xpath := fmt.Sprintf("%s/%s", base, openconfig.OC_NAME_KEY)
		i.append(srlib.ParseSrVal(config.Name, false, srlib.SR_STRING_T, xpath))
	}

	if config.GetChange(openconfig.INTERFACE_TYPE_KEY) {
		xpath := fmt.Sprintf("%s/%s", base, openconfig.INTERFACE_TYPE_KEY)
		i.append(srlib.ParseSrVal(config.Type, false, srlib.SR_IDENTITYREF_T, xpath))
	}

	if config.GetChange(openconfig.OC_ENABLED_KEY) {
		xpath := fmt.Sprintf("%s/%s", base, openconfig.OC_ENABLED_KEY)
		i.append(srlib.ParseSrVal(config.Enabled, false, srlib.SR_BOOL_T, xpath))
	}

	if config.GetChange(openconfig.INTERFACE_MTU_KEY) {
		xpath := fmt.Sprintf("%s/%s", base, openconfig.INTERFACE_MTU_KEY)
		i.append(srlib.ParseSrVal(config.Mtu, false, srlib.SR_UINT16_T, xpath))
	}

	if config.GetChange(openconfig.OC_DESCRIPTION_KEY) {
		xpath := fmt.Sprintf("%s/%s", base, openconfig.OC_DESCRIPTION_KEY)
		i.append(srlib.ParseSrVal(config.Desc, false, srlib.SR_STRING_T, xpath))
	}

	return nil
}

func (i *InterfaceVals) InterfaceEthernetConfig(name string, config *openconfig.InterfaceEthernetConfig) error {
	base := fmt.Sprintf("%s/%s", i.interfaceEthXPath(name), openconfig.OC_CONFIG_KEY)

	if config.GetChange(openconfig.INTERFACE_ETH_MACADDR_KEY) {
		xpath := fmt.Sprintf("%s/%s", base, openconfig.INTERFACE_ETH_MACADDR_KEY)
		i.append(srlib.ParseSrVal(config.MacAddr, false, srlib.SR_STRING_T, xpath))
	}

	return nil
}

func (i *InterfaceVals) Subinterface(name string, index uint32, subif *openconfig.Subinterface) error {
	return nil
}

func (i *InterfaceVals) SubinterfaceConfig(name string, index uint32, config *openconfig.SubinterfaceConfig) error {
	/*
		base := fmt.Sprintf("%s/%s", i.subinterfaceXPath(name, index), openconfig.OC_CONFIG_KEY)

		if config.GetChange(openconfig.OC_INDEX_KEY) {
			xpath := fmt.Sprintf("%s/%s", base, openconfig.OC_INDEX_KEY)
			i.append(srlib.ParseSrVal(config.Index, false, srlib.SR_UINT32_T, xpath))
		}

		if config.GetChange(openconfig.OC_ENABLED_KEY) {
			xpath := fmt.Sprintf("%s/%s", base, openconfig.OC_ENABLED_KEY)
			i.append(srlib.ParseSrVal(config.Enabled, false, srlib.SR_BOOL_T, xpath))
		}

		if config.GetChange(openconfig.OC_DESCRIPTION_KEY) {
			xpath := fmt.Sprintf("%s/%s", base, openconfig.OC_DESCRIPTION_KEY)
			i.append(srlib.ParseSrVal(config.Desc, false, srlib.SR_STRING_T, xpath))
		}
	*/
	return nil
}

func (i *InterfaceVals) SubinterfaceIPv4Config(name string, index uint32, config *openconfig.SubinterfaceIPv4Config) error {
	/*
		base := fmt.Sprintf("%s/%s", i.subinterfaceIPv4XPath(name, index), openconfig.OC_CONFIG_KEY)

		if config.GetChange(openconfig.SUBINTERFACE_MTU_KEY) {
			xpath := fmt.Sprintf("%s/%s", base, openconfig.SUBINTERFACE_MTU_KEY)
			i.append(srlib.ParseSrVal(config.Mtu, false, srlib.SR_UINT16_T, xpath))
		}
	*/
	return nil
}

func (i *InterfaceVals) SubinterfaceIPv4Address(name string, index uint32, ip string, addr *openconfig.IPAddress) error {
	/*
		base := i.subinterfaceIPv4AddrXPath(name, index, ip)
		xpath := fmt.Sprintf("%s/%s", base, openconfig.SUBINTERFACE_ADDR_IP_KEY)
		i.append(srlib.ParseSrVal(addr.IP, false, srlib.SR_STRING_T, xpath))
	*/
	return nil
}

func (i *InterfaceVals) SubinterfaceIPv4AddressConfig(name string, index uint32, ip string, config *openconfig.IPAddressConfig) error {
	/*
		base := fmt.Sprintf("%s/%s", i.subinterfaceIPv4AddrXPath(name, index, ip), openconfig.OC_CONFIG_KEY)

		if config.GetChange(openconfig.SUBINTERFACE_ADDR_IP_KEY) {
			xpath := fmt.Sprintf("%s/%s", base, openconfig.SUBINTERFACE_ADDR_IP_KEY)
			i.append(srlib.ParseSrVal(config.IP, false, srlib.SR_STRING_T, xpath))
		}

		if config.GetChange(openconfig.SUBINTERFACE_ADDR_PREFIXLEN_KEY) {
			xpath := fmt.Sprintf("%s/%s", base, openconfig.SUBINTERFACE_ADDR_PREFIXLEN_KEY)
			i.append(srlib.ParseSrVal(config.PrefixLen, false, srlib.SR_UINT8_T, xpath))
		}
	*/
	return nil
}
func (i *InterfaceVals) SubinterfaceIPv6Config(name string, index uint32, config *openconfig.SubinterfaceIPv6Config) error {
	return nil
}

func (i *InterfaceVals) SubinterfaceIPv6Address(name string, index uint32, ip string, addr *openconfig.IPAddress) error {
	return nil
}

func (i *InterfaceVals) SubinterfaceIPv6AddressConfig(name string, index uint32, ip string, config *openconfig.IPAddressConfig) error {
	return nil
}
