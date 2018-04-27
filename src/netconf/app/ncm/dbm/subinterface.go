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
	"netconf/lib/net"
	"netconf/lib/openconfig"
	"netconf/lib/sysrepo"
)

//
// Subinterface Table
//
type SubinterfaceTable struct {
	session *srlib.SrSession
}

func NewSubinterfaceTable(session *srlib.SrSession) *SubinterfaceTable {
	return &SubinterfaceTable{
		session: session,
	}
}

func (t *SubinterfaceTable) SelectById(ifaceId string) (*openconfig.Subinterface, string, error) {
	name, index, err := ncnet.ParseIFName(ifaceId)
	if err != nil {
		return nil, "", err
	}

	return t.Select(name, index)
}

func (t *SubinterfaceTable) Select(name string, index uint32) (*openconfig.Subinterface, string, error) {

	xpath := fmt.Sprintf("/%s:%s/%s[%s='%s']/%s/%s[%s=%d]//*",
		openconfig.INTERFACES_MODULE, openconfig.INTERFACES_KEY,
		openconfig.INTERFACE_KEY, openconfig.OC_NAME_KEY, name,
		openconfig.SUBINTERFACES_KEY,
		openconfig.SUBINTERFACE_KEY, openconfig.OC_INDEX_KEY, index,
	)

	ifaces := openconfig.NewInterfaces()
	for cv := range t.session.GetItems(xpath) {
		if err := cv.Dispatch(ifaces, nil, nil); err != nil {
			continue
		}
	}

	iface, ok := ifaces[name]
	if !ok {
		return nil, name, fmt.Errorf("Interface not found. %s/%d", name, index)
	}

	subif, ok := iface.Subinterfaces[index]
	if !ok {
		return nil, name, fmt.Errorf("Subinterface not found. %s/%d", name, index)
	}

	return subif, name, nil
}

func (t *SubinterfaceTable) Walk(f func(string, *openconfig.Subinterface)) {

	xpath := fmt.Sprintf("/%s:%s//*",
		openconfig.INTERFACES_MODULE, openconfig.INTERFACES_KEY,
	)

	ifaces := openconfig.NewInterfaces()
	for cv := range t.session.GetItems(xpath) {
		if err := cv.Dispatch(ifaces, nil, nil); err != nil {
			continue
		}
	}

	for ifname, iface := range ifaces {
		for _, subif := range iface.Subinterfaces {
			f(ifname, subif)
		}
	}
}
