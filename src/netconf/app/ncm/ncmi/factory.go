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
	"fmt"
	"net"
	ncianalib "netconf/lib/iana"
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"
)

type InterfaceFactory struct {
	LXCIFNameFmt string
}

func NewInterfaceFactory(lxc string) *InterfaceFactory {
	return &InterfaceFactory{
		LXCIFNameFmt: lxc,
	}
}

func (f *InterfaceFactory) NewIfname(portId PortId) string {
	return fmt.Sprintf(f.LXCIFNameFmt, portId)
}

func (f *InterfaceFactory) NewInterfaceConfig(portId PortId) *openconfig.InterfaceConfig {
	ifname := f.NewIfname(portId)
	ifconf := openconfig.NewInterfaceConfig()
	ifconf.SetName(ifname)
	ifconf.SetType(ncianalib.IANAifType_ethernetCsmacd)
	return ifconf
}

func (f *InterfaceFactory) NewInterfaceEthernet(hwAddr net.HardwareAddr) *openconfig.InterfaceEthernet {
	config := openconfig.NewInterfaceEthernetConfig()
	config.SetMacAddr(hwAddr)
	ethernet := openconfig.NewInterfaceEthernet()
	ethernet.SetConfig(config)
	return ethernet
}

func (f *InterfaceFactory) NewInterface(port *PortEntry) *openconfig.Interface {
	iface := openconfig.NewInterface("")
	iface.SetConfig(f.NewInterfaceConfig(port.PortId))
	if mac, err := net.ParseMAC(port.HwAddr); err == nil {
		iface.SetEthernet(f.NewInterfaceEthernet(mac))
	}
	return iface
}

func (f *InterfaceFactory) NewInterfaces(ports []*PortEntry) openconfig.Interfaces {
	ifaces := openconfig.NewInterfaces()
	for _, port := range ports {
		iface := f.NewInterface(port)
		ifaces[iface.Name] = iface
	}
	return ifaces
}

func (f *InterfaceFactory) NewInterfacesVals(ports []*PortEntry) ([]*srlib.SrVal, error) {
	ifvals := NewInterfaceVals()
	ifaces := f.NewInterfaces(ports)
	if err := openconfig.ProcessInterfaces(ifvals, false, ifaces); err != nil {
		return nil, err
	}
	return ifvals.Vals, nil
}
