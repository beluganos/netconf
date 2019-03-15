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

package ncm

import (
	"fmt"
	ncmcfg "netconf/app/ncm/cfg"
	ncmdbm "netconf/app/ncm/dbm"
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"

	log "github.com/sirupsen/logrus"
)

type NIChangeHandler interface {
	Begin(string, *openconfig.NetworkInstance) error
	Commit() error
	Rollback()
	SetOpt(string, interface{})
}

type NIChangeSet interface {
	Unmarshall(*srlib.SrChangeVal) error
	Walk(srlib.SrChangeOper, func(string, *openconfig.NetworkInstance) error) error
}

type niChangeFactory interface {
	NewHandler(srlib.SrNotifEvent, srlib.SrChangeOper) NIChangeHandler
	NewChangeSet() NIChangeSet
}

type NIChangeController struct {
	factory niChangeFactory
	session *srlib.SrSession
}

func NewNIChangeController(session *srlib.SrSession, factory niChangeFactory) *NIChangeController {
	return &NIChangeController{
		factory: factory,
		session: session,
	}
}

func (c *NIChangeController) Subscribe(flags ...srlib.SrSubscrFlag) (*srlib.Subscriber, error) {
	return srlib.NewModuleChangeSubscriber(
		c.session,
		openconfig.NETWORKINSTANCES_MODULE,
		c,
		srlib.SR_SUBSCR_DEFAULT,
	)
}

func (c *NIChangeController) Notify(session *srlib.SrSession, module string, ev srlib.SrNotifEvent) error {
	log.Debugf("NIChangeController module=%s ev=%s", module, ev)

	chgset, err := c.unmarshall(session, module)
	if err != nil {
		return err
	}

	if err := ncmdbm.Refresh(); err != nil {
		return err
	}

	err = c.callHandlers(chgset, ev,
		srlib.SR_OP_MODIFIED,
		srlib.SR_OP_DELETED,
		srlib.SR_OP_CREATED,
	)
	if err != nil {
		return err
	}

	if ev == srlib.SR_EV_APPLY && ncmcfg.GetConfig().Global.Persist {
		copyConfigAsync(c.session, module)
	}

	return nil
}

func (c *NIChangeController) unmarshall(session *srlib.SrSession, module string) (NIChangeSet, error) {

	chgset := c.factory.NewChangeSet()
	for cv := range session.GetChanges(fmt.Sprintf("/%s:*", module)) {
		log.Debugf("NetworkInstanceChange %s", cv)

		if err := chgset.Unmarshall(cv); err != nil {
			return nil, err
		}
	}

	return chgset, nil
}

func (c *NIChangeController) callHandlers(chgset NIChangeSet, ev srlib.SrNotifEvent, opers ...srlib.SrChangeOper) error {
	for _, oper := range opers {
		if err := c.callHandler(chgset, ev, oper); err != nil {
			return err
		}
	}

	return nil
}

func (c *NIChangeController) callHandler(chgset NIChangeSet, ev srlib.SrNotifEvent, oper srlib.SrChangeOper) error {
	h := c.factory.NewHandler(ev, oper)
	if h == nil {
		return nil
	}

	return chgset.Walk(oper, func(name string, ni *openconfig.NetworkInstance) error {
		log.Debugf("NIChangeController BEGIN(%s/%s). %s", ev, oper, ni)
		if err := h.Begin(name, ni); err != nil {
			log.Infof("NIChangeController ROLLBACK(%s/%s). %s", ev, oper, ni)
			h.Rollback()
			return err
		}

		log.Debugf("NIChangeController COMMIT(%s/%s).", ev, oper)
		if err := h.Commit(); err != nil {
			log.Errorf("NIChangeController COMMIT(%s/%s) error. %s %s", ev, oper, err, ni)
			return err
		}

		log.Infof("NIChangeController COMMIT(%s/%s) Success.", ev, oper)
		return nil
	})
}
