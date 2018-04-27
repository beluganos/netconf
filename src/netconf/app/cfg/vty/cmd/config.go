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

type ConfigCommand struct {
	api.Command
	path string
}

func (c *ConfigCommand) backupPath() string {
	return lib.ConfigBackupPath(c.path)
}

func (c *ConfigCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	c.Command.SetFlags(cmd)
	return cmd
}

func (c *ConfigCommand) SetPathFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "p", lib.FRR_CONF_PATH, "Host port")
	return c.SetFlags(cmd)
}

func (c *ConfigCommand) ShowBackup() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.ShowBackupConfigRun(c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *ConfigCommand) Save() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SaveConfigRun(client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *ConfigCommand) Backup() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.BackupConfigRun(c.path, c.backupPath(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *ConfigCommand) Rollback() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.RollbackConfigRun(c.backupPath(), c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func ConfigCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "frr configuration file commands.",
	}

	show := ConfigCommand{}
	c.AddCommand(show.SetPathFlags(
		&cobra.Command{
			Use:   "show-backup",
			Short: "Show backup configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return show.ShowBackup()
			},
			ValidArgs: []string{"<cr>"},
		},
	))

	save := ConfigCommand{}
	c.AddCommand(save.SetFlags(
		&cobra.Command{
			Use:   "save",
			Short: "Save running configuration to file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return save.Save()
			},
			ValidArgs: []string{"<cr>"},
		},
	))

	backup := ConfigCommand{}
	c.AddCommand(backup.SetPathFlags(
		&cobra.Command{
			Use:   "backup",
			Short: "Backup configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return backup.Backup()
			},
			ValidArgs: []string{"<cr>"},
		},
	))

	rollback := ConfigCommand{}
	c.AddCommand(rollback.SetPathFlags(
		&cobra.Command{
			Use:   "rollback",
			Short: "Replace configuration file to backup.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return rollback.Rollback()
			},
			ValidArgs: []string{"<cr>"},
		},
	))

	return c
}
