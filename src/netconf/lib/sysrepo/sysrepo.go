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
#cgo LDFLAGS: -lsysrepo
#include <stdio.h>
#include <sysrepo.h>
*/
import "C"

import (
	"fmt"
	"strings"
)

//
// sr_type_t
//
type SrType int

const (
	SR_UNKNOWN_T            SrType = C.SR_UNKNOWN_T
	SR_TREE_ITERATOR_T      SrType = C.SR_TREE_ITERATOR_T
	SR_LIST_T               SrType = C.SR_LIST_T
	SR_CONTAINER_T          SrType = C.SR_CONTAINER_T
	SR_CONTAINER_PRESENCE_T SrType = C.SR_CONTAINER_PRESENCE_T
	SR_LEAF_EMPTY_T         SrType = C.SR_LEAF_EMPTY_T
	SR_BINARY_T             SrType = C.SR_BINARY_T
	SR_BITS_T               SrType = C.SR_BITS_T
	SR_BOOL_T               SrType = C.SR_BOOL_T
	SR_DECIMAL64_T          SrType = C.SR_DECIMAL64_T
	SR_ENUM_T               SrType = C.SR_ENUM_T
	SR_IDENTITYREF_T        SrType = C.SR_IDENTITYREF_T
	SR_INSTANCEID_T         SrType = C.SR_INSTANCEID_T
	SR_INT8_T               SrType = C.SR_INT8_T
	SR_INT16_T              SrType = C.SR_INT16_T
	SR_INT32_T              SrType = C.SR_INT32_T
	SR_INT64_T              SrType = C.SR_INT64_T
	SR_STRING_T             SrType = C.SR_STRING_T
	SR_UINT8_T              SrType = C.SR_UINT8_T
	SR_UINT16_T             SrType = C.SR_UINT16_T
	SR_UINT32_T             SrType = C.SR_UINT32_T
	SR_UINT64_T             SrType = C.SR_UINT64_T
	SR_ANYXML_T             SrType = C.SR_ANYXML_T
	SR_ANYDATA_T            SrType = C.SR_ANYDATA_T
)

var srTypeValues = map[string]SrType{
	"SR_UNKNOWN_T":            SR_UNKNOWN_T,
	"SR_TREE_ITERATOR_T":      SR_TREE_ITERATOR_T,
	"SR_LIST_T":               SR_LIST_T,
	"SR_CONTAINER_T":          SR_CONTAINER_T,
	"SR_CONTAINER_PRESENCE_T": SR_CONTAINER_PRESENCE_T,
	"SR_LEAF_EMPTY_T":         SR_LEAF_EMPTY_T,
	"SR_BINARY_T":             SR_BINARY_T,
	"SR_BITS_T":               SR_BITS_T,
	"SR_BOOL_T":               SR_BOOL_T,
	"SR_DECIMAL64_T":          SR_DECIMAL64_T,
	"SR_ENUM_T":               SR_ENUM_T,
	"SR_IDENTITYREF_T":        SR_IDENTITYREF_T,
	"SR_INSTANCEID_T":         SR_INSTANCEID_T,
	"SR_INT8_T":               SR_INT8_T,
	"SR_INT16_T":              SR_INT16_T,
	"SR_INT32_T":              SR_INT32_T,
	"SR_INT64_T":              SR_INT64_T,
	"SR_STRING_T":             SR_STRING_T,
	"SR_UINT8_T":              SR_UINT8_T,
	"SR_UINT16_T":             SR_UINT16_T,
	"SR_UINT32_T":             SR_UINT32_T,
	"SR_UINT64_T":             SR_UINT64_T,
	"SR_ANYXML_T":             SR_ANYXML_T,
	"SR_ANYDATA_T":            SR_ANYDATA_T,
}

