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

package ncccmd

import (
	nc "github.com/hiepon/go-netconf/netconf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type EditConfigCommand struct {
	Command
	Ope    string
	Commit bool
}

func (c *EditConfigCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.Ope, "operation", "d", "replace", "default-operation(replace/none/merge).")
	cmd.PersistentFlags().BoolVarP(&c.Commit, "commit", "c", false, "Commit on success.")
	return c.Command.SetFlags(cmd)
}

func (c *EditConfigCommand) EditConfig(target string, input string) error {
	c.Command.Init()

	data, err := inputData(input)
	if err != nil {
		return err
	}

	client, err := c.Client()
	if err != nil {
		log.Errorf("Client initialize error. %s", err)
		return err
	}
	defer client.Close()

	if _, err := client.Exec(nc.MethodEditConfig(target, c.Ope, string(data))); err != nil {
		log.Errorf("Exec error. %s", err)
		return err
	}

	if c.Commit {
		if _, err := client.Exec(nc.MethodCommit()); err != nil {
			return err
		}
	}

	return nil
}

func EditConfigCmd() *cobra.Command {
	edit := EditConfigCommand{}
	c := edit.SetFlags(&cobra.Command{
		Use:   "edit-config [running/candidate] [filename or -]",
		Short: "edit-config command.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return edit.EditConfig(args[0], args[1])
		},
	})

	return c
}
