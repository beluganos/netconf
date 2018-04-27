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
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"net"
	"strconv"
	"strings"
)

//
// Route Distinguisher Type
//
type RouteDistinguisherType uint16

const (
	RD_TYPE_0    RouteDistinguisherType = 0
	RD_TYPE_1    RouteDistinguisherType = 1
	RD_TYPE_2    RouteDistinguisherType = 2
	RD_TYPE_NONE RouteDistinguisherType = 65535
)

var routeDistinguisherType_names = map[RouteDistinguisherType]string{
	RD_TYPE_0:    "type0",
	RD_TYPE_1:    "type1",
	RD_TYPE_2:    "type2",
	RD_TYPE_NONE: "none",
}

var routeDistinguisherType_values = map[string]RouteDistinguisherType{
	"type0": RD_TYPE_0,
	"type1": RD_TYPE_1,
	"type2": RD_TYPE_2,
}

func (v RouteDistinguisherType) String() string {
	if s, ok := routeDistinguisherType_names[v]; ok {
		return s
	}
	return fmt.Sprintf("RouteDistinguisherType(%d)", v)
}

func ParseRouteDistinguisherType(s string) (RouteDistinguisherType, error) {
	if v, ok := routeDistinguisherType_values[s]; ok {
		return v, nil
	}
	return RD_TYPE_NONE, fmt.Errorf("Invalid RouteDistinguisherType. %s", s)
}

//
// Route Distinguisher
//
type RouteDistinguisher interface {
	xml.Marshaler
	Type() RouteDistinguisherType
	AdminField() []byte
	NumberField() []byte
	Bytes() []byte
	String() string
}

func ParseRouteDistinguisher(s string) (RouteDistinguisher, error) {
	if len(s) == 0 {
		return RouteDistinguisherNone{}, nil
	}

	if rd, err := ParseRouteDistinguisher0(s); err == nil {
		return rd, nil
	}

	if rd, err := ParseRouteDistinguisher1(s); err == nil {
		return rd, nil
	}

	if rd, err := ParseRouteDistinguisher2(s); err == nil {
		return rd, nil
	}

	return RouteDistinguisherNone{}, fmt.Errorf("Invalid RD/RT. %s", s)
}

//
// RD (Type 0)
//
type RouteDistinguisher0 [8]byte

func (r RouteDistinguisher0) Type() RouteDistinguisherType {
	return RD_TYPE_0
}

func (r RouteDistinguisher0) AdminField() []byte {
	return r[2:4]
}

func (r RouteDistinguisher0) NumberField() []byte {
	return r[4:]
}

func (r RouteDistinguisher0) Bytes() []byte {
	r[1] = byte(RD_TYPE_0)
	return r[:]
}

func (r RouteDistinguisher0) String() string {
	admin := binary.BigEndian.Uint16(r.AdminField())
	number := binary.BigEndian.Uint32(r.NumberField())
	return fmt.Sprintf("%d:%d", admin, number)
}

func (r RouteDistinguisher0) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(r.String(), start)
}

func ParseRouteDistinguisher0(s string) (RouteDistinguisher0, error) {
	rd := RouteDistinguisher0{}

	fields := strings.Split(s, ":")
	if len(fields) != 2 {
		return rd, fmt.Errorf("Invalid RD. %s", s)
	}
	if len(fields[0]) == 0 || len(fields[1]) == 0 {
		return rd, fmt.Errorf("Invalid RD. %s", s)
	}

	admin, err := strconv.ParseUint(fields[0], 0, 16)
	if err != nil {
		return rd, fmt.Errorf("Invalid RD. %s %s", s, err)
	}

	number, err := strconv.ParseUint(fields[1], 0, 32)
	if err != nil {
		return rd, fmt.Errorf("Invalid RD. %s %s", s, err)
	}

	binary.BigEndian.PutUint16(rd[2:4], uint16(admin))
	binary.BigEndian.PutUint32(rd[4:], uint32(number))
	return rd, nil
}

