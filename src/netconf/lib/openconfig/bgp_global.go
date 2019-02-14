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
	"strconv"
)

type BgpGlobal struct {
	nclib.SrChanges `xml:"-"`

	Config *BgpGlobalConfig `xml:"config"`
}

type BgpGlobalProcessor interface {
	BgpGlobalConfigProcessor
}

func NewBgpGlobal() *BgpGlobal {
	return &BgpGlobal{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewBgpGlobalConfig(),
	}
}

func (b *BgpGlobal) String() string {
	return fmt.Sprintf("%s{%s} %s",
		OC_GLOBAL_KEY,
		b.Config,
		b.SrChanges,
	)
}

func (b *BgpGlobal) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_CONFIG_KEY:
		if err := b.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpGlobal(p BgpGlobalProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, global *BgpGlobal) error {
	configFunc := func() error {
		if global.GetChange(OC_CONFIG_KEY) {
			return ProcessBgpGlobalConfig(
				p.(BgpGlobalConfigProcessor),
				reverse,
				name,
				key,
				global.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

type BgpGlobalConfig struct {
	nclib.SrChanges `xml:"-"`

	As       uint32 `xml:"as"`
	RouterId string `xml:"router-id"`
}

type BgpGlobalConfigProcessor interface {
	BgpGlobalConfig(string, *NetworkInstanceProtocolKey, *BgpGlobalConfig) error
}

func NewBgpGlobalConfig() *BgpGlobalConfig {
	return &BgpGlobalConfig{
		SrChanges: nclib.NewSrChanges(),
		As:        0,
		RouterId:  "",
	}
}

func (b *BgpGlobalConfig) String() string {
	return fmt.Sprintf("%s{%s=%d, %s='%s'} %s",
		OC_GLOBAL_KEY,
		BGP_AS_KEY, b.As,
		BGP_ROUTERID_KEY, b.RouterId,
		b.SrChanges,
	)
}

func (b *BgpGlobalConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_AS_KEY:
		as, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return fmt.Errorf("Invalid AS. %s", value)
		}
		b.As = uint32(as)

	case BGP_ROUTERID_KEY:
		b.RouterId = value
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpGlobalConfig(p BgpGlobalProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, config *BgpGlobalConfig) error {
	globalFunc := func() error {
		return p.BgpGlobalConfig(name, key, config)
	}

	return nclib.CallFunctions(reverse, globalFunc)
}
