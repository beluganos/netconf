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

package srocgobgp

import (
	"fmt"
	"netconf/lib/openconfig"
	"netconf/lib/sysrepo"
	"testing"
)

func TestProcessor(t *testing.T) {
	var gp openconfig.BgpProcessor
	var pp openconfig.PolicyDefinitionProcessor

	p := NewConfigProcessor()
	gp = p
	pp = p

	if pp == nil {
		t.Errorf("NewConfigProcessor umatch interface.")
	}

	if gp == nil {
		t.Errorf("NewConfigProcessor umatch interface.")
	}
}

func makeBgp(bgp *openconfig.Bgp, xpaths map[string]string) error {
	for xpath, value := range xpaths {
		nodes := srlib.ParseXPath(xpath)
		if err := bgp.Put(nodes[1:], value); err != nil {
			return err
		}
	}
	return nil
}

func cmpLines(src []string, dst []string) error {
	dumpLines(src)

	if len(src) != len(dst) {
		return fmt.Errorf("unmatch line number")
	}

	for index, s := range src {
		if s != dst[index] {
			return fmt.Errorf("unmatch line. %d %s %s", index, s, dst[index])
		}
	}

	return nil
}

func dumpLines(lines []string) {
	for index, line := range lines {
		fmt.Printf("[%d] '%s'\n", index, line)
	}
}

func TestProcessBgpGlobalConfig(t *testing.T) {
	xpaths := map[string]string{
		"/bgp":                         "",
		"/bgp/global":                  "",
		"/bgp/global/config/as":        "65000",
		"/bgp/global/config/router-id": "10.10.10.10",
	}

	d := []string{
		"[global.config]",
		"as = 65000",
		"router-id = \"10.10.10.10\"",
	}

	bgp := openconfig.NewBgp()
	if err := makeBgp(bgp, xpaths); err != nil {
		t.Errorf("makeBgp error. %s", err)
	}

	p := NewConfigProcessor()
	if err := openconfig.ProcessBgp(p, false, "", nil, bgp); err != nil {
		t.Errorf("ProcessNetworkInstances error. %s", err)
	}

	items := p.Items()
	if err := cmpLines(items, d); err != nil {
		t.Errorf("cmpLines error. %s", err)
	}
}

func TestProcessZebra(t *testing.T) {
	xpaths := map[string]string{
		"/bgp/zebra":                            "",
		"/bgp/zebra/config/enabled":             "true",
		"/bgp/zebra/config/version":             "4",
		"/bgp/zebra/config/url":                 "unix:/var/run/frr/zserv.api",
		"/bgp/zebra/config/redistribute-routes": "DIRECTLY_CONNECTED",
	}

	d := []string{
		"[zebra.config]",
		"enabled = true",
		"version = 4",
		"url = \"unix:/var/run/frr/zserv.api\"",
		"redistribute-route-type-list = [connected]",
	}

	bgp := openconfig.NewBgp()
	if err := makeBgp(bgp, xpaths); err != nil {
		t.Errorf("makeBgp error. %s", err)
	}

	p := NewConfigProcessor()
	if err := openconfig.ProcessBgp(p, false, "", nil, bgp); err != nil {
		t.Errorf("ProcessNetworkInstances error. %s", err)
	}

	items := p.Items()
	if err := cmpLines(items, d); err != nil {
		t.Errorf("cmpLines error. %s", err)
	}

}
