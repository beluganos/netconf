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
	"fmt"
	"net"
	"strings"
)

//
// Interface name
//
func NewIFName(dev string, index uint32) string {
	if index == 0 {
		return dev
	} else {
		return fmt.Sprintf("%s.%d", dev, index)
	}
}

func ParseIFName(name string) (string, uint32, error) {
	s := strings.Replace(name, ".", " ", 1)
	var dev string = ""
	var index uint32 = 0
	if n, err := fmt.Sscanf(s, "%s %d", &dev, &index); n == 0 {
		return dev, index, err
	}

	return dev, index, nil
}

//
// router-id
//
type RouterId struct {
	net.IP
}

func NewRouterId(ip net.IP) *RouterId {
	return &RouterId{
		IP: ip,
	}
}

func ParseRouterId(s string) (*RouterId, error) {
	if ip := net.ParseIP(s); ip != nil {
		return NewRouterId(ip), nil
	}
	return nil, fmt.Errorf("Invalid router id. %s", s)
}

func (r *RouterId) String() string {
	return r.IP.String()
}

func (r *RouterId) IPNet() *net.IPNet {
	return IPToIPNet(r.IP, IPBITS_FULL)
}
