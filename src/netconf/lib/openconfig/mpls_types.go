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
	"netconf/lib/xml"
)

type MplsNullLabelType int

const (
	MPLS_TYPES_XMLNS = "http://openconfig.net/yang/mpls-types"
	MPLS_TYPE_MODULE = "oc-mpls-types"
)

const (
	MPLS_NULL_LABEL_TYPE MplsNullLabelType = iota
	MPLS_NULL_LABEL_EXPLICIT
	MPLS_NULL_LABEL_IMPLICIT
)

var mplsNullLabelTypeNames = map[MplsNullLabelType]string{
	MPLS_NULL_LABEL_TYPE:     "NULL_LABEL_TYPE",
	MPLS_NULL_LABEL_EXPLICIT: "EXPLICIT",
	MPLS_NULL_LABEL_IMPLICIT: "IMPLICIT",
}

var mplsNullLabelTypeValues = map[string]MplsNullLabelType{
	"NULL_LABEL_TYPE": MPLS_NULL_LABEL_TYPE,
	"EXPLICIT":        MPLS_NULL_LABEL_EXPLICIT,
	"IMPLICIT":        MPLS_NULL_LABEL_IMPLICIT,
}

func (v MplsNullLabelType) String() string {
	if s, ok := mplsNullLabelTypeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("MplsNullLabelType(%d)", v)
}

func ParseMplsNullLabelType(s string) (MplsNullLabelType, error) {
	_, name := ncxml.ParseXPathName(s)
	if v, ok := mplsNullLabelTypeValues[name]; ok {
		return v, nil
	}
	return MPLS_NULL_LABEL_TYPE, fmt.Errorf("Invalid MplsNullLabelType. %s", s)
}

func (v MplsNullLabelType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	attr := xml.Attr{
		Name:  xml.Name{Local: fmt.Sprintf("xmlns:%s", MPLS_TYPE_MODULE)},
		Value: MPLS_TYPES_XMLNS,
	}
	text := fmt.Sprintf("%s:%s", MPLS_TYPE_MODULE, v)
	start.Attr = append(start.Attr, attr)
	return e.EncodeElement(text, start)
}
