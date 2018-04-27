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

package ncnplib

const YAML_VERSION = 2

type Device struct {
	Dhcp4     bool     `yaml:"dhcp4"`
	Dhcp6     bool     `yaml:"dhcp6"`
	Addresses []string `yaml:"addresses"`
	Gateway4  string   `yaml:"gateway4,omitempty"`
	Gateway6  string   `yaml:"gateway6,omitempty"`
	Mtu       uint16   `yaml:"mtu,omitempty"`
	// AcceptRA  bool     `yaml:"accept-ra"`
}
type Ethernet struct {
	Device `yaml:",inline"`
}

func NewEthernet(device *Device) *Ethernet {
	return &Ethernet{
		Device: *device,
	}
}

type Vlan struct {
	Device `yaml:",inline"`
	Link   string `yaml:"link"`
	Id     uint32 `yaml:"id"`
}

func NewVlan(device *Device, link string, id uint32) *Vlan {
	return &Vlan{
		Device: *device,
		Link:   link,
		Id:     id,
	}
}

type BondParams struct {
	Mode                  BondMode                  `yaml:"mode"`
	LacpRate              uint32                    `yaml:"lacp-rate"`
	MiiMonitorInterval    uint32                    `yaml:"mii-monitor-interval"`
	MinLinks              uint32                    `yaml:"min-links"`
	TransmitHashPolicy    BondTransmitHashPolicy    `yaml:"transmit-hash-policy"`
	AdSelect              BondAdSelect              `yaml:"ad-select"`
	AllSlavesActive       bool                      `yaml:"all-slaves-active"`
	ARPInterval           uint32                    `yaml:"arp-interval"`
	ARPIpTargets          []string                  `yaml:"arp-ip-targets"`
	ARPValidate           BondARPValidate           `yaml:"arp-validate"`
	ARPAllTargets         BondARPAllTargets         `yaml:"arp-all-targets"`
	UpDelay               uint32                    `yaml:"up-delay"`
	DownDelay             uint32                    `yaml:"down-delay"`
	FailOverMACPolicy     BondFailOverMACPolicy     `yaml:"fail-over-mac-policy"`
	GratuitiousARP        uint32                    `yaml:"gratuitious-arp"`
	PacketsPerSlave       uint16                    `yaml:"packets-per-slave"`
	PrimaryReselectPolicy BondPrimaryReselectPolicy `yaml:"primary-reselect-policy"`
	LearnPacketInterval   uint32                    `yaml:"learn-packet-interval"`
	Primary               string                    `yaml:"primary"`
}

type Bond struct {
	Device     `yaml:",inline"`
	Interfaces []string   `yaml:"interfaces"`
	Params     BondParams `yaml:"parameters"`
}

func NewBond(device *Device, ifaces []string, params *BondParams) *Bond {
	return &Bond{
		Device:     *device,
		Interfaces: ifaces,
		Params:     *params,
	}
}

type Network struct {
	Version   uint32               `yaml:"version"`
	Renderer  string               `yaml:"renderer,omitempty"`
	Ethernets map[string]*Ethernet `yaml:"ethernets"`
	Vlans     map[string]*Vlan     `yaml:"vlans"`
	Bonds     map[string]*Bond     `yaml:"bonds"`
}

func NewNetwork() *Network {
	return &Network{
		Version:   YAML_VERSION,
		Ethernets: map[string]*Ethernet{},
		Vlans:     map[string]*Vlan{},
		Bonds:     map[string]*Bond{},
	}
}

func (n *Network) Merge(nw *Network) {
	for ifname, eth := range nw.Ethernets {
		n.Ethernets[ifname] = eth
	}
	for ifname, vlan := range nw.Vlans {
		n.Vlans[ifname] = vlan
	}
	for ifname, bond := range nw.Bonds {
		n.Bonds[ifname] = bond
	}
}

type Config struct {
	Network Network `yaml:"network"`
}

func (c *Config) Merge(cfg *Config) {
	c.Network.Merge(&cfg.Network)
}

func NewConfig() *Config {
	return &Config{
		Network: *NewNetwork(),
	}
}
