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

func convInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		return int(v.(int64))
	}
	return 0
}

func convUint(m map[string]interface{}, key string) uint {
	if v, ok := m[key]; ok {
		return uint(v.(int64))
	}
	return 0
}

func convString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		return v.(string)
	}
	return ""
}

func convBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		return v.(bool)
	}
	return false
}

func convList(m map[string]interface{}, key string) []interface{} {
	if v, ok := m[key]; ok {
		return v.([]interface{})
	}
	return nil
}

func getValue(m map[string]interface{}, key string) interface{} {
	if val, ok := m[key]; ok {
		return val
	}
	return nil
}

type Entries map[string]interface{}

func NewEntries(i interface{}) Entries {
	switch i.(type) {
	case nil:
		return Entries{}
	default:
		return i.(map[string]interface{})
	}
}

func (e Entries) Raw() map[string]interface{} {
	return e
}
