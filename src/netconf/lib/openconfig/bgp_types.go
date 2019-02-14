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
	"encoding/xml"
	"fmt"
	"net"
	ncxml "netconf/lib/xml"
)

type BgpNexthopType string

const (
	BGP_NEXTHOP_TYPE BgpNexthopType = "BGP_NEXTHOP_TYPE"
	BGP_NEXTHOP_SELF BgpNexthopType = "SELF"
)

func ParseBgpNexthopType(s string) (BgpNexthopType, error) {
	_, ss := ncxml.ParseXPathName(s)
	b := BgpNexthopType(ss)
	if _, _, err := b.Values(); err != nil {
		return BGP_NEXTHOP_TYPE, err
	}
	return b, nil
}

func (b BgpNexthopType) Values() (net.IP, BgpNexthopType, error) {
	if b == BGP_NEXTHOP_SELF {
		return nil, b, nil
	}

	if ip := net.ParseIP(string(b)); ip != nil {
		return ip, BGP_NEXTHOP_TYPE, nil
	}

	return nil, BGP_NEXTHOP_TYPE, fmt.Errorf("Invalid Bgp Nexthop. %s", b)
}

func (b BgpNexthopType) String() string {
	return string(b)
}

type BgpAfiSafiType uint32

func NewBgpAfiSafiType(afi uint16, safi uint16) BgpAfiSafiType {
	return BgpAfiSafiType((uint32(afi)<<16 | uint32(safi)))
}

func (b BgpAfiSafiType) AFI() uint16 {
	return uint16((b & 0xffff0000) >> 16)
}

func (b BgpAfiSafiType) SAFI() uint16 {
	return uint16(b & 0x0000ffff)
}

const (
	BGP_AFI_IPV4       = 1
	BGP_AFI_IPV6       = 2
	BGP_AFI_L2VPN      = 25
	BGP_AFI_IPV4_SAFI  = uint32(BGP_AFI_IPV4 << 16)
	BGP_AFI_IPV6_SAFI  = uint32(BGP_AFI_IPV6 << 16)
	BGP_AFI_L2VPN_SAFI = uint32(BGP_AFI_L2VPN << 16)
)

const (
	BGP_AFI_SAFI_TYPE                 = BgpAfiSafiType(0)
	BGP_AFI_SAFI_IPV4_UNICAST         = BgpAfiSafiType(BGP_AFI_IPV4_SAFI | 1)
	BGP_AFI_SAFI_IPV6_UNICAST         = BgpAfiSafiType(BGP_AFI_IPV6_SAFI | 1)
	BGP_AFI_SAFI_IPV4_LABELED_UNICAST = BgpAfiSafiType(BGP_AFI_IPV4_SAFI | 4)
	BGP_AFI_SAFI_IPV6_LABELED_UNICAST = BgpAfiSafiType(BGP_AFI_IPV6_SAFI | 4)
	BGP_AFI_SAFI_L3VPN_IPV4_UNICAST   = BgpAfiSafiType(BGP_AFI_IPV4_SAFI | 128)
	BGP_AFI_SAFI_L3VPN_IPV6_UNICAST   = BgpAfiSafiType(BGP_AFI_IPV6_SAFI | 128)
	BGP_AFI_SAFI_L3VPN_IPV4_MULTICAST = BgpAfiSafiType(BGP_AFI_IPV4_SAFI | 129)
	BGP_AFI_SAFI_L3VPN_IPV6_MULTICAST = BgpAfiSafiType(BGP_AFI_IPV6_SAFI | 129)
	BGP_AFI_SAFI_L2VPN_VPLS           = BgpAfiSafiType(BGP_AFI_L2VPN_SAFI | 65)
	BGP_AFI_SAFI_L2VPN_EVPN           = BgpAfiSafiType(BGP_AFI_L2VPN_SAFI | 70)
)

var bgpAfiSafiTypeNames = map[BgpAfiSafiType]string{
	BGP_AFI_SAFI_TYPE:                 "BGP_AFI_SAFI_TYPE",
	BGP_AFI_SAFI_IPV4_UNICAST:         "IPV4_UNICAST",
	BGP_AFI_SAFI_IPV6_UNICAST:         "IPV6_UNICAST",
	BGP_AFI_SAFI_IPV4_LABELED_UNICAST: "IPV4_LABELED_UNICAST",
	BGP_AFI_SAFI_IPV6_LABELED_UNICAST: "IPV6_LABELED_UNICAST",
	BGP_AFI_SAFI_L3VPN_IPV4_UNICAST:   "L3VPN_IPV4_UNICAST",
	BGP_AFI_SAFI_L3VPN_IPV6_UNICAST:   "L3VPN_IPV6_UNICAST",
	BGP_AFI_SAFI_L3VPN_IPV4_MULTICAST: "L3VPN_IPV4_MULTICAST",
	BGP_AFI_SAFI_L3VPN_IPV6_MULTICAST: "L3VPN_IPV6_MULTICAST",
	BGP_AFI_SAFI_L2VPN_VPLS:           "L2VPN_VPLS",
	BGP_AFI_SAFI_L2VPN_EVPN:           "L2VPN_EVPN",
}

var bgpAfiSafiTypeValues = map[string]BgpAfiSafiType{
	"BGP_AFI_SAFI_TYPE":    BGP_AFI_SAFI_TYPE,
	"IPV4_UNICAST":         BGP_AFI_SAFI_IPV4_UNICAST,
	"IPV6_UNICAST":         BGP_AFI_SAFI_IPV6_UNICAST,
	"IPV4_LABELED_UNICAST": BGP_AFI_SAFI_IPV4_LABELED_UNICAST,
	"IPV6_LABELED_UNICAST": BGP_AFI_SAFI_IPV6_LABELED_UNICAST,
	"L3VPN_IPV4_UNICAST":   BGP_AFI_SAFI_L3VPN_IPV4_UNICAST,
	"L3VPN_IPV6_UNICAST":   BGP_AFI_SAFI_L3VPN_IPV6_UNICAST,
	"L3VPN_IPV4_MULTICAST": BGP_AFI_SAFI_L3VPN_IPV4_MULTICAST,
	"L3VPN_IPV6_MULTICAST": BGP_AFI_SAFI_L3VPN_IPV6_MULTICAST,
	"L2VPN_VPLS":           BGP_AFI_SAFI_L2VPN_VPLS,
	"L2VPN_EVPN":           BGP_AFI_SAFI_L2VPN_EVPN,
}

func (b BgpAfiSafiType) String() string {
	if s, ok := bgpAfiSafiTypeNames[b]; ok {
		return s
	}
	return fmt.Sprintf("BgpAfiSafiType(%d)", b)
}

func ParseBgpAfiSafiType(s string) (BgpAfiSafiType, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := bgpAfiSafiTypeValues[ss]; ok {
		return v, nil
	}
	return BGP_AFI_SAFI_TYPE, fmt.Errorf("Invalid BgpAfiSafiType. %s", s)
}

func (b BgpAfiSafiType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	attr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", BGP_TYPES_MODULE)},
		Value: BGP_TYPES_XMLNS,
	}
	text := fmt.Sprintf("%s:%s", BGP_TYPES_MODULE, b)
	start.Attr = append(start.Attr, attr)
	return e.EncodeElement(text, start)
}
