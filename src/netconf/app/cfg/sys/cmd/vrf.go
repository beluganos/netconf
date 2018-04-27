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

	"github.com/spf13/cobra"
)

type VrfCommand struct {
	api.Command
	path   string
	negate bool
}

func (c *VrfCommand) BackupPath() string {
	return lib.VrfBackupPath(c.path)
}

func (c *VrfCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "p", lib.VRF_CONF_PATH, "config filename")
	return c.Command.SetFlags(cmd)
}

func (c *VrfCommand) SetModFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.negate, "negate", "n", false, "Negate command")
	return c.SetFlags(cmd)
}

func (c *VrfCommand) DoVrf(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	cmd := func() string {
		if c.negate {
			return "del"
		} else {
			return "set"
		}
	}()

	res, err := lib.DoVrfExec(cmd, args, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *VrfCommand) Load() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.LoadVrfExec(client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *VrfCommand) Backup() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.BackupVrfExec(c.path, c.BackupPath(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *VrfCommand) Rollback() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.RollbackVrfExec(c.BackupPath(), c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *VrfCommand) Commit() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.CommitVrfExec(c.BackupPath(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func VrfCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "vrf",
		Short: "vrf configuration commands.",
	}

	cfg := VrfCommand{}
	c.AddCommand(cfg.SetModFlags(
		&cobra.Command{
			Use:   "set ['key=value']",
			Short: "Edit configuration file.",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return cfg.DoVrf(args)
			},
		},
	))

	load := VrfCommand{}
	c.AddCommand(load.SetFlags(
		&cobra.Command{
			Use:   "load",
			Short: "Apply configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return load.Load()
			},
		},
	))

	backup := VrfCommand{}
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

	rollback := VrfCommand{}
	c.AddCommand(rollback.SetFlags(
		&cobra.Command{
			Use:   "rollback",
			Short: "Rollback configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return rollback.Rollback()
			},
		},
	))

	commit := VrfCommand{}
	c.AddCommand(commit.SetFlags(
		&cobra.Command{
			Use:   "commit",
			Short: "Commit configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return commit.Rollback()
			},
		},
	))

	return c
}
