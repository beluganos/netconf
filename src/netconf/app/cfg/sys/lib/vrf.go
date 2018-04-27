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

package cfgsyslib

import (
	"fmt"
	api "netconf/app/cfg/api"

	"golang.org/x/net/context"
)

const (
	VRF_CONF_PATH = "/etc/vrf.conf"
)

func VrfBackupPath(path string) string {
	return fmt.Sprintf("%s.backup", path)
}

func DoVrfExec(cmd string, args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	params := []string{"-vrf", "-cmd", cmd}
	params = append(params, args...)

	shell := api.NewShell("cfgsysctl", params...)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func LoadVrfExec(client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("systemctl", "restart", "vrf")
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func BackupVrfExec(path, backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-f", path, backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func RollbackVrfExec(backup, path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-m", backup, path)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func CommitVrfExec(backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	if reply, err := LoadVrfExec(client); err != nil {
		return reply, err
	}

	shell := api.NewShell("rm", "-f", backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}
