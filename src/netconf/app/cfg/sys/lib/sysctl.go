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
	"time"

	"golang.org/x/net/context"
)

const (
	SYSCTL_CONF_PATH = "/etc/sysctl.d/30-beluganos.conf"
	SYSCTL_WAIT_SEC  = 1 * time.Second
)

func SysctlBackupPath(path string) string {
	return fmt.Sprintf("%s.backup", path)
}

func DoSysctlRun(cmd string, args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	params := []string{"-cmd", cmd}
	params = append(params, args...)

	shell := api.NewShell("cfgsysctl", params...)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func LoadSysctlRun(path string, wait time.Duration, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("sysctl", "-p", path)
	req := api.NewExecuteRequest(shell)
	r, err := client.Execute(context.Background(), req)
	if err == nil {
		time.Sleep(wait)
	}
	return r, err
}

func BackupSysctlRun(path, backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-f", path, backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func RollbackSysctlRun(backup, path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-m", backup, path)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func CommitSysctlRun(path, backup string, wait time.Duration, client api.RpcApiClient) (*api.ExecuteReply, error) {
	if reply, err := LoadSysctlRun(path, wait, client); err != nil {
		return reply, err
	}

	shell := api.NewShell("rm", "-f", backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}
