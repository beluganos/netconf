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

func SetGlobalCmd(negate bool, args []string) []string {
	neg := NegateToStr(negate)
	return []string{
		CMD_CONF_BEGIN,
		fmt.Sprintf("%s%s", neg, joinArgs(args)),
		CMD_CONF_END,
	}
}

func SetGlobalRun(negate bool, args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeVtyExecuteRequest(SetGlobalCmd(negate, args))
	return client.Execute(context.Background(), req)
}

func SetGlobalRouterIdRun(negate bool, routerId string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeVtyExecuteRequest(SetGlobalCmd(negate, []string{"router-id", routerId}))
	return client.Execute(context.Background(), req)
}
