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

type CopyConfigCommand struct {
	Command
}

func (c *CopyConfigCommand) CopyConfig(source string, target string) error {
	c.Command.Init()

	client, err := c.Client()
	if err != nil {
		log.Errorf("Client initialize error. %s", err)
		return err
	}
	defer client.Close()

	if _, err := client.Exec(nc.MethodCopyConfig(target, source)); err != nil {
		log.Errorf("Exec error. %s", err)
		return err
	}

	return nil
}

func CopyConfigCmd() *cobra.Command {
	copy := CopyConfigCommand{}
	c := copy.SetFlags(&cobra.Command{
		Use:   "copy-config [source] [target]",
		Short: "copy-config command.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return copy.CopyConfig(args[0], args[1])
		},
	})

	return c
}
