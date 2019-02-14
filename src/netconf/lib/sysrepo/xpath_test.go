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

import (
	ncxml "netconf/lib/xml"
	"testing"
)

func TestXpath(t *testing.T) {
	xpath := "/local-routing:static-routes/static[ip='192.168.122.0'][prefix-length='24']/name"
	nodes := ParseXPath(xpath)

	for _, node := range nodes {
		t.Log(node)
	}

	if v := len(nodes); v != 3 {
		t.Errorf("ParseXPath unmatch. #nodes=%d", v)
	}

	// nodes[0]
	if v := nodes[0].Name; v != "static-routes" {
		t.Errorf("ParseXPath unmatch. node[0].Name=%s", v)
	}
	if v := nodes[0].Attrs; len(v) != 0 {
		t.Errorf("ParseXPath unmatch. node[0].Attrs=%s", v)
	}

	// nodes[1]
	if v := nodes[1].Name; v != "static" {
		t.Errorf("ParseXPath unmatch. node[1].Name=%s", v)
	}
	if v := nodes[1].Attrs; len(v) != 2 {
		t.Errorf("ParseXPath unmatch. node[1].Attrs=%s", v)
	}
	if v, ok := nodes[1].Attrs["ip"]; !ok || v != "192.168.122.0" {
		t.Errorf("ParseXPath unmatch. node[1].Attrs[ip]=%s", v)
	}
	if v, ok := nodes[1].Attrs["prefix-length"]; !ok || v != "24" {
		t.Errorf("ParseXPath unmatch. node[1].Attrs[prefix-length]=%s", v)
	}

	// nodes[2]
	if v := nodes[2].Name; v != "name" {
		t.Errorf("ParseXPath unmatch. node[2].Name=%s", v)
	}
	if v := nodes[2].Attrs; len(v) != 0 {
		t.Errorf("ParseXPath unmatch. node[2].Attrs=%s", v)
	}
}

func TestNewXPathFromNodes(t *testing.T) {
	xpath := "/ns1:node1/node2/node3[name='NODE3']/ns2:node4[id='ID5']"
	nodes := ParseXPath(xpath)

	s := ncxml.NewXPathFromNodes(nodes)

	if s != "/node1/node2/node3/node4" {
		t.Errorf("NewXPathFromNodes unmatch. '%s'", s)
	}
}

func TestNewXPathFromNodesWithNs(t *testing.T) {
	xpath := "/ns1:node1/node2/node3[name='NODE3']/ns2:node4[id='ID5']"
	nodes := ParseXPath(xpath)

	s := ncxml.NewXPathFromNodesWithNs(nodes)

	if s != "/ns1:node1/node2/node3/ns2:node4" {
		t.Errorf("NewXPathFromNodesWithNs unmatch. '%s'", s)
	}
}
