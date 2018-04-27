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
	"bytes"
	"net"
	"testing"
)

func TestIFAddr_empty(t *testing.T) {
	ifa := IFAddr{}

	if v := ifa.IPNet(); v.String() != "<nil>" {
		t.Errorf("IFAddr.NetworkAddr unmatch. %v", v)
	}
}

func TestIFAddr_ipv4(t *testing.T) {
	ipv4 := net.IPv4(10, 1, 2, 3)
	mask := net.IPv4Mask(255, 255, 255, 0)
	ifv4 := NewIFAddr(ipv4, mask)

	if v := ifv4.String(); v != "10.1.2.3/24" {
		t.Errorf("IFAddr.String unmatch. %s", v)
	}

	if v := ifv4.IPNet(); v.String() != "10.1.2.0/24" {
		t.Errorf("IFAddr.IPNet unmatch. %s", v)
	}
}

func TestIFAddr_plen(t *testing.T) {
	ipv4 := net.IPv4(10, 1, 2, 3)
	plen := uint8(24)
	ifv4 := NewIFAddrWithPlen(ipv4, plen)

	if v := ifv4.String(); v != "10.1.2.3/24" {
		t.Errorf("IFAddr.String unmatch. %s", v)
	}

	if v := ifv4.IPNet(); v.String() != "10.1.2.0/24" {
		t.Errorf("IFAddr.IPNet unmatch. %s", v)
	}
}

func TestParseIFAddr(t *testing.T) {
	ifv4, err := ParseIFAddr("10.1.2.3/24")

	if err != nil {
		t.Errorf("ParseIFAddr error. %s", err)
	}

	if v := ifv4.String(); v != "10.1.2.3/24" {
		t.Errorf("ParseIFAddr unmatch. %s", v)
	}

	if v := ifv4.IPNet(); v.String() != "10.1.2.0/24" {
		t.Errorf("ParseIFAddr unmatch. %s", v)
	}
}

func TestIFName(t *testing.T) {

	if v := NewIFName("eth1", 0); v != "eth1" {
		t.Errorf("NewIFName unmatch. %s", v)
	}

	if v := NewIFName("eth1", 1); v != "eth1.1" {
		t.Errorf("NewIFName unmatch. %s", v)
	}
}

func TestParseIFName(t *testing.T) {
	var name string
	var index uint32
	var err error

	name, index, err = ParseIFName("eth1")
	if err != nil {
		t.Errorf("ParseIFName error. %s", err)
	}
	if name != "eth1" {
		t.Errorf("ParseIFName unmatch. %s", name)
	}
	if index != 0 {
		t.Errorf("ParseIFName unmatch. %d", index)
	}

	name, index, err = ParseIFName("eth1.1")
	if err != nil {
		t.Errorf("ParseIFName error. %s", err)
	}
	if name != "eth1" {
		t.Errorf("ParseIFName unmatch. %s", name)
	}
	if index != 1 {
		t.Errorf("ParseIFName unmatch. %d", index)
	}

	_, _, err = ParseIFName("")
	if err == nil {
		t.Errorf("ParseIFName must be error. %s", err)
	}
}

func TestRouteDistinguisher0(t *testing.T) {
	rd := RouteDistinguisher0{0, 1, 2, 3, 4, 5, 6, 7}

	if v := rd.Type(); v != RD_TYPE_0 {
		t.Errorf("Type unmatch. %s", v)
	}
	if v := rd.Bytes(); bytes.Compare(v, []byte{0, 0, 2, 3, 4, 5, 6, 7}) != 0 {
		t.Errorf("Bytes unmatch. %v", v)
	}
	if v := rd.AdminField(); bytes.Compare(v, []byte{2, 3}) != 0 {
		t.Errorf("AdminField unmatch. %v", v)
	}
	if v := rd.NumberField(); bytes.Compare(v, []byte{4, 5, 6, 7}) != 0 {
		t.Errorf("NumberField unmatch. %v", v)
	}
	if v := rd.String(); v != "515:67438087" {
		t.Errorf("String unmatch. %v", v)
	}
}

