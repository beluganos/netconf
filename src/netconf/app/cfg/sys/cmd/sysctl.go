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
	"time"

	"github.com/spf13/cobra"
)

type SysctlCommand struct {
	api.Command
	path   string
	wait   time.Duration
	negate bool
}

func (c *SysctlCommand) BackupPath() string {
	return lib.SysctlBackupPath(c.path)
}

func (c *SysctlCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "p", lib.SYSCTL_CONF_PATH, "config filename")
	return c.Command.SetFlags(cmd)
}

func (c *SysctlCommand) SetModFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.negate, "negate", "n", false, "Negate command")
	return c.SetFlags(cmd)
}

func (c *SysctlCommand) SetLoadFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().DurationVar(&c.wait, "wait", lib.SYSCTL_WAIT_SEC, "wait seconds after load.")
	return c.SetFlags(cmd)
}

func (c *SysctlCommand) DoSysctl(args []string) error {
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

	res, err := lib.DoSysctlRun(cmd, args, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *SysctlCommand) Load() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.LoadSysctlRun(c.path, c.wait, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *SysctlCommand) Backup() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.BackupSysctlRun(c.path, c.BackupPath(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *SysctlCommand) Rollback() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.RollbackSysctlRun(c.BackupPath(), c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *SysctlCommand) Commit() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.CommitSysctlRun(c.path, c.BackupPath(), c.wait, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func SysctlCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "sysctl",
		Short: "Sysctl configuration commands.",
	}

	cfg := SysctlCommand{}
	c.AddCommand(cfg.SetModFlags(
		&cobra.Command{
			Use:   "set ['key=value']",
			Short: "Edit configuration file.",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return cfg.DoSysctl(args)
			},
		},
	))

	load := SysctlCommand{}
	c.AddCommand(load.SetLoadFlags(
		&cobra.Command{
			Use:   "load",
			Short: "Apply configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return load.Load()
			},
		},
	))

	backup := SysctlCommand{}
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

	rollback := SysctlCommand{}
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

	commit := SysctlCommand{}
	c.AddCommand(commit.SetLoadFlags(
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
