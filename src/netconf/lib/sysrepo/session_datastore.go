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

func (s *SrSession) Commit() error {
	return Commit(s)
}

func Commit(session *SrSession) error {
	ret := C.sr_commit(session.session)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_commit error. %d", ret)
	}

	return nil
}

func (s *SrSession) CopyConfig(moduleName string, src SrDataStore, dst SrDataStore) error {
	return CopyConfig(s, moduleName, src, dst)
}

func CopyConfig(session *SrSession, moduleName string, src SrDataStore, dst SrDataStore) error {
	c_module := C.CString(moduleName)
	defer C.free(unsafe.Pointer(c_module))

	ret := C.sr_copy_config(
		session.session,
		c_module,
		C.sr_datastore_t(src),
		C.sr_datastore_t(dst),
	)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_copy_config error. %d", ret)
	}

	return nil
}

func (s *SrSession) SetItemNull(xpath string, opts SrEditOptions) error {
	return SetItemNull(s, xpath, opts)
}

func SetItemNull(session *SrSession, xpath string, opts SrEditOptions) error {
	c_xpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(c_xpath))

	ret := C.sr_set_item_str(
		session.session,
		c_xpath,
		nil,
		C.sr_edit_options_t(opts),
	)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_set_item_str error. %s %s %d", xpath, opts, ret)
	}

	return nil
}

func (s *SrSession) SetItemStr(xpath string, value string, opts SrEditOptions) error {
	return SetItemStr(s, xpath, value, opts)
}

func SetItemStr(session *SrSession, xpath string, value string, opts SrEditOptions) error {
	c_xpath := C.CString(xpath)
	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_xpath))
	defer C.free(unsafe.Pointer(c_value))

	ret := C.sr_set_item_str(
		session.session,
		c_xpath,
		c_value,
		C.sr_edit_options_t(opts),
	)
	if ret != C.SR_ERR_OK {
		return fmt.Errorf("sr_set_item_str error. %s -> %s %s %d", xpath, value, opts, ret)
	}

	return nil
}

func (s *SrSession) SetItem(val *SrVal, opts SrEditOptions) error {
	return SetItem(s, val, opts)
}

func SetItem(session *SrSession, val *SrVal, opts SrEditOptions) error {
	return SetItemStr(session, val.Xpath, val.Data, opts)
}

func (s *SrSession) GetItems(xpath string) <-chan *SrChangeVal {
	ch := make(chan *SrChangeVal)
	go func() {
		defer close(ch)
		GetItems(s, xpath, func(cv *SrChangeVal) {
			ch <- cv
		})
	}()

	return ch
}

func GetItems(session *SrSession, xpath string, f func(*SrChangeVal)) error {
	c_xpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(c_xpath))

	var iter *C.sr_val_iter_t = nil
	if rc := C.sr_get_items_iter(session.session, c_xpath, &iter); rc != C.SR_ERR_OK {
		return fmt.Errorf("sr_get_items_iter error. %d", rc)
	}
	defer C.sr_free_val_iter(iter)

	for {
		var value *C.sr_val_t = nil
		rc := C.sr_get_item_next(session.session, iter, &value)
		if rc == C.SR_ERR_NOT_FOUND {
			break
		}

		if rc != C.SR_ERR_OK {
			return fmt.Errorf("sr_get_item_next error. %d", rc)
		}

		v := SrChangeVal{
			Oper:   SR_OP_CREATED,
			OldVal: NewSrValFromSr(nil),
			NewVal: NewSrValFromSr(value),
		}
		C.sr_free_val(value)

		f(&v)
	}

	return nil
}

func (s *SrSession) GetChanges(xpath string) <-chan *SrChangeVal {
	ch := make(chan *SrChangeVal)
	go func() {
		defer close(ch)
		GetChanges(s, xpath, func(cv *SrChangeVal) {
			ch <- cv
		})
	}()

	return ch
}

func GetChanges(session *SrSession, xpath string, f func(*SrChangeVal)) error {
	c_xpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(c_xpath))

	var iter *C.sr_change_iter_t = nil
	if rc := C.sr_get_changes_iter(session.session, c_xpath, &iter); rc != C.SR_ERR_OK {
		return fmt.Errorf("sr_get_changes_iter error. %d", rc)
	}
	defer C.sr_free_change_iter(iter)

	for {
		var oper C.sr_change_oper_t
		var old_val *C.sr_val_t = nil
		var new_val *C.sr_val_t = nil

		rc := C.sr_get_change_next(session.session, iter, &oper, &old_val, &new_val)
		if rc != C.SR_ERR_OK {
			break
		}

		v := SrChangeVal{
			Oper:   SrChangeOper(oper),
			OldVal: NewSrValFromSr(old_val),
			NewVal: NewSrValFromSr(new_val),
		}
		f(&v)
	}

	return nil
}
