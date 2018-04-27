//-*- coding: utf-8 -*-

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

func SetIPCmd(negate bool, ipv string, args []string) []string {
	neg := NegateToStr(negate)
	return []string{
		CMD_CONF_BEGIN,
		fmt.Sprintf("%s%s %s", neg, ipv, joinArgs(args)),
		CMD_CONF_END,
	}
}

func SetIPv4(negate bool, args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeVtyExecuteRequest(SetIPCmd(negate, "ip", args))
	return client.Execute(context.Background(), req)
}

func SetIPv6(negate bool, args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeVtyExecuteRequest(SetIPCmd(negate, "ipv6", args))
	return client.Execute(context.Background(), req)
}
