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
	srlib "netconf/lib/sysrepo"
)

//
// Factory - Handler
//
type NIChangeFactoryFunc func(srlib.SrNotifEvent, srlib.SrChangeOper) NIChangeHandler

type NIChangeFactory struct {
	DryRun   bool
	Mtu      uint16
	handlers map[srlib.SrNotifEvent]map[srlib.SrChangeOper]NIChangeFactoryFunc
}

func NewNIChangeFactory() *NIChangeFactory {
	return &NIChangeFactory{
		DryRun: false,
		Mtu:    NIConterinerDefaultMTU,
		handlers: map[srlib.SrNotifEvent]map[srlib.SrChangeOper]NIChangeFactoryFunc{
			srlib.SR_EV_VERIFY: {
				srlib.SR_OP_CREATED:  NewNICreateVerifyHandler,
				srlib.SR_OP_MODIFIED: NewNIModifyVerifyHandler,
				srlib.SR_OP_DELETED:  NewNIDeleteVerifyHandler,
				srlib.SR_OP_MOVED:    NewNIAnyHandler,
			},
			srlib.SR_EV_APPLY: {
				srlib.SR_OP_CREATED:  NewNICreateApplyHandler,
				srlib.SR_OP_MODIFIED: NewNIModifyApplyHandler,
				srlib.SR_OP_DELETED:  NewNIDeleteApplyHandler,
				srlib.SR_OP_MOVED:    NewNIAnyHandler,
			},
			srlib.SR_EV_ABORT: {
				srlib.SR_OP_CREATED:  NewNIAnyHandler,
				srlib.SR_OP_MODIFIED: NewNIAnyHandler,
				srlib.SR_OP_DELETED:  NewNIAnyHandler,
				srlib.SR_OP_MOVED:    NewNIAnyHandler,
			},
			srlib.SR_EV_ENABLED: {
				srlib.SR_OP_CREATED:  NewNIAnyHandler,
				srlib.SR_OP_MODIFIED: NewNIAnyHandler,
				srlib.SR_OP_DELETED:  NewNIAnyHandler,
				srlib.SR_OP_MOVED:    NewNIAnyHandler,
			},
		},
	}
}

func (n *NIChangeFactory) NewHandler(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) NIChangeHandler {
	if opers, ok := n.handlers[ev]; ok {
		if f, ok := opers[oper]; ok && f != nil {
			h := f(ev, oper)
			h.SetOpt("dryrun", n.DryRun)
			h.SetOpt("mtu", n.Mtu)
			return h
		}
	}
	return nil
}

func (f *NIChangeFactory) NewChangeSet() NIChangeSet {
	return NewNetworkInstancesSet()
}
