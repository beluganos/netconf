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

package cfgvtylib

import (
	"fmt"
	api "netconf/app/cfg/api"

	"golang.org/x/net/context"
)

const CMD_LDP_ROUTER = "mpls ldp"

func SetMplsLdpCmd(negate bool, ipv uint, ifname string, args []string) []string {
	neg := NegateToStr(negate)

	cmds := []string{
		CMD_CONF_BEGIN,
		CMD_LDP_ROUTER,
	}

	cmd := fmt.Sprintf("%s%s", neg, joinArgs(args))
	switch ipv {
	case 0:
		cmds = append(cmds, cmd)

	case 4, 6:
		cmds = append(cmds, fmt.Sprintf("address-family ipv%d", ipv))
		if len(ifname) == 0 {
			cmds = append(cmds, cmd)
		} else {
			cmds = append(cmds,
				fmt.Sprintf("interface %s", ifname),
				cmd,
				CMD_EXIT,
			)
		}
		cmds = append(cmds, CMD_EXIT)
	}

	cmds = append(cmds, CMD_CONF_END)
	return cmds
}

func SetMplsLdpRun(negate bool, ipv uint, ifname string, args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeVtyExecuteRequest(SetMplsLdpCmd(negate, ipv, ifname, args))
	return client.Execute(context.Background(), req)
}
