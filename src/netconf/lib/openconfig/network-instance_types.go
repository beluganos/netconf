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
	ncxml "netconf/lib/xml"
)

const (
	NETWORK_INSTANCE_TYPES_XMLNS  = "http://openconfig.net/yang/network-instance-types"
	NETWORK_INSTANCE_TYPES_MODULE = "oc-ni-types"
)

//
// NETWORK_INSTANCE_TYPE
//
type NetworkInstanceType int

const (
	NETWORK_INSTANCE_TYPE NetworkInstanceType = iota
	NETWORK_INSTANCE_DEFAULT
	NETWORK_INSTANCE_L3VRF
	NETWORK_INSTANCE_L2VSI
	NETWORK_INSTANCE_L2P2P
	NETWORK_INSTANCE_L2L3
)

var networkInstanceType_name = map[NetworkInstanceType]string{
	NETWORK_INSTANCE_TYPE:    "NETWORK_INSTANCE_TYPE",
	NETWORK_INSTANCE_DEFAULT: "DEFAULT_INSTANCE",
	NETWORK_INSTANCE_L3VRF:   "L3VRF",
	NETWORK_INSTANCE_L2VSI:   "L2VSI",
	NETWORK_INSTANCE_L2P2P:   "L2P2P",
	NETWORK_INSTANCE_L2L3:    "L2L3",
}

var networkInstanceType_values = map[string]NetworkInstanceType{
	"NETWORK_INSTANCE_TYPE": NETWORK_INSTANCE_TYPE,
	"DEFAULT_INSTANCE":      NETWORK_INSTANCE_DEFAULT,
	"L3VRF":                 NETWORK_INSTANCE_L3VRF,
	"L2VSI":                 NETWORK_INSTANCE_L2VSI,
	"L2P2P":                 NETWORK_INSTANCE_L2P2P,
	"L2L3":                  NETWORK_INSTANCE_L2L3,
}

func (t NetworkInstanceType) String() string {
	if s, ok := networkInstanceType_name[t]; ok {
		return s
	}
	return fmt.Sprintf("NetworkInstanceType(%d)", t)
}

func ParseNetworkInstanceType(s string) (NetworkInstanceType, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := networkInstanceType_values[ss]; ok {
		return v, nil
	}
	return NETWORK_INSTANCE_TYPE, fmt.Errorf("Invalid NetworkInstanceType. %s", s)
}

func (t NetworkInstanceType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	attr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", NETWORK_INSTANCE_TYPES_MODULE)},
		Value: NETWORK_INSTANCE_TYPES_XMLNS,
	}
	text := fmt.Sprintf("%s:%s", NETWORK_INSTANCE_TYPES_MODULE, t)
	start.Attr = append(start.Attr, attr)
	return e.EncodeElement(text, start)
}

//
// ENDPOINT_TYPE
//
type EndpointType int

const (
	ENDPOINT_TYPE EndpointType = iota
	ENDPOINT_LOCAL
	ENDPOINT_REMOTE
)

var endpointType_names = map[EndpointType]string{
	ENDPOINT_TYPE:   "ENDPOINT_TYPE",
	ENDPOINT_LOCAL:  "LOCAL",
	ENDPOINT_REMOTE: "REMOTE",
}

var endpointType_values = map[string]EndpointType{
	"ENDPOINT_TYPE": ENDPOINT_TYPE,
	"LOCAL":         ENDPOINT_LOCAL,
	"REMOTE":        ENDPOINT_REMOTE,
}

func (v EndpointType) String() string {
	if s, ok := endpointType_names[v]; ok {
		return s
	}
	return fmt.Sprintf("EndpointType(%d)", v)
}

func ParseEndpointType(s string) (EndpointType, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := endpointType_values[ss]; ok {
		return v, nil
	}
	return ENDPOINT_TYPE, fmt.Errorf("Invalid EndpointType. %s", s)
}

func (t EndpointType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	attr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", NETWORK_INSTANCE_TYPES_MODULE)},
		Value: NETWORK_INSTANCE_TYPES_XMLNS,
	}
	text := fmt.Sprintf("%s:%s", NETWORK_INSTANCE_TYPES_MODULE, t)
	start.Attr = append(start.Attr, attr)
	return e.EncodeElement(text, start)
}

//
// LABEL_ALLOCATION_MODE
//
type LabelAllocationMode int

const (
	LABEL_ALLOCATION_MODE LabelAllocationMode = iota
	LABEL_ALLOCATION_PER_PREFIX
	LABEL_ALLOCATION_PER_NEXTHOP
	LABEL_ALLOCATION_INSTANCE_LABEL
)

var labelAllocationMode_names = map[LabelAllocationMode]string{
	LABEL_ALLOCATION_MODE:           "LABEL_ALLOCATION_MODE",
	LABEL_ALLOCATION_PER_PREFIX:     "PER_PREFIX",
	LABEL_ALLOCATION_PER_NEXTHOP:    "PER_NEXTHOP",
	LABEL_ALLOCATION_INSTANCE_LABEL: "INSTANCE_LABEL",
}

