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

package ncxml

import (
	"fmt"
	"strings"
)

type XPathNode struct {
	Name  string
	Ns    string
	Attrs map[string]string
}

func NewXPathNode(ns, name string, attrs map[string]string) *XPathNode {
	return &XPathNode{
		Name:  name,
		Ns:    ns,
		Attrs: attrs,
	}
}

func (n *XPathNode) NodeName() string {
	if len(n.Ns) == 0 {
		return n.Name
	}
	return fmt.Sprintf("%s:%s", n.Ns, n.Name)
}

func (n *XPathNode) Attr(name string, defaultVal string) string {
	if v, ok := n.Attrs[name]; ok {
		return v
	}
	return defaultVal
}

func XPathNodeNames(nodes []*XPathNode) []string {
	names := make([]string, len(nodes))
	for i, node := range nodes {
		names[i] = node.Name
	}
	return names
}

func XPathNodeNamesWithNs(nodes []*XPathNode) []string {
	names := make([]string, len(nodes))
	for i, node := range nodes {
		names[i] = node.NodeName()
	}
	return names
}

func NewXPathFromNodes(nodes []*XPathNode) string {
	names := XPathNodeNames(nodes)
	return fmt.Sprintf("/%s", strings.Join(names, "/"))
}

func NewXPathFromNodesWithNs(nodes []*XPathNode) string {
	names := XPathNodeNamesWithNs(nodes)
	return "/" + strings.Join(names, "/")
}

func IndexOfXPathNode(name string, nodes []*XPathNode, start int) int {
	for start < len(nodes) {
		if node := nodes[start]; name == node.Name {
			return start
		}
		start++
	}
	return -1
}

func IndexOfXPathNodeWithNs(ns, name string, nodes []*XPathNode, start int) int {
	for start < len(nodes) {
		if node := nodes[start]; name == node.Name && ns == node.Ns {
			return start
		}
		start++
	}
	return -1
}
