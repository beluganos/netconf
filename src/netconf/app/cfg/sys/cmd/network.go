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

package cfgsyscmd

import (
	api "netconf/app/cfg/api"
	lib "netconf/app/cfg/sys/lib"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

type NetworkCommand struct {
	api.Command
	path   string
	wait   time.Duration
	negate bool
	mtu    uint16
	addrs  lib.NetowrkAddrs
}

func parseVid(s string) (uint, error) {
	if len(s) == 0 {
		return 0, nil
	}

	v, err := strconv.ParseUint(s, 0, 16)
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (c *NetworkCommand) BackupPath() string {
	return lib.NetworkBackupPath(c.path)
}

func (c *NetworkCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "p", lib.NETPLAN_CONF_PATH, "config filename")
	cmd.PersistentFlags().DurationVar(&c.wait, "wait", lib.NETPLAN_WAIT_SEC, "wait seconds after load.")

	return c.Command.SetFlags(cmd)
}

func (c *NetworkCommand) SetEthFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().Uint16Var(&c.mtu, "mtu", 0, "MTU.")
	cmd.PersistentFlags().VarP(&c.addrs, "addr", "a", "Interface address.")
	cmd.PersistentFlags().BoolVarP(&c.negate, "negate", "n", false, "Negate command")
	return c.SetFlags(cmd)
}

func (c *NetworkCommand) DoNetwork(device string, vid string, slaves []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	v, err := parseVid(vid)
	if err != nil {
		return err
	}

	cmd := func() string {
		if c.negate {
			return "del"
		} else {
			return "set"
		}
	}()

	res, err := lib.DoNetworkRun(cmd, device, v, uint(c.mtu), c.addrs.Strings(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *NetworkCommand) Backup(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.BackupNetwotkRun(c.path, c.BackupPath(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *NetworkCommand) Rollback(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.BackupNetwotkRun(c.BackupPath(), c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *NetworkCommand) Load(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.LoadNetworkRun(c.wait, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func NetworkCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "network",
		Short: "Neteork configuration file commands.",
	}

	c_set := &cobra.Command{
		Use:   "set",
		Short: "Edit configuration file commands.",
	}
	c.AddCommand(c_set)

	eth := NetworkCommand{}
	c_set.AddCommand(eth.SetEthFlags(
		&cobra.Command{
			Use:   "ethernet [device]",
			Short: "Ethernet device configuration.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return eth.DoNetwork(args[0], "", []string{})
			},
		},
	))

	vlan := NetworkCommand{}
	c_set.AddCommand(vlan.SetEthFlags(
		&cobra.Command{
			Use:   "vlan [device] [vid]",
			Short: "VLAN device configuration.",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return vlan.DoNetwork(args[0], args[1], []string{})
			},
		},
	))

	bond := NetworkCommand{}
	c_set.AddCommand(bond.SetEthFlags(
		&cobra.Command{
			Use:   "bond [device] [slave...]",
			Short: "Bonding device configuration.",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return bond.DoNetwork(args[0], "", args[1:])
			},
		},
	))

	load := NetworkCommand{}
	c.AddCommand(load.SetFlags(
		&cobra.Command{
			Use:   "load",
			Short: "Apply setting.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return load.Load(args)
			},
		},
	))

	backup := NetworkCommand{}
	c.AddCommand(backup.SetFlags(
		&cobra.Command{
			Use:   "backup",
			Short: "Backup configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return backup.Backup(args)
			},
		},
	))

	rollback := NetworkCommand{}
	c.AddCommand(rollback.SetFlags(
		&cobra.Command{
			Use:   "rollback",
			Short: "Replace configuration file to backup.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return rollback.Rollback(args)
			},
		},
	))

	return c
}
