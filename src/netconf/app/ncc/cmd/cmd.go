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
	"os"

	nc "github.com/hiepon/go-netconf/netconf"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Command struct {
	Host     string
	Username string
	Password string
	Verbose  bool
}

func (c *Command) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.Host, "host", "H", "localhost:830", "Host name.")
	cmd.PersistentFlags().StringVarP(&c.Username, "user", "U", os.Getenv("USER"), "Username.")
	cmd.PersistentFlags().StringVarP(&c.Password, "passwd", "P", "", "Password.")
	cmd.PersistentFlags().BoolVarP(&c.Verbose, "verbose", "v", false, "Show detail messages.")
	return cmd
}

func (c *Command) Init() {
	if c.Verbose {
		log.SetLevel(log.DebugLevel)
	}
	nc.SetLog(log.StandardLogger())
}

func (c *Command) Client() (session *nc.Session, err error) {
	config := nc.SSHConfigPassword(c.Username, c.Password)
	session, err = nc.DialSSH(c.Host, config)
	if session != nil {
		log.Debugf("SESSION: %d", session.SessionID)
		for index, capa := range session.ServerCapabilities {
			log.Debugf("CAPA[%d]: %s", index, capa)
		}
	}
	return
}
