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

package ncgobgp

import (
	"fmt"
	"testing"
)

func TestViper(t *testing.T) {
	f := "test/gobgp_test.toml"
	cfg, err := ReadConfigFile(f, "toml")
	if err != nil {
		t.Errorf("ReadConfig error. %s", err)
	}

	for _, key := range cfg.AllKeys() {
		fmt.Printf("KEY:%s VALUE:%v\n", key, cfg.Get(key))
	}

	g := cfg.Global()
	fmt.Printf("global           %v\n", g)
	fmt.Printf("global.config    %v\n", g.Config())
	fmt.Printf("global.config.as %v\n", g.Config().As())

	g.Config().SetAs(100)
	cfg.SetGlobal(g)

	zebra := cfg.Zebra()
	zebra.Config().SetEnabled(true)
	zebra.Config().SetVersion(4)
	zebra.Config().SetRedistributeRouteTypeList([]string{"connected"})
	cfg.SetZebra(zebra)

	if err := cfg.WriteConfigAs("test/gobgp_test.new.json"); err != nil {
		t.Errorf("WriteConfigAs error. %s", err)
	}

	cfg.SetZebra(NewZebra(nil))
	if err := cfg.WriteConfigAs("test/gobgp_test.nozebra.json"); err != nil {
		t.Errorf("WriteConfigAs error. %s", err)
	}

	for _, n := range cfg.Neighbors() {
		fmt.Printf("neighbor: %v\n", n)
		fmt.Printf("neighbor.config: %v\n", n.Config())
		fmt.Printf("neighbor.transport: %v\n", n.Transport())
		for _, afisafi := range n.AfiSafis() {
			fmt.Printf("neighbor.afisafi: %v\n", afisafi)
		}
	}

	for _, d := range cfg.PolicyDefinitions() {
		fmt.Printf("policy-definition: %v\n", d)
	}
}
