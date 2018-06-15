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
	"fmt"
	"net"
	"netconf/lib/lxd"

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
	dns     bool
	mngif   string
}

func (c *Command) SetFlags(cmd *cobra.Command) *cobra.Command {
	cmd.PersistentFlags().StringVarP(&c.host, "host", "H", "localhost", "Host Name")
	cmd.PersistentFlags().UintVarP(&c.port, "port", "P", LISTEN_PORT, "Host port")
	cmd.PersistentFlags().BoolVarP(&c.dns, "dns", "", false, "resolve host by dns.")
	cmd.PersistentFlags().StringVarP(&c.mngif, "mngif", "", "eth0", "Management interface on Container.")
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

	if !c.dns {
		ip, err := resolveName(c.host, c.mngif)
		if err != nil {
			log.Errorf("resolveName error. %s %s %s", c.host, c.mngif, err)
			return
		}

		log.Debugf("%s/%s resolved as %s", c.host, c.mngif, ip)
		c.host = ip.String()
	}
}

func resolveName(name string, ifname string) (net.IP, error) {
	client := lxdlib.NewClient("")
	if err := client.Connect(); err != nil {
		log.Errorf("Connect to lxd error.%s", err)
		return nil, err
	}

	state, _, err := lxdlib.GetState(client.Server, name)
	if err != nil {
		log.Errorf("GetState(%s) error. %s", name, err)
		return nil, err
	}

	nw, ok := state.Network[ifname]
	if !ok {
		log.Errorf("Network not found. name: %s, ifname: %s", name, ifname)
		return nil, fmt.Errorf("Network not found. name: %s, ifname: %s", name, ifname)
	}

	var ipv6 net.IP = nil

	for _, addr := range nw.Addresses {
		if addr.Scope == "global" {
			switch addr.Family {
			case "inet":
				if ip := net.ParseIP(addr.Address); ip.IsGlobalUnicast() {
					log.Debugf("state.NetworkAddr: IPv4 %s %s %s", name, ifname, ip)
					return ip, nil
				} else {
					log.Debugf("state.NetworkAddr: Not Global %s %s %v", name, ifname, ip)
				}
			case "inet6":
				if ip := net.ParseIP(addr.Address); ip.IsGlobalUnicast() {
					log.Debugf("state.NetworkAddr: IPv6 %s %s %s", name, ifname, ip)
					ipv6 = ip
				} else {
					log.Debugf("state.NetworkAddr: Not Global %s %s %s", name, ifname, ip)
				}
			default:
				log.Warnf("state.NetworkAddr: Unknown %s %s %v", name, ifname, addr)
			}
		} else {
			log.Debugf("state.NetworkAddr: Not global %s %s %v", name, ifname, addr)
		}
	}

	if ipv6 == nil {
		return nil, fmt.Errorf("state.NetworkAddr: global address not found. %s %s", name, ifname)
	}

	return ipv6, nil
}
