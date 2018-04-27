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
#include <sysrepo/values.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//
// sr_val_t
//
type SrVal struct {
	Data    string
	DefFlag bool
	Type    SrType
	Xpath   string
}

func NewSrVal(data string, defflag bool, srType SrType, xpath string) *SrVal {
	return &SrVal{
		Data:    data,
		DefFlag: defflag,
		Type:    srType,
		Xpath:   xpath,
	}
}

func NewSrValFromSr(v *C.sr_val_t) *SrVal {
	if v == nil {
		return nil
	}

	c_data := C.sr_val_to_str(v)
	defer C.free(unsafe.Pointer(c_data))

	return &SrVal{
		Data:    C.GoString(c_data),
		DefFlag: bool(v.dflt),
		Type:    SrType(v._type),
		Xpath:   C.GoString(v.xpath),
	}
}

func ParseSrVal(v interface{}, defflag bool, srType SrType, xpath string) *SrVal {
	data := func() string {
		switch srType {
		case SR_BOOL_T:
			return fmt.Sprintf("%t", v)

		case SR_DECIMAL64_T, SR_INT8_T, SR_INT16_T, SR_INT32_T, SR_INT64_T, SR_UINT8_T, SR_UINT16_T, SR_UINT32_T, SR_UINT64_T:
			return fmt.Sprintf("%d", v)

		case SR_ENUM_T, SR_IDENTITYREF_T, SR_STRING_T, SR_INSTANCEID_T:
			return fmt.Sprintf("%s", v)

		case SR_BINARY_T:
			return fmt.Sprintf("%x", v)

		case SR_TREE_ITERATOR_T, SR_LIST_T, SR_CONTAINER_T, SR_CONTAINER_PRESENCE_T, SR_LEAF_EMPTY_T:
			return fmt.Sprintf("%v", v)

		default:
			return fmt.Sprintf("%v", v)
		}
	}()

	return NewSrVal(data, defflag, srType, xpath)
}

func (v *SrVal) DataString() string {
	if v == nil {
		return "<nil>"
	}

	switch v.Type {
	case SR_ENUM_T, SR_IDENTITYREF_T, SR_STRING_T, SR_INSTANCEID_T:
		return fmt.Sprintf("'%s'", v.Data)
	case SR_BINARY_T:
		return fmt.Sprintf("[%s]", v.Data)
	default:
		return fmt.Sprintf("%s", v.Data)
	}
}

func (v *SrVal) String() string {
	if v == nil {
		return ""
	}

	dflt := func() string {
		if v.DefFlag {
			return "default"
		}
		return ""
	}()

	if data := v.DataString(); len(data) == 0 {
		return fmt.Sprintf("%s %s %s", v.Xpath, v.Type, dflt)
	} else {
		return fmt.Sprintf("%s = %s %s %s", v.Xpath, data, v.Type, dflt)
	}
}

func SrValGoString(v *C.sr_val_t) string {
	s, err := SrValGoStringSafe(v, "")
	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	return s
}

func SrValGoStringSafe(v *C.sr_val_t, defaultStr string) (string, error) {
	if v == nil {
		return defaultStr, nil
	}

	var mem *C.char = nil
	if rc := C.sr_print_val_mem(&mem, v); rc != C.SR_ERR_OK {
		return "", fmt.Errorf("sr_print_val_mem error. %d", rc)
	}

	defer C.free(unsafe.Pointer(mem))
	return C.GoString(mem), nil
}
