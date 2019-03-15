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

//
// Interface address
//
type IFAddr struct {
	IP   net.IP
	Mask net.IPMask
}

func NewIFAddr(ip net.IP, mask net.IPMask) *IFAddr {
	if ip == nil {
		ip = net.IP{}
	}

	if mask == nil {
		mask = net.IPMask{}
	}

	return &IFAddr{
		IP:   ip,
		Mask: mask,
	}
}

func NewIFAddrWithPlen(ip net.IP, plen uint8) *IFAddr {
	ifa := NewIFAddr(ip, nil)
	ifa.SetPLen(plen)
	return ifa
}

func ParseIFAddr(s string) (*IFAddr, error) {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}

	return NewIFAddr(ip, ipnet.Mask), nil
}

func (i *IFAddr) String() string {
	ipn := net.IPNet{
		IP:   i.IP,
		Mask: i.Mask,
	}
	return ipn.String()
}

func (i *IFAddr) SetPLen(plen uint8) {
	i.Mask = net.CIDRMask(int(plen), IPToBitlen(i.IP))
}

func (i *IFAddr) PLen() uint8 {
	ones, _ := i.Mask.Size()
	return uint8(ones)
}

func (i *IFAddr) IPNet() *net.IPNet {
	return &net.IPNet{
		IP:   i.IP.Mask(i.Mask),
		Mask: i.Mask,
	}
}

func (i *IFAddr) IPVer() int {
	return IPToVersion(i.IP)
}
