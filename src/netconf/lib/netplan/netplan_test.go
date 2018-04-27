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

package ncnplib

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	c, err := ReadConfigFile("netplan_test.yaml")

	if err != nil {
		t.Errorf("Read error. %s", err)
	}

	fmt.Println(c)

	b := &bytes.Buffer{}
	if err := WriteConfig(b, c); err != nil {
		t.Errorf("Write error. %s", err)
	}

	fmt.Println(string(b.Bytes()))
}