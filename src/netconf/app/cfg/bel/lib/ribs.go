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

package cfgbellib

import (
	"fmt"
	api "netconf/app/cfg/api"

	"golang.org/x/net/context"
)

const RIBX_CONF_PATH = "/etc/beluganos/ribxd.conf"

func RibxBackupPath(path string) string {
	return fmt.Sprintf("%s.backup", path)
}

func SetRibsVrf(rt string, rd string, path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	params := []string{"ribs", "set", "vrf", "-t", rt, "-d", rd, "-f", path}
	shell := api.NewShell("cfgbelug", params...)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func BackupRibx(path string, backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-f", path, backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func RollbackRibx(backup string, path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-m", backup, path)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func LoadRibs(client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("systemctl", "restart", "ribs")
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)

}

func CommitRibx(backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	if r, err := LoadRibs(client); err != nil {
		return r, err
	}

	shell := api.NewShell("rm", "-f", backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}
