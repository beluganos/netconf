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

package cfgsyslib

import (
	"net"
	"strings"
)

type NetowrkAddrs []*net.IPNet

func (n *NetowrkAddrs) Set(value string) error {
	ip, nw, err := net.ParseCIDR(value)
	if err != nil {
		return err
	}
	a := &net.IPNet{
		IP:   ip,
		Mask: nw.Mask,
	}
	*n = append(*n, a)
	return nil
}

func (n NetowrkAddrs) Type() string {
	return "cfgsyslib.NetowrkAddrs"
}

func (n NetowrkAddrs) String() string {
	return strings.Join(n.Strings(), "|")
}

func (n NetowrkAddrs) Strings() []string {
	ss := make([]string, len(n))
	for i, v := range n {
		ss[i] = v.String()
	}
	return ss
}

type NetworkSlaves []string

func (n *NetworkSlaves) Set(value string) error {
	*n = append(*n, value)
	return nil
}

func (n NetworkSlaves) String() string {
	return strings.Join(n, "|")
}
