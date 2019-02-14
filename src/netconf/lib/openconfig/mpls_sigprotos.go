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
	nclib "netconf/lib"
	ncxml "netconf/lib/xml"
)

//
// mpls/signaling-protocols
//
type MplsSigProtocols struct {
	nclib.SrChanges `xml:"-"`

	Ldp *MplsLdp `xml:"ldp"`
}

type MplsSigProtocolsProcessor interface {
	MplsLdpProcessor
}

func NewMplsSigProtocols() *MplsSigProtocols {
	return &MplsSigProtocols{
		SrChanges: nclib.NewSrChanges(),
		Ldp:       NewMplsLdp(),
	}
}

func (m *MplsSigProtocols) String() string {
	return fmt.Sprintf("%s{%s} %s",
		MPLS_SIGPROTOS_KEY,
		m.Ldp,
		m.SrChanges,
	)
}

func (m *MplsSigProtocols) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_LDP_KEY:
		if err := m.Ldp.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsSigProtocols(p MplsSigProtocolsProcessor, reverse bool, name string, protos *MplsSigProtocols) error {
	ldpFunc := func() error {
		if protos.GetChange(MPLS_LDP_KEY) {
			return ProcessMplsLdp(
				p.(MplsLdpProcessor),
				reverse,
				name,
				protos.Ldp,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, ldpFunc)
}
