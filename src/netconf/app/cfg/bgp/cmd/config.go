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

package cfgbgpcmd

import (
	api "netconf/app/cfg/api"
	lib "netconf/app/cfg/bgp/lib"
	"time"

	"github.com/spf13/cobra"
)

type GoBgpCommand struct {
	api.Command
	path   string
	wait   time.Duration
	negate bool
}

func (c *GoBgpCommand) BackupPath() string {
	return lib.GobgpBackupPath(c.path)
}

func (c *GoBgpCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	c.Command.SetFlags(cmd)
	return cmd
}

func (c *GoBgpCommand) SetModFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "p", lib.GOBGP_CONF_PATH, "config filename")
	cmd.PersistentFlags().BoolVarP(&c.negate, "negate", "n", false, "Negate command")
	return c.SetFlags(cmd)
}

func (c *GoBgpCommand) SetFileFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "p", lib.GOBGP_CONF_PATH, "config filename")
	return c.SetFlags(cmd)
}

func (c *GoBgpCommand) SetLoadFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().DurationVar(&c.wait, "wait", lib.GOBGP_WAIT_SEC, "wait seconds after load.")
	return c.SetFlags(cmd)
}

func (c *GoBgpCommand) DoGobgp(args []string) error {
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

	res, err := lib.DoGobgpRun(cmd, c.path, args, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *GoBgpCommand) Load() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.LoadGobgpRun(c.wait, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *GoBgpCommand) Backup() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.BackupGobgpRun(c.path, c.BackupPath(), client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *GoBgpCommand) Rollback() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.RollbackGobgpRun(c.BackupPath(), c.path, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *GoBgpCommand) Commit() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.CommitGobgpRun(c.BackupPath(), c.wait, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func (c *GoBgpCommand) Restart() error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.RestartGobgpRun(client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func ConfigCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "GoBGP configuration commands.",
	}

	cfg := GoBgpCommand{}
	c.AddCommand(cfg.SetModFlags(
		&cobra.Command{
			Use:   "set [filepath or '-']",
			Short: "Merge configuration file.",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return cfg.DoGobgp(args)
			},
		},
	))

	load := GoBgpCommand{}
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

	backup := GoBgpCommand{}
	c.AddCommand(backup.SetFileFlags(
		&cobra.Command{
			Use:   "backup",
			Short: "Backup configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return backup.Backup()
			},
		},
	))

	rollback := GoBgpCommand{}
	c.AddCommand(rollback.SetFileFlags(
		&cobra.Command{
			Use:   "rollback",
			Short: "Rollback configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return rollback.Rollback()
			},
		},
	))

	commit := GoBgpCommand{}
	c.AddCommand(commit.SetModFlags(
		&cobra.Command{
			Use:   "commit",
			Short: "Commit configuration file.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return commit.Commit()
			},
		},
	))

	restart := GoBgpCommand{}
	c.AddCommand(restart.SetModFlags(
		&cobra.Command{
			Use:   "restart",
			Short: "Restart gobgpd.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return restart.Restart()
			},
		},
	))

	return c
}
