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
// sr_connect_ctx
//
type SrConnection struct {
	conn *C.sr_conn_ctx_t
}

func NewSrConnection() *SrConnection {
	return &SrConnection{
		conn: nil,
	}
}

func (c *SrConnection) Connect(appName string) error {
	c_appname := C.CString(appName)
	defer C.free(unsafe.Pointer(c_appname))

	rc := C.sr_connect(c_appname, C.SR_CONN_DEFAULT, &c.conn)
	if rc != C.SR_ERR_OK {
		return fmt.Errorf("sr_connect error. %d", rc)
	}

	return nil
}

func (c *SrConnection) Disconnect() {
	if c.conn != nil {
		C.sr_disconnect(c.conn)
		c.conn = nil
	}
}
