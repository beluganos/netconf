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
	"strconv"
)

type BgpZebra struct {
	nclib.SrChanges `xml:"-"`

	Config *BgpZebraConfig `xml:"config"`
}

type BgpZebraProcessor interface {
	BgpZebraConfigProcessor
}

func NewBgpZebra() *BgpZebra {
	return &BgpZebra{
		SrChanges: nclib.NewSrChanges(),
		Config:    NewBgpZebraConfig(),
	}
}

func (b *BgpZebra) String() string {
	return fmt.Sprintf("%s{%s} %s",
		BGP_ZEBRA_KEY,
		b.Config,
		b.SrChanges,
	)
}

func (b *BgpZebra) Put(nodes []*ncxml.XPathNode, value string) error {
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

func ProcessBgpZebra(p BgpZebraProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, zebra *BgpZebra) error {

	configFunc := func() error {
		if zebra.GetChange(OC_CONFIG_KEY) {
			return ProcessBgpZebraConfig(
				p.(BgpZebraConfigProcessor),
				reverse,
				name,
				key,
				zebra.Config,
			)
		}
		return nil
	}

	return nclib.CallFunctions(reverse, configFunc)
}

type BgpZebraConfig struct {
	nclib.SrChanges `xml:"-"`

	Enabled      bool                  `xml:"enabled"`
	Version      uint32                `xml:"version"`
	Url          string                `xml:"url"`
	RedistRoutes []InstallProtocolType `xml:"redistribute"`
}

type BgpZebraConfigProcessor interface {
	BgpZebraConfig(string, *NetworkInstanceProtocolKey, *BgpZebraConfig) error
}

func NewBgpZebraConfig() *BgpZebraConfig {
	return &BgpZebraConfig{
		SrChanges:    nclib.NewSrChanges(),
		Enabled:      false,
		Version:      0,
		Url:          "",
		RedistRoutes: []InstallProtocolType{},
	}
}

func (b *BgpZebraConfig) String() string {
	return fmt.Sprintf("%s{%s=%t, %s=%d, %s='%s', %s=%v} %s",
		OC_CONFIG_KEY,
		OC_ENABLED_KEY, b.Enabled,
		BGP_ZEBRA_VERSION_KEY, b.Version,
		BGP_ZEBRA_URL_KEY, b.Url,
		BGP_ZEBRA_REDISTROUTES_KEY, b.RedistRoutes,
		b.SrChanges,
	)
}

func (b *BgpZebraConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case OC_ENABLED_KEY:
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		b.Enabled = enabled

	case BGP_ZEBRA_VERSION_KEY:
		version, err := strconv.ParseUint(value, 0, 32)
		if err != nil {
			return err
		}
		b.Version = uint32(version)

	case BGP_ZEBRA_URL_KEY:
		b.Url = value

	case BGP_ZEBRA_REDISTROUTES_KEY:
		ptype, err := ParseInstallProtocolType(value)
		if err != nil {
			return err
		}
		b.RedistRoutes = append(b.RedistRoutes, ptype)
	}

	b.SetChange(nodes[0].Name)
	return nil
}

func ProcessBgpZebraConfig(p BgpZebraConfigProcessor, reverse bool, name string, key *NetworkInstanceProtocolKey, config *BgpZebraConfig) error {
	configFunc := func() error {
		return p.BgpZebraConfig(name, key, config)
	}

	return nclib.CallFunctions(reverse, configFunc)
}
