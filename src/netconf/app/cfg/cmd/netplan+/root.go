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

package main

import (
	"os"

	"github.com/spf13/cobra"
)

func rootCmd(name string) *cobra.Command {
	var showCompletion bool
	c := &cobra.Command{
		Use:   name,
		Short: "extended netplan command.",
		Run: func(cmd *cobra.Command, args []string) {
			if showCompletion {
				cmd.GenBashCompletion(os.Stdout)
			} else {
				cmd.Usage()
			}
		},
	}

	c.PersistentFlags().BoolVar(&showCompletion, "show-completion", false, "Show bash-comnpletion")

	init := NpCommand{}
	c.AddCommand(init.SetInitFlags(
		&cobra.Command{
			Use:   "init",
			Short: "apply settings(not supported by netplan.)",
			Args:  cobra.NoArgs,
			RunE: func(md *cobra.Command, args []string) error {
				return init.init().Init()
			},
		},
	))

	apply := NpCommand{}
	c.AddCommand(apply.SetInitFlags(
		&cobra.Command{
			Use:   "apply",
			Short: "run netplan apply and init",
			Args:  cobra.NoArgs,
			RunE: func(md *cobra.Command, args []string) error {
				return apply.init().Apply()
			},
		},
	))

	gen := NpCommand{}
	c.AddCommand(gen.SetFlags(
		&cobra.Command{
			Use:   "generate",
			Short: "run netplan generate",
			Args:  cobra.NoArgs,
			RunE: func(md *cobra.Command, args []string) error {
				return gen.init().Run("generate")
			},
		},
	))

	ifupdown := NpCommand{}
	c.AddCommand(ifupdown.SetFlags(
		&cobra.Command{
			Use:   "ifupdown-migrate",
			Short: "run netplan ifupdown-migrate",
			Args:  cobra.NoArgs,
			RunE: func(md *cobra.Command, args []string) error {
				return ifupdown.init().Run("ifupdown-migrate")
			},
		},
	))

	return c
}
