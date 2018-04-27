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

type IPv6Comand struct {
	api.Command
	negate bool
}

func (c *IPv6Comand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.negate, "negate", "n", false, "Negate command")
	return c.Command.SetFlags(cmd)
}

func (c *IPv6Comand) IPv6(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SetIPv6(c.negate, args, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil

}

func IPv6Cmd() *cobra.Command {
	ipv6 := IPv6Comand{}
	c := ipv6.SetFlags(
		&cobra.Command{
			Use:   "ipv6 [command...]",
			Short: "IPv6 configuration",
			RunE: func(cmd *cobra.Command, args []string) error {
				return ipv6.IPv6(args)
			},
		},
	)

	return c
}
