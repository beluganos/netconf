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

package cfgrpcapi

import (
	"bytes"
	"fmt"
	"netconf/lib"
	"strings"

	"google.golang.org/grpc"
)

const LISTEN_PORT = 50081

func NewShell(cmd string, args ...string) *Shell {
	return NewShellIn(cmd, []byte{}, args...)
}

func NewShellIn(cmd string, in []byte, args ...string) *Shell {
	return &Shell{
		Cmd:  cmd,
		Args: args,
		In:   in,
	}
}

func (s *Shell) ToNative() *nclib.Shell {
	if s.In == nil || len(s.In) == 0 {
		return nclib.NewShell(s.Cmd, s.Args...)
	} else {
		r := bytes.NewReader(s.In)
		return nclib.NewShellIn(s.Cmd, r, s.Args...)
	}
}

//
// Result
//
func (r *Result) Strings() []string {
	return strings.Split(string(r.Output), "\n")
}

func NewResult(output []byte) *Result {
	return &Result{
		Output: output,
	}
}

func NewExecuteRequest(shells ...*Shell) *ExecuteRequest {
	return &ExecuteRequest{
		Shells: shells,
	}
}

func NewExecuteReply(results ...*Result) *ExecuteReply {
	return &ExecuteReply{
		Results: results,
	}
}

//
// ApiClient
//
func NewInsecureClient(host string, port uint, opts ...grpc.DialOption) (RpcApiClient, *grpc.ClientConn, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return nil, nil, err
	}

	return NewRpcApiClient(conn), conn, nil
}
