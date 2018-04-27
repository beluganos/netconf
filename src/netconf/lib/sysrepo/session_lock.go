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

func (s *SrSession) LockDatastore() error {
	return LockDatastore(s)
}

func LockDatastore(session *SrSession) error {
	ret := C.sr_lock_datastore(session.session)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_lock_datastore error. %d", ret)
	}
	return nil
}

func (s *SrSession) UnlockDatastore() error {
	return UnlockDatastore(s)
}

func UnlockDatastore(session *SrSession) error {
	ret := C.sr_unlock_datastore(session.session)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_lock_datastore error. %d", ret)
	}
	return nil
}

func (s *SrSession) LockModule(moduleName string) error {
	return LockModule(s, moduleName)
}

func LockModule(session *SrSession, moduleName string) error {
	c_module := C.CString(moduleName)
	defer C.free(unsafe.Pointer(c_module))

	ret := C.sr_lock_module(session.session, c_module)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_lock_module error. %s %d", moduleName, ret)
	}
	return nil
}

func (s *SrSession) UnlockModule(moduleName string) error {
	return UnlockModule(s, moduleName)
}

func UnlockModule(session *SrSession, moduleName string) error {
	c_module := C.CString(moduleName)
	defer C.free(unsafe.Pointer(c_module))

	ret := C.sr_unlock_module(session.session, c_module)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_unlock_module error. %s %d", moduleName, ret)
	}
	return nil
}