var srTypeNames = map[SrType]string{
	SR_UNKNOWN_T:            "SR_UNKNOWN_T",
	SR_TREE_ITERATOR_T:      "SR_TREE_ITERATOR_T",
	SR_LIST_T:               "SR_LIST_T",
	SR_CONTAINER_T:          "SR_CONTAINER_T",
	SR_CONTAINER_PRESENCE_T: "SR_CONTAINER_PRESENCE_T",
	SR_LEAF_EMPTY_T:         "SR_LEAF_EMPTY_T",
	SR_BINARY_T:             "SR_BINARY_T",
	SR_BITS_T:               "SR_BITS_T",
	SR_BOOL_T:               "SR_BOOL_T",
	SR_DECIMAL64_T:          "SR_DECIMAL64_T",
	SR_ENUM_T:               "SR_ENUM_T",
	SR_IDENTITYREF_T:        "SR_IDENTITYREF_T",
	SR_INSTANCEID_T:         "SR_INSTANCEID_T",
	SR_INT8_T:               "SR_INT8_T",
	SR_INT16_T:              "SR_INT16_T",
	SR_INT32_T:              "SR_INT32_T",
	SR_INT64_T:              "SR_INT64_T",
	SR_STRING_T:             "SR_STRING_T",
	SR_UINT8_T:              "SR_UINT8_T",
	SR_UINT16_T:             "SR_UINT16_T",
	SR_UINT32_T:             "SR_UINT32_T",
	SR_UINT64_T:             "SR_UINT64_T",
	SR_ANYXML_T:             "SR_ANYXML_T",
	SR_ANYDATA_T:            "SR_ANYDATA_T",
}

func (v SrType) String() string {
	if s, ok := srTypeNames[v]; ok {
		return s
	}
	return fmt.Sprintf("SrType(%d)", v)
}

//
// sr_error_t
//
type SrError int

const (
	SR_ERR_OK                SrError = C.SR_ERR_OK
	SR_ERR_INVAL_ARG         SrError = C.SR_ERR_INVAL_ARG
	SR_ERR_NOMEM             SrError = C.SR_ERR_NOMEM
	SR_ERR_NOT_FOUND         SrError = C.SR_ERR_NOT_FOUND
	SR_ERR_INTERNAL          SrError = C.SR_ERR_INTERNAL
	SR_ERR_INIT_FAILED       SrError = C.SR_ERR_INIT_FAILED
	SR_ERR_IO                SrError = C.SR_ERR_IO
	SR_ERR_DISCONNECT        SrError = C.SR_ERR_DISCONNECT
	SR_ERR_MALFORMED_MSG     SrError = C.SR_ERR_MALFORMED_MSG
	SR_ERR_UNSUPPORTED       SrError = C.SR_ERR_UNSUPPORTED
	SR_ERR_UNKNOWN_MODEL     SrError = C.SR_ERR_UNKNOWN_MODEL
	SR_ERR_BAD_ELEMENT       SrError = C.SR_ERR_BAD_ELEMENT
	SR_ERR_VALIDATION_FAILED SrError = C.SR_ERR_VALIDATION_FAILED
	SR_ERR_OPERATION_FAILED  SrError = C.SR_ERR_OPERATION_FAILED
	SR_ERR_DATA_EXISTS       SrError = C.SR_ERR_DATA_EXISTS
	SR_ERR_DATA_MISSING      SrError = C.SR_ERR_DATA_MISSING
	SR_ERR_UNAUTHORIZED      SrError = C.SR_ERR_UNAUTHORIZED
	SR_ERR_INVAL_USER        SrError = C.SR_ERR_INVAL_USER
	SR_ERR_LOCKED            SrError = C.SR_ERR_LOCKED
	SR_ERR_TIME_OUT          SrError = C.SR_ERR_TIME_OUT
	SR_ERR_RESTART_NEEDED    SrError = C.SR_ERR_RESTART_NEEDED
	SR_ERR_VERSION_MISMATCH  SrError = C.SR_ERR_VERSION_MISMATCH
)

var srErrorNames = map[SrError]string{
	SR_ERR_OK:                "SR_ERR_OK",
	SR_ERR_INVAL_ARG:         "SR_ERR_INVAL_ARG",
	SR_ERR_NOMEM:             "SR_ERR_NOMEM",
	SR_ERR_NOT_FOUND:         "SR_ERR_NOT_FOUND",
	SR_ERR_INTERNAL:          "SR_ERR_INTERNAL",
	SR_ERR_INIT_FAILED:       "SR_ERR_INIT_FAILED",
	SR_ERR_IO:                "SR_ERR_IO",
	SR_ERR_DISCONNECT:        "SR_ERR_DISCONNECT",
	SR_ERR_MALFORMED_MSG:     "SR_ERR_MALFORMED_MSG",
	SR_ERR_UNSUPPORTED:       "SR_ERR_UNSUPPORTED",
	SR_ERR_UNKNOWN_MODEL:     "SR_ERR_UNKNOWN_MODEL",
	SR_ERR_BAD_ELEMENT:       "SR_ERR_BAD_ELEMENT",
	SR_ERR_VALIDATION_FAILED: "SR_ERR_VALIDATION_FAILED",
	SR_ERR_OPERATION_FAILED:  "SR_ERR_OPERATION_FAILED",
	SR_ERR_DATA_EXISTS:       "SR_ERR_DATA_EXISTS",
	SR_ERR_DATA_MISSING:      "SR_ERR_DATA_MISSING",
	SR_ERR_UNAUTHORIZED:      "SR_ERR_UNAUTHORIZED",
	SR_ERR_INVAL_USER:        "SR_ERR_INVAL_USER",
	SR_ERR_LOCKED:            "SR_ERR_LOCKED",
	SR_ERR_TIME_OUT:          "SR_ERR_TIME_OUT",
	SR_ERR_RESTART_NEEDED:    "SR_ERR_RESTART_NEEDED",
	SR_ERR_VERSION_MISMATCH:  "SR_ERR_VERSION_MISMATCH",
}

