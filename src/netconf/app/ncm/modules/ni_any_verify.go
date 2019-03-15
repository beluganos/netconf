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

package ncm

import (
	"fmt"
	"net"
	ncmdbm "netconf/app/ncm/dbm"
	"netconf/lib/openconfig"
)

func VerifyNIInterface(iface *openconfig.NetworkInstanceInterface) error {
	subif, _, err := ncmdbm.Subinterfaces().SelectById(iface.Id)
	if err != nil {
		return err
	}

	if subif.Index == 0 {
		if _, err := ncmdbm.Interfaces().Select(iface.Id); err != nil {
			return err
		}
	}
	return nil
}

func VerifyNIInterfaceRefConfig(config *openconfig.InterfaceRefConfig) error {
	keys := []string{openconfig.INTERFACE_KEY, openconfig.SUBINTERFACE_KEY}
	if chg := config.GetChanges(keys...); !chg {
		return fmt.Errorf("interface or subinterface not specified.")
	}

	return nil
}

func VerifyNIInterfaceConfig(id string, config *openconfig.NetworkInstanceInterfaceConfig) error {
	keys := []string{openconfig.INTERFACE_KEY, openconfig.SUBINTERFACE_KEY}

	if config.GetChanges(keys...) {
		if id != config.IFName() {
			return fmt.Errorf("Invalid Interface Id. %s", config)
		}

	} else if config.Compare(openconfig.SUBINTERFACE_KEY) {
		// pass
	} else if config.OneOfChange(keys...) {
		return fmt.Errorf("Interface or Subinterface not exist. %s", config)
	}

	return nil
}

func verifyIPAndPrefixLen(ip net.IP, plen uint8) error {
	if ip.To4() != nil {
		// ip is ipv4 address
		if plen > 32 {
			return fmt.Errorf("invalid prefix-length. %s/%d", ip, plen)
		}
	}

	return nil
}

func VerifyNILoopbackAddrConfig(config *openconfig.NetworkInstanceLoopbackAddrConfig) error {
	keys := []string{openconfig.NETWORKINSTANCE_LO_IP_KEY, openconfig.NETWORKINSTANCE_LO_PLEN_KEY}

	if chg := config.GetChanges(keys...); !chg {
		return fmt.Errorf("ip or prefix not specified. %s/%d", config.Ip, config.PrefixLen)
	}

	if err := verifyIPAndPrefixLen(config.Ip, config.PrefixLen); err != nil {
		return err
	}

	return nil
}

func VerifyNIStaticRouteConfig(config *openconfig.StaticRouteConfig) error {
	keys := []string{openconfig.STATICROUTE_IP_KEY, openconfig.STATICROUTE_PREFIXLEN_KEY}

	if chg := config.GetChanges(keys...); !chg {
		return fmt.Errorf("ip or prefix not specified. %s/%d", config.Ip, config.PrefixLen)
	}

	if err := verifyIPAndPrefixLen(config.Ip, config.PrefixLen); err != nil {
		return err
	}

	return nil
}
