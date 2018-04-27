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

package cfgvtycmd

import (
	api "netconf/app/cfg/api"
	lib "netconf/app/cfg/vty/lib"

	"github.com/spf13/cobra"
)

type MplsCommand struct {
	api.Command
	negate bool
}

func (c *MplsCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.negate, "negate", "n", false, "Negate command")
	return c.Command.SetFlags(cmd)
}

func (c *MplsCommand) Ldp(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SetMplsLdpRun(c.negate, 0, "", args, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *MplsCommand) AddrssFamily(ipver uint, ifname string, args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SetMplsLdpRun(c.negate, ipver, ifname, args, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func MplsCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "mpls",
		Short: "MPLS LDP configuration commands.",
	}

	// mpls ldp
	// commands...
	ldp := MplsCommand{}
	c_ldp := ldp.SetFlags(&cobra.Command{
		Use:   "ldp [command...]",
		Short: "MPLS LDP configuration.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ldp.Ldp(args)
		},
	})
	c.AddCommand(c_ldp)

	// mpls ldp
	// address-family ipv4/6
	// commands...
	c_af := &cobra.Command{
		Use:   "address-family",
		Short: "MPLS LDP(address-family) configuration.",
	}
	c_ldp.AddCommand(c_af)

	ipv4 := MplsCommand{}
	c_ipv4 := ipv4.SetFlags(
		&cobra.Command{
			Use:   "ipv4 [command...]",
			Short: "MPLS LDP(IPv4) configuration.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return ipv4.AddrssFamily(4, "", args)
			},
		},
	)
	c_af.AddCommand(c_ipv4)

	ifv4 := MplsCommand{}
	c_ifv4 := ifv4.SetFlags(
		&cobra.Command{
			Use:   "interface [ifname] [command...]",
			Short: "MPLS LDP(IPv4) configuration.",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return ifv4.AddrssFamily(4, args[0], args[1:])
			},
		},
	)
	c_ipv4.AddCommand(c_ifv4)

	ipv6 := MplsCommand{}
	c_ipv6 := ipv6.SetFlags(
		&cobra.Command{
			Use:   "ipv6 [command...]",
			Short: "MPLS LDP(IPv6) configuration.",
			RunE: func(cmd *cobra.Command, args []string) error {
				return ipv6.AddrssFamily(6, "", args)
			},
		},
	)
	c_af.AddCommand(c_ipv6)

	ifv6 := MplsCommand{}
	c_ifv6 := ifv6.SetFlags(
		&cobra.Command{
			Use:   "interface [ifname] [command...]",
			Short: "MPLS LDP interface (IPv6) configuration.",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return ifv6.AddrssFamily(6, args[0], args[1:])
			},
		},
	)
	c_ipv6.AddCommand(c_ifv6)

	return c
}
