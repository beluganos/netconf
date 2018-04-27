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

package ncmdbm

import (
	"fmt"
	"netconf/lib/openconfig"
	"netconf/lib/sysrepo"
)

//
// Interface Table
//
type InterfaceTable struct {
	session *srlib.SrSession
}

func NewInterfaceTable(session *srlib.SrSession) *InterfaceTable {
	return &InterfaceTable{
		session: session,
	}
}

func (t *InterfaceTable) Select(name string) (*openconfig.Interface, error) {

	xpath := fmt.Sprintf("/%s:%s/%s[%s='%s']//*",
		openconfig.INTERFACES_MODULE, openconfig.INTERFACES_KEY,
		openconfig.INTERFACE_KEY, openconfig.OC_NAME_KEY, name,
	)

	ifaces := openconfig.NewInterfaces()
	for cv := range t.session.GetItems(xpath) {
		cv.Dispatch(ifaces, nil, nil)
	}

	iface, ok := ifaces[name]
	if !ok {
		return nil, fmt.Errorf("Interface not found. %s", name)
	}

	return iface, nil
}

func (t *InterfaceTable) Walk(f func(string, *openconfig.Interface)) {

	xpath := fmt.Sprintf("/%s:%s//*",
		openconfig.INTERFACES_MODULE, openconfig.INTERFACES_KEY,
	)

	ifaces := openconfig.NewInterfaces()
	for cv := range t.session.GetItems(xpath) {
		cv.Dispatch(ifaces, nil, nil)
	}

	for name, iface := range ifaces {
		f(name, iface)
	}
}
