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
// ModuleChangeHandler
//
type ModuleChangeHandler interface {
	Notify(session *SrSession, module string, ev SrNotifEvent) error
}

var moduleChangesHandlers = map[string]ModuleChangeHandler{}

func registerModuleChanges(name string, handler ModuleChangeHandler) error {
	if _, ok := moduleChangesHandlers[name]; ok {
		return fmt.Errorf("module already exists. %s", name)
	}

	moduleChangesHandlers[name] = handler
	return nil
}

func unregisterModuleChanges(name string) {
	delete(moduleChangesHandlers, name)
}

func getModuleChange(name string) (ModuleChangeHandler, error) {
	if h, ok := moduleChangesHandlers[name]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("module not found. %s", name)
}

//
// Subscriber(ModuleChange)
//
func NewModuleChangeSubscriber(session *SrSession, module string, handler ModuleChangeHandler, flags ...SrSubscrFlag) (*Subscriber, error) {
	if err := registerModuleChanges(module, handler); err != nil {
		return nil, err
	}

	c_module := C.CString(module)
	defer C.free(unsafe.Pointer(c_module))

	subscr := NewSubscriber(session, flags...)
	ret := C.sr_module_change_subscribe(
		session.session,
		c_module,
		C.sr_module_change_cb(C.module_change_cb),
		nil,
		0,
		subscr.opts,
		&subscr.subscr,
	)
	if ret != C.SR_ERR_OK {
		unregisterModuleChanges(module)
		return nil, fmt.Errorf("sr_module_change_subscribe error. %d", ret)
	}

	return subscr, nil
}

//export Go_module_change_cb
func Go_module_change_cb(c_session *C.sr_session_ctx_t, c_module *C.char, c_event C.sr_notif_event_t, key unsafe.Pointer) C.int {
	session := NewSrSession(c_session)
	module := C.GoString(c_module)
	event := SrNotifEvent(c_event)

	handler, err := getModuleChange(module)
	if err != nil {
		log.Errorf("Go_module_change_cb error. %s", err)
		session.SetError(err, module)
		return C.SR_ERR_INTERNAL
	}

	if err := handler.Notify(session, module, event); err != nil {
		log.Errorf("handler.Notify error. %s", err)
		session.SetError(err, module)
		if event == SR_EV_VERIFY {
			return C.SR_ERR_VALIDATION_FAILED
		}
		return C.SR_ERR_INTERNAL
	}

	return C.SR_ERR_OK
}
