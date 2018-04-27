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

package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

const (
	CMD_SET = "set"
	CMD_DEL = "del"
)

const NETPLAN_CONF_PATH = "/etc/netplan/02-beluganos.yaml"

type Args struct {
	Cmd     string
	Path    string
	Backup  string
	Device  string
	Vid     uint
	Mtu     uint
	Addrs   Addrs
	Slaves  Slaves
	Verbose bool
	Args    []string
}

func (a *Args) Parse() error {
	flag.StringVar(&a.Cmd, "cmd", "", "'set' or 'del'")
	flag.StringVar(&a.Path, "path", NETPLAN_CONF_PATH, "Config filename")
	flag.StringVar(&a.Backup, "backup", "", "Backup filename")
	flag.StringVar(&a.Device, "device", "", "Device name")
	flag.UintVar(&a.Vid, "vid", 0, "VLAN-ID")
	flag.UintVar(&a.Mtu, "mtu", 0, "MTU")
	flag.Var(&a.Addrs, "a", "Interface addresses (ip/prefix-len)")
	flag.Var(&a.Slaves, "s", "Slave interfaces")
	flag.BoolVar(&a.Verbose, "v", false, "show detail message")
	flag.Parse()
	a.Args = flag.Args()

	if len(a.Device) == 0 {
		return fmt.Errorf("Invalid argument(path).")
	}

	return nil
}

func (a *Args) BackupPath() string {
	if len(a.Backup) == 0 {
		return fmt.Sprintf("%s.backup", a.Path)
	}

	return a.Backup
}

func (a *Args) IFName() string {
	if a.Vid == 0 {
		return a.Device
	}

	return fmt.Sprintf("%s.%d", a.Device, a.Vid)
}

type Addrs []*net.IPNet

func (n *Addrs) Set(value string) error {
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

func (n Addrs) String() string {
	return strings.Join(n.Strings(), "|")
}

func (n Addrs) Strings() []string {
	ss := make([]string, len(n))
	for i, v := range n {
		ss[i] = v.String()
	}
	return ss
}

type Slaves []string

func (n *Slaves) Set(value string) error {
	*n = append(*n, value)
	return nil
}

func (n Slaves) String() string {
	return strings.Join(n, "|")
}
