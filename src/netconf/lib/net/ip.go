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

package ncnet

import (
	"net"
)

const (
	IPVER4      = 4
	IPVER6      = 6
	IPBITS4     = 32
	IPBITS6     = 128
	IPBITS_FULL = -1
)

func IPStringToVersion(s string) (int, error) {
	ip := net.ParseIP(s)
	if ip == nil {
		var err error
		if ip, _, err = net.ParseCIDR(s); err != nil {
			return 0, err
		}
	}
	return IPToVersion(ip), nil
}

func IPToVersion(ip net.IP) int {
	if ip.To4() == nil {
		return IPVER6
	}

	return IPVER4
}

func IPToBitlen(ip net.IP) int {
	if IPToVersion(ip) == 4 {
		return IPBITS4
	}

	return IPBITS6
}

func IPToIPNet(ip net.IP, plen int) *net.IPNet {
	if ip == nil {
		return nil
	}

	bitlen := IPToBitlen(ip)

	if plen == IPBITS_FULL {
		plen = bitlen
	}

	return &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(plen, bitlen),
	}
}
