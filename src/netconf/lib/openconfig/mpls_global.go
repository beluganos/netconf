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
// mpls/global
type MplsGlobal struct {
	nclib.SrChanges `xml:"-"`

	Config     *MplsConfig        `xml:"config"`
	IfaceAttrs MplsInterfaceAttrs `xml:"interface-attributes"`
}

type MplsGlobalProcessor interface {
	MplsConfigProcessor
	MplsInterfaceAttrProcessor
}

func NewMplsGlobal() *MplsGlobal {
	return &MplsGlobal{
		SrChanges:  nclib.NewSrChanges(),
		Config:     NewMplsConfig(),
		IfaceAttrs: NewMplsInterfaceAttrs(),
	}
}

func (m *MplsGlobal) String() string {
	return fmt.Sprintf("%s{%s, %s} %s",
		OC_GLOBAL_KEY,
		m.Config,
		m.IfaceAttrs,
		m.SrChanges,
	)
}

func (m *MplsGlobal) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := m.Config.Put(nodes[1:], value); err != nil {
			return nil
		}

	case MPLS_IFATTRS_KEY:
		if err := m.IfaceAttrs.Put(nodes[1:], value); err != nil {
			return nil
		}
	}

	m.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsGlobal(p MplsGlobalProcessor, reverse bool, name string, global *MplsGlobal) error {
	configFunc := func() error {
		if global.GetChange(OC_CONFIG_KEY) {
			return ProcessMplsConfig(
				p.(MplsConfigProcessor),
				reverse,
				name,
				global.Config,
			)
		}
		return nil
	}

	attrsFunc := func() error {
		if global.GetChange(MPLS_IFATTRS_KEY) {
			return ProcessMplsInterfaceAttrs(
				p.(MplsInterfaceAttrProcessor),
				reverse,
				name,
				global.IfaceAttrs,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc, attrsFunc)
}

//
// mpls/global/config
//
type MplsConfig struct {
	nclib.SrChanges `xml:"-"`

	NullLabel MplsNullLabelType `xml:"null-label"`
}

type MplsConfigProcessor interface {
	MplsConfig(string, *MplsConfig) error
}

func NewMplsConfig() *MplsConfig {
	return &MplsConfig{
		SrChanges: nclib.NewSrChanges(),
		NullLabel: 0,
	}
}

func (c *MplsConfig) String() string {
	return fmt.Sprintf("%s{%s=%s} %s",
		OC_CONFIG_KEY,
		MPLS_NULLLABEL_KEY, c.NullLabel,
		c.SrChanges,
	)
}

func (c *MplsConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case MPLS_NULLLABEL_KEY:
		v, err := ParseMplsNullLabelType(value)
		if err != nil {
			return err
		}
		c.NullLabel = v
	}

	c.SetChange(nodes[0].Name)
	return nil
}

func ProcessMplsConfig(p MplsConfigProcessor, reverse bool, name string, config *MplsConfig) error {
	configFunc := func() error {
		return p.MplsConfig(name, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
