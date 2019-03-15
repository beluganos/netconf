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

package cfglxccmd

import (
	lib "netconf/app/cfg/lxc/lib"

	"github.com/spf13/cobra"
)

const InterfaceDefaultMTU = 9000

type InterfaceCommand struct {
	Command
	keep   bool
	prefix string
	mtu    uint16
	negate bool
}

func (c *InterfaceCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.keep, "keep", "k", false, "keep file on error.")
	cmd.PersistentFlags().StringVarP(&c.prefix, "prefix", "p", "v", "Prefix of ifname on host.")
	cmd.PersistentFlags().Uint16VarP(&c.mtu, "mtu", "m", InterfaceDefaultMTU, "MTU of NIC devices.")
	return c.Command.SetFlags(cmd)
}

func (c *InterfaceCommand) SetSetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.keep, "keep", "k", false, "keep file on error.")
	cmd.PersistentFlags().BoolVarP(&c.negate, "no", "n", false, "negate command.")
	return c.Command.SetFlags(cmd)
}

func (c *InterfaceCommand) Add(name string, ifname string, hwaddr string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		return err
	}

	if err := lib.AddInterface(client, name, c.prefix, ifname, hwaddr, c.mtu); err != nil {
		return err
	}

	return nil
}

func (c *InterfaceCommand) Delete(name string, ifname string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		return err
	}

	if err := lib.DeleteInterface(client, name, ifname); err != nil {
		return err
	}

	return nil
}

func (c *InterfaceCommand) Set(name string, ifname string, args []string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		return err
	}

	if err := lib.SetInterface(client, name, ifname, c.negate, args...); err != nil {
		return err
	}

	return nil
}

func InterfaceCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "interface",
		Short: "LXD interface commands.",
	}

	add := InterfaceCommand{}
	c.AddCommand(add.SetFlags(
		&cobra.Command{
			Use:   "add [cotainer name] [ifname] [hwaddr]",
			Short: "Add interface to container.",
			Args:  cobra.RangeArgs(2, 3),
			RunE: func(cmd *cobra.Command, args []string) error {
				hwaddr := func() string {
					if len(args) == 2 {
						return ""
					}
					return args[2]
				}()
				return add.Add(args[0], args[1], hwaddr)
			},
		},
	))

	delete := InterfaceCommand{}
	c.AddCommand(delete.SetFlags(
		&cobra.Command{
			Use:   "delete [cotainer name] [ifname]",
			Short: "Delete interface from container.",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return delete.Delete(args[0], args[1])
			},
		},
	))

	set := InterfaceCommand{}
	c.AddCommand(set.SetSetFlags(
		&cobra.Command{
			Use:   "set [container name] [ifname] [key=value...]",
			Short: "Set interface config",
			Args:  cobra.MinimumNArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				return set.Set(args[0], args[1], args[2:])
			},
		},
	))

	return c
}
