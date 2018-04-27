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

package cfglxclib

import (
	"fmt"
	"strings"
)

type Operation int

const (
	OP_HELP Operation = iota
	OP_SAVE
	OP_LOAD
	OP_CLEAR
	OP_ADD
	OP_DEL
	OP_ADD_IFACE
	OP_DEL_IFACE
)

var opNames = map[Operation]string{
	OP_HELP:      "help",
	OP_SAVE:      "save",
	OP_LOAD:      "load",
	OP_CLEAR:     "clear",
	OP_ADD:       "add",
	OP_DEL:       "del",
	OP_ADD_IFACE: "addif",
	OP_DEL_IFACE: "delif",
}

var opValues = map[string]Operation{
	"help":  OP_HELP,
	"save":  OP_SAVE,
	"load":  OP_LOAD,
	"clear": OP_CLEAR,
	"add":   OP_ADD,
	"del":   OP_DEL,
	"addif": OP_ADD_IFACE,
	"delif": OP_DEL_IFACE,
}

func StrOperations() string {
	ss := make([]string, len(opNames))
	for i := 0; i < len(opNames); i++ {
		ss[i] = Operation(i).String()
	}
	return strings.Join(ss, "/")
}

func ParseOperation(s string) (Operation, error) {
	if op, ok := opValues[s]; ok {
		return op, nil
	}
	return OP_HELP, fmt.Errorf("Invalid operation. %s", s)
}

func (o Operation) String() string {
	if s, ok := opNames[o]; ok {
		return s
	}
	return fmt.Sprintf("Operation(%d)", o)
}
