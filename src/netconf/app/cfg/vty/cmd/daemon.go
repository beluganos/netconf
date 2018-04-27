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

type DaemonCommand struct {
	api.Command
}

func (c *DaemonCommand) Enable(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SetDaemonRun(args, "set", client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *DaemonCommand) Disable(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SetDaemonRun(args, "del", client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *DaemonCommand) Show() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.DaemonConfigRun(client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func DaemonCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "daemon",
		Short: "daemons configuration file commands.",
	}

	enable := DaemonCommand{}
	c.AddCommand(enable.SetFlags(
		&cobra.Command{
			Use:   "enable [daemon...]",
			Short: "Enable daemon(s).",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return enable.Enable(args)
			},
		},
	))

	disable := DaemonCommand{}
	c.AddCommand(disable.SetFlags(
		&cobra.Command{
			Use:   "disable [daemon...]",
			Short: "Disable daemon(s).",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return disable.Disable(args)
			},
		},
	))

	show := DaemonCommand{}
	c.AddCommand(show.SetFlags(
		&cobra.Command{
			Use:   "show",
			Short: "Show daemons settings.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return show.Show()
			},
		},
	))

	return c
}
