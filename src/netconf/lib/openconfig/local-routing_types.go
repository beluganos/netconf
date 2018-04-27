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
	"fmt"
	"net"
)

type LocalDefinedNexthop int

const (
	LOCAL_DEFINED_NEXT_HOP LocalDefinedNexthop = iota
	LOCAL_DEFINED_NEXT_HOP_DROP
	LOCAL_DEFINED_NEXT_HOP_LOCAL_LINK
)

var localDefinedNexthopNames = map[LocalDefinedNexthop]string{
	LOCAL_DEFINED_NEXT_HOP:            "LOCAL_DEFINED_NEXT_HOP",
	LOCAL_DEFINED_NEXT_HOP_DROP:       "DROP",
	LOCAL_DEFINED_NEXT_HOP_LOCAL_LINK: "LOCAL_LINK",
}

var localDefinedNexthopValues = map[string]LocalDefinedNexthop{
	"LOCAL_DEFINED_NEXT_HOP": LOCAL_DEFINED_NEXT_HOP,
	"DROP":       LOCAL_DEFINED_NEXT_HOP_DROP,
	"LOCAL_LINK": LOCAL_DEFINED_NEXT_HOP_LOCAL_LINK,
}

func (v LocalDefinedNexthop) String() string {
	if s, ok := localDefinedNexthopNames[v]; ok {
		return s
	}
	return fmt.Sprintf("LocalDefinedNexthop(%d)", v)
}

func ParseLocalDefinedNexthop(s string) (LocalDefinedNexthop, error) {
	if v, ok := localDefinedNexthopValues[s]; ok {
		return v, nil
	}
	return LOCAL_DEFINED_NEXT_HOP, fmt.Errorf("Invalid LocalDefinedNexthop. %s", s)
}

func ParseLocalDefinedNexthops(s string) (net.IP, LocalDefinedNexthop, error) {
	if nh, err := ParseLocalDefinedNexthop(s); err == nil {
		return nil, nh, nil
	}

	if ip := net.ParseIP(s); ip != nil {
		return ip, LOCAL_DEFINED_NEXT_HOP, nil
	}

	return nil, LOCAL_DEFINED_NEXT_HOP, fmt.Errorf("Invalid Nexthop. %s", s)
}
