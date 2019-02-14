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
	"fmt"
	srlib "netconf/lib/sysrepo"

	log "github.com/sirupsen/logrus"
)

type SubscribeController struct {
	session *srlib.SrSession
	NoCopy  bool
	Debug   bool
}

func NewSubscribeController(session *srlib.SrSession, noCopy bool, debug bool) *SubscribeController {
	return &SubscribeController{
		session: session,
		NoCopy:  noCopy,
		Debug:   debug,
	}
}

func (s *SubscribeController) Start(name string, done chan struct{}) (*srlib.Subscriber, error) {
	return srlib.NewModuleChangeSubscriber(
		s.session,
		name,
		s,
		srlib.SR_SUBSCR_DEFAULT,
	)
}

func (s *SubscribeController) Notify(session *srlib.SrSession, module string, ev srlib.SrNotifEvent) error {
	log.Debugf("Notify module=%s ev=%s", module, ev)

	if s.Debug {
		for cv := range session.GetChanges(fmt.Sprintf("/%s:*", module)) {
			log.Debugf("%s", cv)
		}
	}

	if ev == srlib.SR_EV_APPLY && !s.NoCopy {
		s.CopyConfig(s.session, module)
	}

	return nil
}

func (s *SubscribeController) CopyConfig(session *srlib.SrSession, module string) {
	err := s.session.CopyConfig(module, srlib.SR_DS_RUNNING, srlib.SR_DS_STARTUP)
	if err != nil {
		log.Errorf("CopyConfig(%s) error. %s", module, err)
	} else {
		log.Infof("CopyConfig(%s) success.", module)
	}
}
