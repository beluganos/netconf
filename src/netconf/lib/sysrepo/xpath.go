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
#include <sysrepo/xpath.h>
*/
import "C"
import (
	"netconf/lib/xml"
	"unsafe"
)

func xpathAttrs(state *C.sr_xpath_ctx_t) map[string]string {
	attrs := map[string]string{}
	for {
		c_name := C.sr_xpath_next_key_name(nil, state)
		if c_name == nil {
			break
		}
		name := C.GoString(c_name)

		value := func() string {
			c_value := C.sr_xpath_next_key_value(nil, state)
			if c_value == nil {
				return ""
			}
			return C.GoString(c_value)
		}()

		attrs[name] = value
	}
	return attrs
}

func ParseXPath(xpath string) []*ncxml.XPathNode {
	nodes := []*ncxml.XPathNode{}
	ParseXPathNodes(xpath, func(node *ncxml.XPathNode) bool {
		nodes = append(nodes, node)
		return true
	})
	return nodes
}

func ParseXPathNodes(xpath string, f func(*ncxml.XPathNode) bool) {
	c_xpath := C.CString(xpath)
	defer C.free(unsafe.Pointer(c_xpath))

	var state C.sr_xpath_ctx_t
	c_node := C.sr_xpath_next_node_with_ns(c_xpath, &state)
	for c_node != nil {
		ns, name := ncxml.ParseXPathName(C.GoString(c_node))
		node := ncxml.NewXPathNode(ns, name, xpathAttrs(&state))
		if ok := f(node); !ok {
			break
		}
		c_node = C.sr_xpath_next_node_with_ns(nil, &state)
	}
}
