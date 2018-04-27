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

package srlib

/*
#include <stdio.h>
#include <sysrepo.h>
#include "helper.h"
*/
import "C"

type Subscriber struct {
	session *SrSession
	subscr  *C.sr_subscription_ctx_t
	opts    C.sr_subscr_options_t
}

func NewSubscriber(session *SrSession, flags ...SrSubscrFlag) *Subscriber {
	return &Subscriber{
		session: session,
		subscr:  nil,
		opts:    JoinSrSubscrFlags(flags...).C(),
	}
}

func (s *Subscriber) Stop() {
	if s.subscr != nil {
		C.sr_unsubscribe(s.session.session, s.subscr)
		s.subscr = nil
	}
}
