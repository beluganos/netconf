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

package cfgbelcmd

import (
	api "netconf/app/cfg/api"
	lib "netconf/app/cfg/bel/lib"

	"github.com/spf13/cobra"
)

type RibsCommand struct {
	api.Command
	path string
	rt   string
	rd   string
}

func (c *RibsCommand) BackupPath() string {
	return lib.RibxBackupPath(c.path)
}

func (c *RibsCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "p", lib.RIBX_CONF_PATH, "config filename")
	return c.Command.SetFlags(cmd)
}

func (c *RibsCommand) SetVrfFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.rt, "RT", "t", "", "route target")
	cmd.PersistentFlags().StringVarP(&c.rd, "RD", "d", "", "route distinguisher")
	return c.SetFlags(cmd)
}

func (c *RibsCommand) SetVrf() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SetRibsVrf(c.rt, c.rd, c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *RibsCommand) Backup() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.BackupRibx(c.path, c.BackupPath(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *RibsCommand) Rollback() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.RollbackRibx(c.BackupPath(), c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *RibsCommand) Load() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.LoadRibs(client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func RibsCmd() *cobra.Command {

	c := &cobra.Command{
		Use:   "ribs",
		Short: "Ribs configuration file commands.",
	}

	set := &cobra.Command{
		Use:   "set",
		Short: "Ribs set configuration.",
	}
	c.AddCommand(set)

	vrf := RibsCommand{}
	set.AddCommand(vrf.SetVrfFlags(
		&cobra.Command{
			Use:   "vrf",
			Short: "Set Vrf config value.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vrf.SetVrf()
			},
		},
	))

	load := RibsCommand{}
	c.AddCommand(load.SetFlags(
		&cobra.Command{
			Use:   "load",
			Short: "Apply setting.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return load.Load()
			},
		},
	))

	backup := RibsCommand{}
	c.AddCommand(backup.SetFlags(
		&cobra.Command{
			Use:   "backup",
			Short: "Backup configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return backup.Backup()
			},
		},
	))

	rollback := RibsCommand{}
	c.AddCommand(rollback.SetFlags(
		&cobra.Command{
			Use:   "rollback",
			Short: "rollback configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return rollback.Rollback()
			},
		},
	))

	return c
}
