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

const (
	FRR_CONF_PATH = "/etc/frr/frr.conf"
	DAEMONS_PATH  = "/etc/frr/daemons"
)

func ConfigBackupPath(path string) string {
	return fmt.Sprintf("%s.backup", path)
}

func ShowConfigCmd(args []string) []string {
	cmd := fmt.Sprintf("show %s", joinArgs(args))
	return []string{cmd}
}

func ShowConfigRun(args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeVtyExecuteRequest(ShowConfigCmd(args))
	return client.Execute(context.Background(), req)
}

func ShowBackupConfigRun(path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeExecuteRequest("cat", path)
	return client.Execute(context.Background(), req)
}

func SaveConfigCmd() []string {
	return []string{
		"write file",
	}
}

func SaveConfigRun(client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeVtyExecuteRequest(SaveConfigCmd())
	return client.Execute(context.Background(), req)
}

func RollbackConfigRun(backup, path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeExecuteRequest("cp", "-f", backup, path)
	return client.Execute(context.Background(), req)
}

func BackupConfigRun(path, backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeExecuteRequest("cp", "-f", path, backup)
	return client.Execute(context.Background(), req)
}

func DaemonConfigRun(client api.RpcApiClient) (*api.ExecuteReply, error) {
	req := makeExecuteRequest("cat", DAEMONS_PATH)
	return client.Execute(context.Background(), req)
}

func SetDaemonRun(args []string, cmd string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	params := []string{"-cmd", cmd}
	params = append(params, args...)
	req := makeExecuteRequest("cfgfrr", params...)
	return client.Execute(context.Background(), req)
}
