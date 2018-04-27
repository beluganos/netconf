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

package nclib

import (
	"testing"
)

func TestChanges_empty(t *testing.T) {
	c := NewSrChanges()

	if ok := c.Compare(); !ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}
}

func TestChanges_x1(t *testing.T) {
	c := NewSrChanges()

	c.SetChange("test1")
	if ok := c.GetChange("test1"); !ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.GetChange("test2"); ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.Compare("test1"); !ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.Compare(); ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.Compare("test1", "test2"); ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}
}

func TestChanges_x2(t *testing.T) {
	c := NewSrChanges()

	c.SetChange("test1")
	c.SetChange("test2")
	if ok := c.GetChange("test1"); !ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.GetChange("test2"); !ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.GetChange("test3"); ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.Compare("test1"); ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.Compare("test1", "test2"); !ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}

	if ok := c.Compare("test1", "test3"); ok {
		t.Errorf("SrChanges.Compare unmatch. %t", ok)
	}
}

func TestChanges_string(t *testing.T) {
	c := NewSrChanges()

	if s := c.String(); s != "" {
		t.Errorf("SrChanges.String unmatch. %s", s)
	}

	c.SetChange("test1")

	if s := c.String(); s != "test1" {
		t.Errorf("SrChanges.String unmatch. %s", s)
	}

	c.SetChange("test2")

	if s := c.String(); s != "test1|test2" && s != "test2|test1" {
		t.Errorf("SrChanges.String unmatch. %s", s)
	}
}

func TestChanges_changes(t *testing.T) {
	c := NewSrChanges()

	if ok := c.GetChanges(); !ok {
		t.Errorf("SrChanges.GetChanges unmatch. %t", ok)
	}

	c.SetChanges("test1", "test2")

	if ok := c.GetChanges("test1"); !ok {
		t.Errorf("SrChanges.GetChanges unmatch. %t", ok)
	}

	if ok := c.GetChanges("test2"); !ok {
		t.Errorf("SrChanges.GetChanges unmatch. %t", ok)
	}

	if ok := c.GetChanges("test2", "test1"); !ok {
		t.Errorf("SrChanges.GetChanges unmatch. %t", ok)
	}

	if ok := c.GetChanges("test"); ok {
		t.Errorf("SrChanges.GetChanges unmatch. %t", ok)
	}

	if ok := c.GetChanges("test1", "test3"); ok {
		t.Errorf("SrChanges.GetChanges unmatch. %t", ok)
	}
}

func TestChanges_oneof(t *testing.T) {
	c := NewSrChanges()
	c.SetChanges("test1", "test2")

	if ok := c.OneOfChange(); ok {
		t.Errorf("StChanges.OneOfChange must be unmatch. %t", ok)
	}

	if ok := c.OneOfChange("dummy1"); ok {
		t.Errorf("StChanges.OneOfChange must be unmatch. %t", ok)
	}

	if ok := c.OneOfChange("dummy1", "dummy2"); ok {
		t.Errorf("StChanges.OneOfChange must be unmatch. %t", ok)
	}

	if ok := c.OneOfChange("test1"); !ok {
		t.Errorf("StChanges.OneOfChange unmatch. %t", ok)
	}

	if ok := c.OneOfChange("test1", "dummy"); !ok {
		t.Errorf("StChanges.OneOfChange unmatch. %t", ok)
	}

	if ok := c.OneOfChange("dummy", "test1"); !ok {
		t.Errorf("StChanges.OneOfChange unmatch. %t", ok)
	}

	if ok := c.OneOfChange("test1", "test2"); !ok {
		t.Errorf("StChanges.OneOfChange unmatch. %t", ok)
	}

	if ok := c.OneOfChange("dummy", "test1", "test2"); !ok {
		t.Errorf("StChanges.OneOfChange unmatch. %t", ok)
	}

	if ok := c.OneOfChange("test1", "dummy", "test2"); !ok {
		t.Errorf("StChanges.OneOfChange unmatch. %t", ok)
	}

	if ok := c.OneOfChange("test1", "test2", "dummy"); !ok {
		t.Errorf("StChanges.OneOfChange unmatch. %t", ok)
	}
}
