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
	"fmt"
	lib "netconf/app/cfg/lxc/lib"
	"path"

	"github.com/spf13/cobra"
)

type ProfileCommand struct {
	Command
	backupDir string
}

func (c *ProfileCommand) BackupPath(name string) string {
	return path.Join(c.backupDir, fmt.Sprintf("backup_lxd_profile_%s.yml", name))
}

func (c *ProfileCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.backupDir, "backup", "b", lib.PROFILE_BACKUP_DIR, "backuo directory")
	return c.Command.SetFlags(cmd)
}

func (c *ProfileCommand) Load(name string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		return err
	}

	if err := lib.LoadProfile(client, name, c.BackupPath(name)); err != nil {
		return err
	}

	return nil
}

func (c *ProfileCommand) Save(name string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		return err
	}

	if err := lib.SaveProfile(client, name, c.BackupPath(name)); err != nil {
		return err
	}

	return nil
}

func (c *ProfileCommand) Clear(name string) error {
	c.Command.Init()

	if err := lib.ClearBackup(c.BackupPath(name)); err != nil {
		return err
	}

	return nil
}

func ProfileCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "profile",
		Short: "LXD profile commands.",
	}

	load := ProfileCommand{}
	c.AddCommand(load.SetFlags(
		&cobra.Command{
			Use:   "load [cotainer name]",
			Short: "Load profile from file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return load.Load(args[0])
			},
		},
	))

	save := ProfileCommand{}
	c.AddCommand(save.SetFlags(
		&cobra.Command{
			Use:   "save [cotainer name]",
			Short: "Save profile to file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return save.Save(args[0])
			},
		},
	))

	clear := ProfileCommand{}
	c.AddCommand(clear.SetFlags(
		&cobra.Command{
			Use:   "clear [cotainer name]",
			Short: "Clear backup profile file.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return clear.Save(args[0])
			},
		},
	))

	return c
}
