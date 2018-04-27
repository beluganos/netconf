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

package main

import (
	api "netconf/app/cfg/api"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func RegisterRpcApiServer(s *grpc.Server, srv *RpcApiServer) {
	api.RegisterRpcApiServer(s, srv)
}

type RpcApiServer struct {
}

func NewRpcApiServer() *RpcApiServer {
	return &RpcApiServer{}
}

func (s *RpcApiServer) Execute(ctxt context.Context, req *api.ExecuteRequest) (*api.ExecuteReply, error) {
	log.Debugf("Execute")

	results := []*api.Result{}
	for _, s := range req.Shells {
		output, err := s.ToNative().Exec()
		results = append(results, api.NewResult(output))

		if err != nil {
			log.Errorf("Execute: %s %s %s", s, err, string(output))
			return api.NewExecuteReply(results...), err
		}

		log.Debugf("Execute: %s %s", s, string(output))
	}

	return api.NewExecuteReply(results...), nil
}
