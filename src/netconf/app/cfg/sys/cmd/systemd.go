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

type SystemdCommand struct {
	api.Command
}

func (c *SystemdCommand) Systemctl(args []string) error {
	c.Command.Init()

	client, conn, err := c.Client()
	if err != nil {
		return err
	}
	defer conn.Close()

	res, err := lib.SystemctlRun(args, client)
	if err != nil {
		return err
	}

	api.PrintReply(res)
	return nil
}

func SystemdCmd() *cobra.Command {
	cfg := SystemdCommand{}

	c := &cobra.Command{
		Use:   "systemctl [option...]",
		Short: "Systemd configuration commands.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cfg.Systemctl(args)
		},
	}

	return cfg.SetFlags(c)
}
