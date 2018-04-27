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

package vtylib

import (
	"fmt"
)

type DaemonState int

const (
	DAEMON_STATE DaemonState = iota
	DAEMON_STATE_YES
	DAEMON_STATE_NO
)

var daemonStateNames = map[DaemonState]string{
	DAEMON_STATE:     "DAEMON_STATE",
	DAEMON_STATE_YES: "yes",
	DAEMON_STATE_NO:  "no",
}

var daemonStateValues = map[string]DaemonState{
	"DAEMON_STATE": DAEMON_STATE,
	"yes":          DAEMON_STATE_YES,
	"no":           DAEMON_STATE_NO,
}

func ParseDaemonState(s string) (DaemonState, error) {
	if v, ok := daemonStateValues[s]; ok {
		return v, nil
	}
	return DAEMON_STATE, fmt.Errorf("Invalid DaemonState. %s", s)
}

func (v DaemonState) String() string {
	if s, ok := daemonStateNames[v]; ok {
		return s
	}
	return fmt.Sprintf("DaemonState(%d)", v)
}

type DaemonType int

const (
	DAEMON_TYPE DaemonType = iota
	DAEMON_ZEBRA
	DAEMON_OSPF
	DAEMON_OSPF6
	DAEMON_BGP
	DAEMON_EIGRP
	DAEMON_ISIS
	DAEMON_LDP
	DAEMON_NHRP
	DAEMON_PIM
	DAEMON_RIP
	DAEMON_RIPNG
)

var daemonTypeName = map[DaemonType]string{
	DAEMON_ZEBRA: "zebra",
	DAEMON_OSPF:  "ospfd",
	DAEMON_OSPF6: "ospf6d",
	DAEMON_BGP:   "bgpd",
	DAEMON_EIGRP: "eigrp",
	DAEMON_ISIS:  "isisd",
	DAEMON_LDP:   "ldpd",
	DAEMON_NHRP:  "nhrpd",
	DAEMON_PIM:   "pimd",
	DAEMON_RIP:   "ripd",
	DAEMON_RIPNG: "ripngd",
}

var daemonTypeVal = map[string]DaemonType{
	"zebra":  DAEMON_ZEBRA,
	"ospfd":  DAEMON_OSPF,
	"ospf6d": DAEMON_OSPF6,
	"bgpd":   DAEMON_BGP,
	"eigrp":  DAEMON_EIGRP,
	"isisd":  DAEMON_ISIS,
	"ldpd":   DAEMON_LDP,
	"nhrpd":  DAEMON_NHRP,
	"pimd":   DAEMON_PIM,
	"ripd":   DAEMON_RIP,
	"ripngd": DAEMON_RIPNG,
}

func ParseDaemonType(s string) (DaemonType, error) {
	if v, ok := daemonTypeVal[s]; ok {
		return v, nil
	}
	return DAEMON_TYPE, fmt.Errorf("Invalid DaemonType. %s", s)
}

func (d DaemonType) String() string {
	if s, ok := daemonTypeName[d]; ok {
		return s
	}
	return fmt.Sprintf("DaemonType(%d)", d)
}
