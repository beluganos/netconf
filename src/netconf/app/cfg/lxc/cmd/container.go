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

type ContainerCommand struct {
	Command
	keep bool
}

func (c *ContainerCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().BoolVarP(&c.keep, "keep", "k", false, "keep file on error.")
	return c.Command.SetFlags(cmd)
}

func (c *ContainerCommand) Create(name string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		return err
	}

	if err := lib.CreateContainer(client, name, c.keep); err != nil {
		return err
	}

	return nil
}

func (c *ContainerCommand) Delete(name string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		return err
	}

	lib.DeleteContainer(client, name)

	return nil
}

func ContainerCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "container",
		Short: "LXD container commands.",
	}

	create := ContainerCommand{}
	c.AddCommand(create.SetFlags(
		&cobra.Command{
			Use:   "create [cotainer name]",
			Short: "Create container.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return create.Create(args[0])
			},
		},
	))

	delete := ContainerCommand{}
	c.AddCommand(delete.SetFlags(
		&cobra.Command{
			Use:   "delete [cotainer name]",
			Short: "Delete container.",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return delete.Delete(args[0])
			},
		},
	))

	return c
}
