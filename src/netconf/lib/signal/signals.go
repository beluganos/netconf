// -*- coding: utf-8 -*

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

package ncsignal

import (
	"os"
	"os/signal"
)

type SignalFunc func(os.Signal)

type SignalServer map[os.Signal]SignalFunc

func NewServer() SignalServer {
	return SignalServer{}
}

func (s SignalServer) Register(sig os.Signal, f SignalFunc) SignalServer {
	s[sig] = f
	return s
}

func (s SignalServer) Registers(m map[os.Signal]SignalFunc) SignalServer {
	for k, v := range m {
		s[k] = v
	}
	return s
}

func (s SignalServer) signals() []os.Signal {
	sigs := []os.Signal{}
	for k, _ := range s {
		sigs = append(sigs, k)
	}
	return sigs
}

func (s SignalServer) Start(done <-chan struct{}) {
	go s.Serve(done)
}

func (s SignalServer) Serve(done <-chan struct{}) {
	ch := make(chan os.Signal)
	defer close(ch)

	signal.Notify(ch, s.signals()...)

	for {
		select {
		case sig := <-ch:
			if f, ok := s[sig]; ok && f != nil {
				f(sig)
			}

		case <-done:
			break
		}
	}
}