var srErrorValues = map[string]SrError{
	"SR_ERR_OK":                SR_ERR_OK,
	"SR_ERR_INVAL_ARG":         SR_ERR_INVAL_ARG,
	"SR_ERR_NOMEM":             SR_ERR_NOMEM,
	"SR_ERR_NOT_FOUND":         SR_ERR_NOT_FOUND,
	"SR_ERR_INTERNAL":          SR_ERR_INTERNAL,
	"SR_ERR_INIT_FAILED":       SR_ERR_INIT_FAILED,
	"SR_ERR_IO":                SR_ERR_IO,
	"SR_ERR_DISCONNECT":        SR_ERR_DISCONNECT,
	"SR_ERR_MALFORMED_MSG":     SR_ERR_MALFORMED_MSG,
	"SR_ERR_UNSUPPORTED":       SR_ERR_UNSUPPORTED,
	"SR_ERR_UNKNOWN_MODEL":     SR_ERR_UNKNOWN_MODEL,
	"SR_ERR_BAD_ELEMENT":       SR_ERR_BAD_ELEMENT,
	"SR_ERR_VALIDATION_FAILED": SR_ERR_VALIDATION_FAILED,
	"SR_ERR_OPERATION_FAILED":  SR_ERR_OPERATION_FAILED,
	"SR_ERR_DATA_EXISTS":       SR_ERR_DATA_EXISTS,
	"SR_ERR_DATA_MISSING":      SR_ERR_DATA_MISSING,
	"SR_ERR_UNAUTHORIZED":      SR_ERR_UNAUTHORIZED,
	"SR_ERR_INVAL_USER":        SR_ERR_INVAL_USER,
	"SR_ERR_LOCKED":            SR_ERR_LOCKED,
	"SR_ERR_TIME_OUT":          SR_ERR_TIME_OUT,
	"SR_ERR_RESTART_NEEDED":    SR_ERR_RESTART_NEEDED,
	"SR_ERR_VERSION_MISMATCH":  SR_ERR_VERSION_MISMATCH,
}

func (v SrError) String() string {
	if s, ok := srErrorNames[v]; ok {
		return s
	}
	return fmt.Sprintf("SrError(%d)", v)
}

//
// sr_change_oper_t
//
type SrChangeOper int

const (
	SR_OP_CREATED  SrChangeOper = C.SR_OP_CREATED
	SR_OP_MODIFIED SrChangeOper = C.SR_OP_MODIFIED
	SR_OP_DELETED  SrChangeOper = C.SR_OP_DELETED
	SR_OP_MOVED    SrChangeOper = C.SR_OP_MOVED
)

var srChangeOperNames = map[SrChangeOper]string{
	SR_OP_CREATED:  "SR_OP_CREATED",
	SR_OP_MODIFIED: "SR_OP_MODIFIED",
	SR_OP_DELETED:  "SR_OP_DELETED",
	SR_OP_MOVED:    "SR_OP_MOVED",
}

var srChangeOperValues = map[string]SrChangeOper{
	"SR_OP_CREATED":  SR_OP_CREATED,
	"SR_OP_MODIFIED": SR_OP_MODIFIED,
	"SR_OP_DELETED":  SR_OP_DELETED,
	"SR_OP_MOVED":    SR_OP_MOVED,
}

func (v SrChangeOper) String() string {
	if s, ok := srChangeOperNames[v]; ok {
		return s
	}
	return fmt.Sprintf("SrChangeOper(%d)", v)
}