//
// RD (Type 1)
//
type RouteDistinguisher1 [8]byte

func (r RouteDistinguisher1) Type() RouteDistinguisherType {
	return RD_TYPE_1
}

func (r RouteDistinguisher1) AdminField() []byte {
	return r[2:6]
}

func (r RouteDistinguisher1) NumberField() []byte {
	return r[6:]
}

func (r RouteDistinguisher1) Bytes() []byte {
	r[1] = byte(RD_TYPE_1)
	return r[:]
}

func (r RouteDistinguisher1) String() string {
	admin := net.IP(r.AdminField())
	number := binary.BigEndian.Uint16(r.NumberField())
	return fmt.Sprintf("%s:%d", admin, number)
}

func (r RouteDistinguisher1) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(r.String(), start)
}

func ParseRouteDistinguisher1(s string) (RouteDistinguisher1, error) {
	rd := RouteDistinguisher1{}

	fields := strings.Split(s, ":")
	if len(fields) != 2 {
		return rd, fmt.Errorf("Invalid RD. %s", s)
	}
	if len(fields[0]) == 0 || len(fields[1]) == 0 {
		return rd, fmt.Errorf("Invalid RD. %s", s)
	}

	admin := net.ParseIP(fields[0])
	if admin == nil || admin.To4() == nil {
		return rd, fmt.Errorf("Invalid RD. %s", s)
	}

	number, err := strconv.ParseUint(fields[1], 0, 16)
	if err != nil {
		return rd, fmt.Errorf("Invalid RD. %s %s", s, err)
	}

	copy(rd[2:], admin.To4())
	binary.BigEndian.PutUint16(rd[6:], uint16(number))
	return rd, nil
}

//
// RD (Type2)
//
type RouteDistinguisher2 [8]byte

func (r RouteDistinguisher2) Type() RouteDistinguisherType {
	return RD_TYPE_2
}

func (r RouteDistinguisher2) AdminField() []byte {
	return r[2:6]
}

func (r RouteDistinguisher2) NumberField() []byte {
	return r[6:]
}

func (r RouteDistinguisher2) Bytes() []byte {
	r[1] = byte(RD_TYPE_2)
	return r[:]
}

func (r RouteDistinguisher2) String() string {
	admin := binary.BigEndian.Uint32(r.AdminField())
	number := binary.BigEndian.Uint16(r.NumberField())
	return fmt.Sprintf("%d:%d", admin, number)
}

func (r RouteDistinguisher2) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(r.String(), start)
}

func ParseRouteDistinguisher2(s string) (RouteDistinguisher2, error) {
	rd := RouteDistinguisher2{}

	fields := strings.Split(s, ":")
	if len(fields) != 2 {
		return rd, fmt.Errorf("Invalid RD. %s", s)
	}
	if len(fields[0]) == 0 || len(fields[1]) == 0 {
		return rd, fmt.Errorf("Invalid RD. %s", s)
	}

	admin, err := strconv.ParseUint(fields[0], 0, 32)
	if err != nil {
		return rd, fmt.Errorf("Invalid RD. %s %s", s, err)
	}

	number, err := strconv.ParseUint(fields[1], 0, 16)
	if err != nil {
		return rd, fmt.Errorf("Invalid RD. %s %s", s, err)
	}

	binary.BigEndian.PutUint32(rd[2:6], uint32(admin))
	binary.BigEndian.PutUint16(rd[6:], uint16(number))
	return rd, nil
}

//
// RD (NONE)
//
type RouteDistinguisherNone struct {
}

func (r RouteDistinguisherNone) Type() RouteDistinguisherType {
	return RD_TYPE_NONE
}

func (r RouteDistinguisherNone) AdminField() []byte {
	return []byte{}
}

func (r RouteDistinguisherNone) NumberField() []byte {
	return []byte{}
}

func (r RouteDistinguisherNone) Bytes() []byte {
	return []byte{}
}

func (r RouteDistinguisherNone) String() string {
	return ""
}

func (r RouteDistinguisherNone) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return nil
}
