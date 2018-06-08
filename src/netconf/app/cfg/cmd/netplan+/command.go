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
	"os/exec"

	"netconf/lib/netplan"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
)

const (
	DEFAULT_CFG_DIR = "/etc/netplan"
)

type NpCommand struct {
	dir     string
	verbose bool
	force   bool
}

func (c *NpCommand) init() *NpCommand {
	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}

	return c
}

func (c *NpCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "show detail message.")
	cmd.PersistentFlags().BoolVarP(&c.force, "force", "f", false, "ignore error.")
	return cmd
}

func (c *NpCommand) SetInitFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.dir, "dir", "d", DEFAULT_CFG_DIR, "config dir name.")
	return c.SetFlags(cmd)
}

func (c *NpCommand) Exec(cmd string, args ...string) ([]byte, error) {
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		log.Errorf("%s", out)
		if c.force {
			err = nil
		}
	}

	return out, err
}

func (c *NpCommand) SetMtu(ifname string, device *ncnplib.Device) error {
	if device.Mtu == 0 {
		log.Debugf("[%s] MTU: not changed.(%d)", ifname, device.Mtu)
		return nil
	}

	link, err := netlink.LinkByName(ifname)
	if err != nil {
		if c.force {
			log.Warnf("%s", err)
			return nil
		} else {
			return err
		}
	}

	linkMtu := link.Attrs().MTU

	if linkMtu == int(device.Mtu) {
		log.Debugf("[%s] MTU: not changed.(%d->%d)", ifname, linkMtu, device.Mtu)
		return nil
	}

	log.Debugf("[%s] MTU: %d -> %d", ifname, linkMtu, device.Mtu)

	if err := netlink.LinkSetMTU(link, int(device.Mtu)); err != nil {
		if c.force {
			log.Warnf("%s", err)
		} else {
			return err
		}
	}

	return nil
}

func (c *NpCommand) Init() error {
	return c.DoInit(false) // force flag: false
}

func (c *NpCommand) DoInit(force bool) error {
	cfg, err := ncnplib.ReadConfigDir(c.dir)
	if err != nil {
		log.Errorf("%s %s", c.dir, err)
		return err
	}

	for ifname, eth := range cfg.Network.Ethernets {
		log.Debugf("Ethernet[%s] %v", ifname, eth)
		if err := c.SetMtu(ifname, &eth.Device); err != nil {
			if !force {
				log.Errorf("%s %s", ifname, err)
				return err
			}
			log.Warnf("%s %s. but ignored.", ifname, err)
		}
	}

	for ifname, vlan := range cfg.Network.Vlans {
		log.Debugf("VLAN[%s] %v", ifname, vlan)
		if err := c.SetMtu(ifname, &vlan.Device); err != nil {
			if !force {
				log.Errorf("%s %s", ifname, err)
				return err
			}
			log.Warnf("%s %s. but ignored.", ifname, err)
		}
	}

	for ifname, bond := range cfg.Network.Bonds {
		log.Debugf("BOND[%s] %v", ifname, bond)
		if err := c.SetMtu(ifname, &bond.Device); err != nil {
			if !force {
				log.Errorf("%s %s", ifname, err)
				return err
			}
			log.Errorf("%s %s. but ignored.", ifname, err)
		}
	}

	return nil
}

func (c *NpCommand) Run(arg string) error {
	_, err := c.Exec("netplan", arg)
	return err
}

func (c *NpCommand) Apply() error {
	c.DoInit(true) // force flag: true

	if err := c.Run("apply"); err != nil {
		return err
	}

	return c.DoInit(false) // force flag: false
}