func NewSrChangeOperFromSr(s string) (SrChangeOper, error) {
	if v, ok := srChangeOperValues[s]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("unknown SrChangeOper %s", s)
}

//
// sr_notif_event_t
//
type SrNotifEvent int

const (
	SR_EV_VERIFY  SrNotifEvent = C.SR_EV_VERIFY
	SR_EV_APPLY   SrNotifEvent = C.SR_EV_APPLY
	SR_EV_ABORT   SrNotifEvent = C.SR_EV_ABORT
	SR_EV_ENABLED SrNotifEvent = C.SR_EV_ENABLED
)

var srNotifEventNames = map[SrNotifEvent]string{
	SR_EV_VERIFY:  "SR_EV_VERIFY",
	SR_EV_APPLY:   "SR_EV_APPLY",
	SR_EV_ABORT:   "SR_EV_ABORT",
	SR_EV_ENABLED: "SR_EV_ENABLED",
}

var srNotifEventValues = map[string]SrNotifEvent{
	"SR_EV_VERIFY":  SR_EV_VERIFY,
	"SR_EV_APPLY":   SR_EV_APPLY,
	"SR_EV_ABORT":   SR_EV_ABORT,
	"SR_EV_ENABLED": SR_EV_ENABLED,
}

func (v SrNotifEvent) String() string {
	if s, ok := srNotifEventNames[v]; ok {
		return s
	}
	return fmt.Sprintf("SrNotifEvent(%d)", v)
}

//
// sr_datastore_t
//
type SrDataStore int

const (
	SR_DS_STARTUP   SrDataStore = C.SR_DS_STARTUP
	SR_DS_RUNNING   SrDataStore = C.SR_DS_RUNNING
	SR_DS_CANDIDATE SrDataStore = C.SR_DS_CANDIDATE
)

var srDataStoreNames = map[SrDataStore]string{
	SR_DS_STARTUP:   "SR_DS_STARTUP",
	SR_DS_RUNNING:   "SR_DS_RUNNING",
	SR_DS_CANDIDATE: "SR_DS_CANDIDATE",
}

var stDataStoreValues = map[string]SrDataStore{
	"SR_DS_STARTUP":   SR_DS_STARTUP,
	"SR_DS_RUNNING":   SR_DS_RUNNING,
	"SR_DS_CANDIDATE": SR_DS_CANDIDATE,
}

func (v SrDataStore) String() string {
	if s, ok := srDataStoreNames[v]; ok {
		return s
	}
	return fmt.Sprintf("SrDataStore(%d)", v)
}

func (v SrDataStore) C() C.sr_datastore_t {
	return C.sr_datastore_t(v)
}

//
// sr_edit_options_t
//
type SrEditOptions uint32

const (
	SR_EDIT_DEFAULT       SrEditOptions = C.SR_EDIT_DEFAULT
	SR_EDIT_NON_RECURSIVE SrEditOptions = C.SR_EDIT_NON_RECURSIVE
	SR_EDIT_STRICT        SrEditOptions = C.SR_EDIT_STRICT
)

var srEditOptionsNames = map[SrEditOptions]string{
	SR_EDIT_DEFAULT:       "SR_EDIT_DEFAULT",
	SR_EDIT_NON_RECURSIVE: "SR_EDIT_NON_RECURSIVE",
	SR_EDIT_STRICT:        "SR_EDIT_STRICT",
}

var srEditOptionsValues = map[string]SrEditOptions{
	"SR_EDIT_DEFAULT":       SR_EDIT_DEFAULT,
	"SR_EDIT_NON_RECURSIVE": SR_EDIT_NON_RECURSIVE,
	"SR_EDIT_STRICT":        SR_EDIT_STRICT,
}

func (v SrEditOptions) String() string {
	names := []string{}
	for value, name := range srEditOptionsNames {
		if (v & value) != 0 {
			names = append(names, name)
		}
	}
	return strings.Join(names, "|")
}

//
// sr_subscr_flag_t
//
type SrSubscrFlag uint32

