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

package cfgrpcapi

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func PrintReply(reply *ExecuteReply) {
	for i, r := range reply.Results {
		log.Infof("-- Reply#%d --", i)
		for _, s := range r.Strings() {
			log.Infof("%s", s)
		}
	}
}

type Command struct {
	host    string
	port    uint
	verbose bool
}

func (c *Command) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.host, "host", "H", "localhost", "Host Name")
	cmd.PersistentFlags().UintVarP(&c.port, "port", "P", LISTEN_PORT, "Host port")
	cmd.PersistentFlags().BoolVarP(&c.verbose, "verbose", "v", false, "Host port")

	return cmd
}

func (c *Command) Client() (RpcApiClient, *grpc.ClientConn, error) {
	return NewInsecureClient(c.host, c.port)
}

func (c *Command) Init() {
	if c.verbose {
		log.SetLevel(log.DebugLevel)
	}
}
