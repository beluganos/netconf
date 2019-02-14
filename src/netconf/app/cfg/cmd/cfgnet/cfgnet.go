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
	ncnplib "netconf/lib/netplan"

	log "github.com/sirupsen/logrus"
)

func deleteSlice(slice []string, strs ...string) []string {
	m := map[string]struct{}{}
	for _, s := range slice {
		m[s] = struct{}{}
	}
	for _, s := range strs {
		delete(m, s)
	}
	result := []string{}
	for k, _ := range m {
		result = append(result, k)
	}
	return result
}

func mergeDevice(device *ncnplib.Device, src *ncnplib.Device) {
	for _, address := range src.Addresses {
		device.Addresses = append(device.Addresses, address)
		log.Debugf("Ethernet/IP = %s", address)
	}
	device.Addresses = deleteSlice(device.Addresses)

	if mtu := uint16(src.Mtu); mtu != 0 {
		device.Mtu = mtu
		log.Debugf("Ethernet/MTU = %d", mtu)
	}
}

func mergeEthernet(ethernet *ncnplib.Ethernet, src *ncnplib.Ethernet) {
	mergeDevice(&ethernet.Device, &src.Device)
}

func mergeVlan(vlan *ncnplib.Vlan, src *ncnplib.Vlan) {
	mergeDevice(&vlan.Device, &src.Device)

	if link := src.Link; len(link) != 0 {
		vlan.Link = link
		log.Debugf("VLAN/Link = %s", link)
	}

	if id := src.Id; id != 0 {
		vlan.Id = id
		log.Debugf("VLAN/Id = %d", id)
	}
}

func setConfig(cfg *ncnplib.Config, args *Args) error {
	ifname := args.IFName()
	device := &ncnplib.Device{
		Addresses: args.Addrs.Strings(),
		Mtu:       uint16(args.Mtu),
	}

	if vid := uint32(args.Vid); vid == 0 {
		src := ncnplib.NewEthernet(device)
		if ethernet, ok := cfg.Network.Ethernets[ifname]; ok {
			mergeEthernet(ethernet, src)
		} else {
			cfg.Network.Ethernets[ifname] = src
		}
	} else {
		src := ncnplib.NewVlan(device, args.Device, vid)
		if vlan, ok := cfg.Network.Vlans[ifname]; ok {
			mergeVlan(vlan, src)
		} else {
			cfg.Network.Vlans[ifname] = src
		}
	}
	return nil
}

func deleteDevice(device *ncnplib.Device, src *ncnplib.Device) {
	if len(src.Addresses) != 0 {
		device.Addresses = deleteSlice(device.Addresses, src.Addresses...)
		log.Debugf("Ethernet/IP = %s DELETED", src.Addresses)
	}

	if src.Mtu != 0 {
		device.Mtu = 0
		log.Debugf("Ethernet/MTU = %d DELETED", src.Mtu)
	}
}

func deleteEthernet(ethernet *ncnplib.Ethernet, src *ncnplib.Ethernet) {
	deleteDevice(&ethernet.Device, &src.Device)
}

func deleteVlan(vlan *ncnplib.Vlan, src *ncnplib.Vlan) {
	deleteDevice(&vlan.Device, &src.Device)

	/*
	   if link := src.Link; len(link) != 0 {
	       vlan.Link = ""
	       log.Debugf("SetNetwork DEL Ethernet(%s)/Link = %s", name, link)
	   }

	   if id := src.Id; id != 0 {
	       vlan.Id = 0
	       log.Debugf("SetNetwork DEL Ethernet(%s)/Id = %d", name, id)
	   }
	*/
}

func delConfig(cfg *ncnplib.Config, args *Args) error {
	ifname := args.IFName()
	device := &ncnplib.Device{
		Addresses: args.Addrs.Strings(),
		Mtu:       uint16(args.Mtu),
	}

	if vid := uint32(args.Vid); vid == 0 {
		if ethernet, ok := cfg.Network.Ethernets[ifname]; ok {
			if len(args.Addrs) == 0 && args.Mtu == 0 {
				delete(cfg.Network.Ethernets, ifname)
			} else {
				deleteEthernet(ethernet, ncnplib.NewEthernet(device))
			}
		} else {
			log.Warnf("%s not found.", ifname)
		}
	} else {
		if vlan, ok := cfg.Network.Vlans[ifname]; ok {
			if len(args.Addrs) == 0 && args.Mtu == 0 {
				delete(cfg.Network.Vlans, ifname)
			} else {
				deleteVlan(vlan, ncnplib.NewVlan(device, args.Device, vid))
			}
		} else {
			log.Warnf("%s not found.", ifname)
		}
	}

	return nil
}