const (
	SR_SUBSCR_DEFAULT                  SrSubscrFlag = C.SR_SUBSCR_DEFAULT
	SR_SUBSCR_CTX_REUSE                SrSubscrFlag = C.SR_SUBSCR_CTX_REUSE
	SR_SUBSCR_PASSIVE                  SrSubscrFlag = C.SR_SUBSCR_PASSIVE
	SR_SUBSCR_APPLY_ONLY               SrSubscrFlag = C.SR_SUBSCR_APPLY_ONLY
	SR_SUBSCR_EV_ENABLED               SrSubscrFlag = C.SR_SUBSCR_EV_ENABLED
	SR_SUBSCR_NO_ABORT_FOR_REFUSED_CFG SrSubscrFlag = C.SR_SUBSCR_NO_ABORT_FOR_REFUSED_CFG
	SR_SUBSCR_NOTIF_REPLAY_FIRST       SrSubscrFlag = C.SR_SUBSCR_NOTIF_REPLAY_FIRST
)

var srSubscrFlagNames = map[SrSubscrFlag]string{
	SR_SUBSCR_DEFAULT:                  "DEFAULT",
	SR_SUBSCR_CTX_REUSE:                "CTX_REUSE",
	SR_SUBSCR_PASSIVE:                  "PASSIVE",
	SR_SUBSCR_APPLY_ONLY:               "APPLY_ONLY",
	SR_SUBSCR_EV_ENABLED:               "EV_ENABLED",
	SR_SUBSCR_NO_ABORT_FOR_REFUSED_CFG: "NO_ABORT_FOR_REFUSED_CFG",
	SR_SUBSCR_NOTIF_REPLAY_FIRST:       "NOTIF_REPLAY_FIRST",
}

var srSubscrFlagValues = map[string]SrSubscrFlag{
	"DEFAULT":                  SR_SUBSCR_DEFAULT,
	"CTX_REUSE":                SR_SUBSCR_CTX_REUSE,
	"PASSIVE":                  SR_SUBSCR_PASSIVE,
	"APPLY_ONLY":               SR_SUBSCR_APPLY_ONLY,
	"EV_ENABLED":               SR_SUBSCR_EV_ENABLED,
	"NO_ABORT_FOR_REFUSED_CFG": SR_SUBSCR_NO_ABORT_FOR_REFUSED_CFG,
	"NOTIF_REPLAY_FIRST":       SR_SUBSCR_NOTIF_REPLAY_FIRST,
}

func (v SrSubscrFlag) String() string {
	names := []string{}
	for value, name := range srSubscrFlagNames {
		if (v & value) != 0 {
			names = append(names, name)
		}
	}
	return strings.Join(names, "|")
}

func (v SrSubscrFlag) C() C.sr_subscr_options_t {
	return C.sr_subscr_options_t(v)
}

func ParseSrSubscrFlag(s string) (SrSubscrFlag, error) {
	if v, ok := srSubscrFlagValues[s]; ok {
		return v, nil
	}
	return SR_SUBSCR_DEFAULT, fmt.Errorf("Invalid SrSubscrFlag. %s", s)
}

func JoinSrSubscrFlags(flags ...SrSubscrFlag) SrSubscrFlag {
	var v SrSubscrFlag = 0
	for _, flag := range flags {
		v = v | flag
	}
	return v
}

//
// log_lecel
//
type SrLogLevel int

const (
	SR_LL_NONE SrLogLevel = C.SR_LL_NONE
	SR_LL_ERR  SrLogLevel = C.SR_LL_ERR
	SR_LL_WRN  SrLogLevel = C.SR_LL_WRN
	SR_LL_INF  SrLogLevel = C.SR_LL_INF
	SR_LL_DBG  SrLogLevel = C.SR_LL_DBG
)

var srLogLevelNames = map[SrLogLevel]string{
	SR_LL_NONE: "NONE",
	SR_LL_ERR:  "ERR",
	SR_LL_WRN:  "WRN",
	SR_LL_INF:  "INF",
	SR_LL_DBG:  "CBG",
}

var srLogLevelValues = map[string]SrLogLevel{
	"NONE": SR_LL_NONE,
	"ERR":  SR_LL_ERR,
	"WRN":  SR_LL_WRN,
	"INF":  SR_LL_INF,
	"CBG":  SR_LL_DBG,
}

func (v SrLogLevel) String() string {
	if s, ok := srLogLevelNames[v]; ok {
		return s
	}
	return fmt.Sprintf("SrLogLevel(%d)", v)
}

func (v SrLogLevel) C() C.sr_log_level_t {
	return C.sr_log_level_t(v)
}

func ParseSrLogLevel(s string) (SrLogLevel, error) {
	if v, ok := srLogLevelValues[s]; ok {
		return v, nil
	}
	return SR_LL_NONE, fmt.Errorf("Invalid SrLogLevel. %s", s)
}
