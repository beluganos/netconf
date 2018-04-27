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

package openconfig

import (
	"fmt"
	"netconf/lib"
	"netconf/lib/xml"
)

const (
	MPLS_KEY           = "mpls"
	MPLS_NULLLABEL_KEY = "null-label"
	MPLS_IFATTRS_KEY   = "interface-attributes"
	MPLS_SIGPROTOS_KEY = "signaling-protocols"
	MPLS_LDP_KEY       = "ldp"
)

//
// mpls
//
type Mpls struct {
	nclib.SrChanges `xml:"-"`

	Global    *MplsGlobal       `xml:"global"`
	SigProtos *MplsSigProtocols `xml:"signaling-protocols"`
}

type MplsProcessor interface {
	MplsGlobalProcessor
	MplsSigProtocolsProcessor
}

func NewMpls() *Mpls {
	return &Mpls{
		SrChanges: nclib.NewSrChanges(),
		Global:    NewMplsGlobal(),
		SigProtos: NewMplsSigProtocols(),
	}
}

func (m *Mpls) String() string {
	return fmt.Sprintf("%s{%s, %s} %s",
		MPLS_KEY,
		m.Global,
		m.SigProtos,
		m.SrChanges,
	)
}

func (m *Mpls) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_GLOBAL_KEY:
		if err := m.Global.Put(nodes[1:], value); err != nil {
			return err
		}

	case MPLS_SIGPROTOS_KEY:
		if err := m.SigProtos.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMpls(p MplsProcessor, reverse bool, name string, mpls *Mpls) error {
	globalFunc := func() error {
		if mpls.GetChange(OC_GLOBAL_KEY) {
			return ProcessMplsGlobal(
				p.(MplsGlobalProcessor),
				reverse,
				name,
				mpls.Global,
			)
		}
		return nil
	}

	spFunc := func() error {
		if mpls.GetChange(MPLS_SIGPROTOS_KEY) {
			return ProcessMplsSigProtocols(
				p.(MplsSigProtocolsProcessor),
				reverse,
				name,
				mpls.SigProtos,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, globalFunc, spFunc)
}