var labelAllocationMode_values = map[string]LabelAllocationMode{
	"LABEL_ALLOCATION_MODE": LABEL_ALLOCATION_MODE,
	"PER_PREFIX":            LABEL_ALLOCATION_PER_PREFIX,
	"PER_NEXTHOP":           LABEL_ALLOCATION_PER_NEXTHOP,
	"INSTANCE_LABEL":        LABEL_ALLOCATION_INSTANCE_LABEL,
}

func (v LabelAllocationMode) String() string {
	if s, ok := labelAllocationMode_names[v]; ok {
		return s
	}
	return fmt.Sprintf("LabelAllocationMode(%d)", v)
}

func ParseLabelAllocationMode(s string) (LabelAllocationMode, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := labelAllocationMode_values[ss]; ok {
		return v, nil
	}
	return LABEL_ALLOCATION_MODE, fmt.Errorf("Invalid LabelAllocationMode. %s", s)
}

func (t LabelAllocationMode) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	attr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", NETWORK_INSTANCE_TYPES_MODULE)},
		Value: NETWORK_INSTANCE_TYPES_XMLNS,
	}
	text := fmt.Sprintf("%s:%s", NETWORK_INSTANCE_TYPES_MODULE, t)
	start.Attr = append(start.Attr, attr)
	return e.EncodeElement(text, start)
}

//
// ENCAPSULATION
//
type Encapsulation int

const (
	ENCAPSULATION Encapsulation = iota
	ENCAPSULATION_MPLS
	ENCAPSULATION_VXLAN
)

var encapsulation_names = map[Encapsulation]string{
	ENCAPSULATION:       "ENCAPSULATION",
	ENCAPSULATION_MPLS:  "ENCAPSULATION_MPLS",
	ENCAPSULATION_VXLAN: "ENCAPSULATION_VXLAN",
}

var encapsulation_values = map[string]Encapsulation{
	"MPLS":  ENCAPSULATION_MPLS,
	"VXLAN": ENCAPSULATION_VXLAN,
}

func (v Encapsulation) String() string {
	if s, ok := encapsulation_names[v]; ok {
		return s
	}
	return fmt.Sprintf("Encapsulation(%d)", v)
}

func ParseEncapsulation(s string) (Encapsulation, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := encapsulation_values[ss]; ok {
		return v, nil
	}
	return ENCAPSULATION, fmt.Errorf("Invalid Encapsulation. %s", s)
}

func (t Encapsulation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	attr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", NETWORK_INSTANCE_TYPES_MODULE)},
		Value: NETWORK_INSTANCE_TYPES_XMLNS,
	}
	text := fmt.Sprintf("%s:%s", NETWORK_INSTANCE_TYPES_MODULE, t)
	start.Attr = append(start.Attr, attr)
	return e.EncodeElement(text, start)
}

//
// SIGNALLING_PROTOCOL
//
type SignallingProtocol int

const (
	SIGNALLING_PROTOCOL SignallingProtocol = iota
	SIGNALLING_PROTOCOL_LDP
	SIGNALLING_PROTOCOL_BGP_VPLS
	SIGNALLING_PROTOCOL_BGP_EVPN
)

var signallingProtocol_names = map[SignallingProtocol]string{
	SIGNALLING_PROTOCOL:          "SIGNALLING_PROTOCOL",
	SIGNALLING_PROTOCOL_LDP:      "LDP",
	SIGNALLING_PROTOCOL_BGP_VPLS: "BGP_VPLS",
	SIGNALLING_PROTOCOL_BGP_EVPN: "BGP_EVPN",
}

var signallingProtocol_values = map[string]SignallingProtocol{
	"SIGNALLING_PROTOCOL": SIGNALLING_PROTOCOL,
	"LDP":                 SIGNALLING_PROTOCOL_LDP,
	"BGP_VPLS":            SIGNALLING_PROTOCOL_BGP_VPLS,
	"BGP_EVPN":            SIGNALLING_PROTOCOL_BGP_EVPN,
}

func (v SignallingProtocol) String() string {
	if s, ok := signallingProtocol_names[v]; ok {
		return s
	}
	return fmt.Sprintf("SignallingProtocol(%d)", v)
}

func ParseSignallingProtocol(s string) (SignallingProtocol, error) {
	_, ss := ncxml.ParseXPathName(s)
	if v, ok := signallingProtocol_values[ss]; ok {
		return v, nil
	}
	return SIGNALLING_PROTOCOL, fmt.Errorf("Invalid SignallingProtocol. %s", s)
}

func (t SignallingProtocol) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	attr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", NETWORK_INSTANCE_TYPES_MODULE)},
		Value: NETWORK_INSTANCE_TYPES_XMLNS,
	}
	text := fmt.Sprintf("%s:%s", NETWORK_INSTANCE_TYPES_MODULE, t)
	start.Attr = append(start.Attr, attr)
	return e.EncodeElement(text, start)
}
