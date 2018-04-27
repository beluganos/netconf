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

//
// sr_log_stderr
//
func SrLogStderr(level SrLogLevel) {
	C.sr_log_stderr(level.C())
}

//
// sr_log_syslog
//
func SrLogSyslog(level SrLogLevel) {
	C.sr_log_syslog(level.C())
}

var logCb func(SrLogLevel, string) = nil

//
// sr_log_Set_cb
//
func SrLogSetCb(f func(SrLogLevel, string)) {
	logCb = f
}

//export Go_log_cb
func Go_log_cb(c_level C.sr_log_level_t, c_message *C.char) {
	if logCb != nil {
		logCb(SrLogLevel(c_level), C.GoString(c_message))
	}

}
