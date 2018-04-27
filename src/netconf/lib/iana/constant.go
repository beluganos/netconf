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

package ncianalib

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type IANAifType string

const (
	IANAifType_XMLNS  = "urn:ietf:params:xml:ns:yang:iana-if-type"
	IANAifType_MODULE = "iana-if-type"
)

const (
	IANAifType_other          = IANAifType("other")
	IANAifType_regular1822    = IANAifType("regular1822")
	IANAifType_hdh1822        = IANAifType("hdh1822")
	IANAifType_ddnX25         = IANAifType("ddnX25")
	IANAifType_rfc877x25      = IANAifType("rfc877x25")
	IANAifType_ethernetCsmacd = IANAifType("ethernetCsmacd")
)

var IANAifType_Values = map[IANAifType]uint32{
	IANAifType_other:          1,
	IANAifType_regular1822:    2,
	IANAifType_hdh1822:        3,
	IANAifType_ddnX25:         4,
	IANAifType_rfc877x25:      5,
	IANAifType_ethernetCsmacd: 6,
}

var IANAifType_Names = map[uint32]IANAifType{
	1: IANAifType_other,
	2: IANAifType_regular1822,
	3: IANAifType_hdh1822,
	4: IANAifType_ddnX25,
	5: IANAifType_rfc877x25,
	6: IANAifType_ethernetCsmacd,
}

func ParseIANAifType(s string) (IANAifType, error) {
	ss := strings.SplitN(s, ":", 2)
	ianaType := IANAifType(ss[len(ss)-1])
	if _, ok := IANAifType_Values[ianaType]; ok {
		return ianaType, nil
	}

	return IANAifType_other, fmt.Errorf("Invalid IANAifType. %s", s)
}

func (i IANAifType) String() string {
	return fmt.Sprintf("%s:%s", IANAifType_MODULE, string(i))
}

func (i IANAifType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	xmlNsAttr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", IANAifType_MODULE)},
		Value: IANAifType_XMLNS,
	}
	start.Attr = append(start.Attr, xmlNsAttr)
	return e.EncodeElement(i.String(), start)
}
