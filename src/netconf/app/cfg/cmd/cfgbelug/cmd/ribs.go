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

package cfgbelugcmd

import (
	"github.com/spf13/cobra"
)

const RIBS_CONFIG_FILE = "/etc/beluganos/ribxd.conf"

type RibsCommand struct {
	Command
	path string
	rt   string
	rd   string
}

func (c *RibsCommand) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.path, "path", "f", RIBS_CONFIG_FILE, "config filename")
	cmd.PersistentFlags().StringVarP(&c.rt, "RT", "t", "", "route target")
	cmd.PersistentFlags().StringVarP(&c.rd, "RD", "d", "", "route distinguisher")
	return c.Command.SetFlags(cmd)
}

func (c *RibsCommand) SetVrf() error {
	cfg := RibsConfig{}
	if err := ReadConfig(c.path, &cfg); err != nil {
		return err
	}

	if len(c.rt) != 0 {
		if err := cfg.SetRibsVrf("rt", c.rt); err != nil {
			return err
		}
	}

	if len(c.rd) != 0 {
		if err := cfg.SetRibsVrf("rd", c.rd); err != nil {
			return err
		}
	}

	return WriteConfig(c.path, &cfg)
}

func RibsCmd() *cobra.Command {

	c := &cobra.Command{
		Use:   "ribs",
		Short: "Ribs configuration command.",
	}

	set := &cobra.Command{
		Use:   "set",
		Short: "Ribs set configuration.",
	}
	c.AddCommand(set)

	vrf := RibsCommand{}
	set.AddCommand(vrf.SetFlags(
		&cobra.Command{
			Use:   "vrf",
			Short: "Set Vrf config value.",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				return vrf.SetVrf()
			},
		},
	))

	return c
}
