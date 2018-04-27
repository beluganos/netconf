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

const NETPLAN_CONF_PATH = "/etc/netplan/02-beluganos.yaml"
const NETPLAN_WAIT_SEC = 1 * time.Second
const NETPLAN_CMD = "netplan+"

func NetworkBackupPath(path string) string {
	return fmt.Sprintf("%s.backup", path)
}

func makeCfgnetParams(cmd string, device string, vid uint, mtu uint, addrs []string) []string {
	params := []string{
		"-device", device,
		"-cmd", cmd,
		"-vid", fmt.Sprintf("%d", vid),
		"-mtu", fmt.Sprintf("%d", mtu),
	}
	for _, addr := range addrs {
		params = append(params, "-a", addr)
	}
	return params
}

func DoNetworkRun(cmd string, device string, vid uint, mtu uint, addrs []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	params := makeCfgnetParams(cmd, device, vid, mtu, addrs)
	shell := api.NewShell("cfgnet", params...)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func SetNetworkRun(device string, vid uint, mtu uint, addrs []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	return DoNetworkRun("set", device, vid, mtu, addrs, client)
}

func DelNetworkRun(device string, vid uint, mtu uint, addrs []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	return DoNetworkRun("del", device, vid, mtu, addrs, client)
}

func LoadNetworkRun(wait time.Duration, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell(NETPLAN_CMD, "apply")
	req := api.NewExecuteRequest(shell)
	r, err := client.Execute(context.Background(), req)
	if err == nil {
		time.Sleep(wait)
	}
	return r, err
}

func BackupNetwotkRun(path, backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-f", path, backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func RollbackNetworkRun(backup, path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-m", backup, path)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func CommitNetworkRun(backup string, wait time.Duration, client api.RpcApiClient) (*api.ExecuteReply, error) {
	if reply, err := LoadNetworkRun(wait, client); err != nil {
		return reply, err
	}

	shell := api.NewShell("rm", "-f", backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}
