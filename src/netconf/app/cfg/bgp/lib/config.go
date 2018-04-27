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

package cfgbgplib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	api "netconf/app/cfg/api"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	GOBGP_CONF_PATH     = "/etc/frr/gobgpd.toml"
	GOBGP_WAIT_SEC      = 1 * time.Second
	GOBGP_PROCESS_NAME  = "gobgpd"
	GOBGP_RELOAD_SIGNAL = "-HUP"
)

func GobgpBackupPath(path string) string {
	return fmt.Sprintf("%s.backup", path)
}

func readConfig(path string) ([]byte, error) {
	var r io.Reader = nil
	if len(path) == 0 || path == "-" {
		r = os.Stdin

		log.Debugf("readConfig(stdin)")

	} else {
		f, err := os.Open(path)
		if err != nil {
			log.Errorf("readConfig error. %s", err)
			return nil, err
		}
		defer f.Close()
		r = f

		log.Debugf("readConfig(%s)", path)
	}

	return ioutil.ReadAll(r)
}

func readConfigs(paths []string) ([]byte, error) {
	var buf bytes.Buffer

	for _, path := range paths {
		config, err := readConfig(path)
		if err != nil {
			return nil, err
		}
		buf.Write(config)
		buf.WriteByte('\n')

		log.Debugf("readConfig. %s", string(config))
	}

	return buf.Bytes(), nil
}

func DoGobgpRun(cmd string, path string, args []string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	config, err := readConfigs(args)
	if err != nil {
		return nil, err
	}

	shell := api.NewShellIn("cfgbgp", config, "-c", path, "-cmd", cmd)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func LoadGobgpRun(wait time.Duration, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("pkill", GOBGP_RELOAD_SIGNAL, GOBGP_PROCESS_NAME)
	req := api.NewExecuteRequest(shell)
	r, err := client.Execute(context.Background(), req)
	if err == nil {
		time.Sleep(wait)
	}
	return r, err
}

func BackupGobgpRun(path, backup string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-f", path, backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func RollbackGobgpRun(backup, path string, client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("cfgcp", "-m", backup, path)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func CommitGobgpRun(backup string, wait time.Duration, client api.RpcApiClient) (*api.ExecuteReply, error) {
	if reply, err := LoadGobgpRun(wait, client); err != nil {
		return reply, err
	}

	shell := api.NewShell("rm", "-f", backup)
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}

func RestartGobgpRun(client api.RpcApiClient) (*api.ExecuteReply, error) {
	shell := api.NewShell("systemctl", "restart", "gobgpd")
	req := api.NewExecuteRequest(shell)
	return client.Execute(context.Background(), req)
}