func TestParseRouteDistinguisher0(t *testing.T) {
	s := "291:1164413194" // 0x0123:0x4567890a
	rd, err := ParseRouteDistinguisher0(s)
	if err != nil {
		t.Errorf("ParseRouteDistinguisher0 error. %s", err)
	}
	if v := rd.AdminField(); bytes.Compare(v, []byte{0x01, 0x23}) != 0 {
		t.Errorf("AdminField unmatch. %v", v)
	}
	if v := rd.NumberField(); bytes.Compare(v, []byte{0x45, 0x67, 0x89, 0x0a}) != 0 {
		t.Errorf("NumberField unmatch. %v", v)
	}
	if v := rd.String(); v != s {
		t.Errorf("String unmatch. %s", v)
	}
}

func TestRouteDistinguisher1(t *testing.T) {
	rd := RouteDistinguisher1{0, 1, 2, 3, 4, 5, 6, 7}

	if v := rd.Type(); v != RD_TYPE_1 {
		t.Errorf("Type unmatch. %s", v)
	}
	if v := rd.Bytes(); bytes.Compare(v, []byte{0, 1, 2, 3, 4, 5, 6, 7}) != 0 {
		t.Errorf("Bytes unmatch. %v", v)
	}
	if v := rd.AdminField(); bytes.Compare(v, []byte{2, 3, 4, 5}) != 0 {
		t.Errorf("AdminField unmatch. %v", v)
	}
	if v := rd.NumberField(); bytes.Compare(v, []byte{6, 7}) != 0 {
		t.Errorf("NumberField unmatch. %v", v)
	}
	if v := rd.String(); v != "2.3.4.5:1543" {
		t.Errorf("String unmatch. %v", v)
	}
}

func TestParseRouteDistinguisher1(t *testing.T) {
	s := "1.2.3.4:5678"
	rd, err := ParseRouteDistinguisher1(s)
	if err != nil {
		t.Errorf("ParseRouteDistinguisher1 error. %s", err)
	}
	if v := rd.String(); v != s {
		t.Errorf("String unmatch. %s %#v", v, rd)
	}
}

func TestRouteDistinguisher2(t *testing.T) {
	rd := RouteDistinguisher2{0, 1, 2, 3, 4, 5, 6, 7}

	if v := rd.Type(); v != RD_TYPE_2 {
		t.Errorf("Type unmatch. %s", v)
	}
	if v := rd.Bytes(); bytes.Compare(v, []byte{0, 2, 2, 3, 4, 5, 6, 7}) != 0 {
		t.Errorf("Bytes unmatch. %v", v)
	}
	if v := rd.AdminField(); bytes.Compare(v, []byte{2, 3, 4, 5}) != 0 {
		t.Errorf("AdminField unmatch. %v", v)
	}
	if v := rd.NumberField(); bytes.Compare(v, []byte{6, 7}) != 0 {
		t.Errorf("NumberField unmatch. %v", v)
	}
	if v := rd.String(); v != "33752069:1543" {
		t.Errorf("String unmatch. %v", v)
	}
}

func TestParseRouteDistinguisher2(t *testing.T) {
	s := "19088743:35082" // 0x01234567:0x890a
	rd, err := ParseRouteDistinguisher2(s)
	if err != nil {
		t.Errorf("ParseRouteDistinguisher0 error. %s", err)
	}
	if v := rd.AdminField(); bytes.Compare(v, []byte{0x01, 0x23, 0x45, 0x67}) != 0 {
		t.Errorf("AdminField unmatch. %v", v)
	}
	if v := rd.NumberField(); bytes.Compare(v, []byte{0x89, 0x0a}) != 0 {
		t.Errorf("NumberField unmatch. %v", v)
	}
	if v := rd.String(); v != s {
		t.Errorf("String unmatch. %s", v)
	}
}
