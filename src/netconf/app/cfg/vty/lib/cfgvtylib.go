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
	api "netconf/app/cfg/api"
	"strings"
)

const (
	CONFIG_STARTUP = "startup-config"
	CONFIG_RUNNING = "running-config"

	CMD_CONF_BEGIN = "configure terminal"
	CMD_CONF_END   = "end"
	CMD_EXIT       = "exit"
)

func NegateToStr(negate bool) string {
	if negate {
		return "no "
	}
	return ""
}

func joinArgs(args []string) string {
	return strings.Join(args, " ")
}

func makeVtyArgs(cmds []string) []string {
	args := make([]string, len(cmds)*2)
	for i, cmd := range cmds {
		args[i*2] = "-c"
		args[i*2+1] = cmd
	}
	return args
}

func makeVtyExecuteRequest(args []string) *api.ExecuteRequest {
	return makeExecuteRequest("vtysh", makeVtyArgs(args)...)
}

func makeExecuteRequest(cmd string, args ...string) *api.ExecuteRequest {
	shell := api.NewShell(cmd, args...)
	return api.NewExecuteRequest(shell)

}
