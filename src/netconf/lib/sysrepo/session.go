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
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//
// sr_session_ctx
//
type SrSession struct {
	session *C.sr_session_ctx_t
}

func NewSrSession(session *C.sr_session_ctx_t) *SrSession {
	return &SrSession{
		session: session,
	}
}

func (s *SrSession) Start(conn *SrConnection, ds SrDataStore) error {
	rc := C.sr_session_start(conn.conn, ds.C(), C.SR_SESS_DEFAULT, &s.session)
	if rc != C.SR_ERR_OK {
		return fmt.Errorf("sr_session_start error. %d", rc)
	}

	return nil
}

func (s *SrSession) Stop() {
	if s.session != nil {
		C.sr_session_stop(s.session)
		s.session = nil
	}
}

func (s *SrSession) Refresh() error {
	return RefreshSrSession(s)
}

func RefreshSrSession(s *SrSession) error {
	rc := C.sr_session_refresh(s.session)
	if rc != C.SR_ERR_OK {
		return fmt.Errorf("sr_session_refresh error. %d", rc)
	}
	return nil
}

func (s *SrSession) SetError(err error, xpath string) {
	SetSrSessionError(s, err, xpath)
}

func SetSrSessionError(s *SrSession, err error, xpath string) {
	c_msg := C.CString(err.Error())
	defer C.free(unsafe.Pointer(c_msg))

	c_xpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(c_xpath))

	C.sr_set_error(s.session, c_msg, c_xpath)
}
