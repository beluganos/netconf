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
import (
	"fmt"
	"unsafe"

	log "github.com/sirupsen/logrus"
)

//
// SubtreeChangeHandler
//
type SubtreeChangeHandler interface {
	Notify(session *SrSession, xpath string, ev SrNotifEvent) error
}

var subtreeChangeHandlers = map[string]SubtreeChangeHandler{}

func registerSubtreeChanges(xpath string, handler SubtreeChangeHandler) error {
	if _, ok := subtreeChangeHandlers[xpath]; ok {
		return fmt.Errorf("xpath already exists. %s", xpath)
	}

	subtreeChangeHandlers[xpath] = handler
	return nil
}

func unregisterSubtreeChanges(xpath string) {
	delete(subtreeChangeHandlers, xpath)
}

func getSubtreeChanges(xpath string) (SubtreeChangeHandler, error) {
	if h, ok := subtreeChangeHandlers[xpath]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("xpath not found. %s", xpath)
}

//
// Subscriber(SubtreeChange)
//
func NewSubtreeChangeSubscriber(session *SrSession, xpath string, handler SubtreeChangeHandler) (*Subscriber, error) {
	if err := registerSubtreeChanges(xpath, handler); err != nil {
		return nil, err
	}

	c_xpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(c_xpath))

	subscr := NewSubscriber(session)
	ret := C.sr_module_change_subscribe(
		session.session,
		c_xpath,
		C.sr_subtree_change_cb(C.subtree_change_cb),
		nil,
		0,
		subscr.opts,
		&subscr.subscr,
	)
	if ret != C.SR_ERR_OK {
		unregisterSubtreeChanges(xpath)
		return nil, fmt.Errorf("sr_module_change_subscribe error. %d", ret)
	}

	return subscr, nil
}

//export Go_subtree_change_cb
func Go_subtree_change_cb(c_session *C.sr_session_ctx_t, c_xpath *C.char, c_event C.sr_notif_event_t, key unsafe.Pointer) C.int {
	session := NewSrSession(c_session)
	xpath := C.GoString(c_xpath)
	event := SrNotifEvent(c_event)

	handler, err := getSubtreeChanges(xpath)
	if err != err {
		log.Errorf("Go_subtree_change_cb error. %s", err)
		session.SetError(err, xpath)
		return C.SR_ERR_INTERNAL
	}

	if err := handler.Notify(session, xpath, event); err != nil {
		log.Errorf("handler.Notify error. %s", err)
		session.SetError(err, xpath)
		if event == SR_EV_VERIFY {
			return C.SR_ERR_VALIDATION_FAILED
		}
		return C.SR_ERR_INTERNAL
	}

	return C.SR_ERR_OK
}
