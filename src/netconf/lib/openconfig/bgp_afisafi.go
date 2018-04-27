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
	"encoding/xml"
	"fmt"
	"netconf/lib"
	"netconf/lib/xml"
)

//
// afi-safis
//
type BgpAfiSafis map[string]*BgpAfiSafi

func NewBgpAfiSafis() BgpAfiSafis {
	return BgpAfiSafis{}

}

func (b BgpAfiSafis) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	name, ok := nodes[0].Attrs[BGP_AFISAFI_NAME_KEY]
	if !ok {
		return fmt.Errorf("%s@%s not found. %s", BGP_AFISAFI_KEY, BGP_AFISAFI_NAME_KEY, nodes[0])
	}

	afisafi, ok := b[name]
	if !ok {
		afisafi = NewBgpAfiSafi(name)
		b[name] = afisafi
	}

	return afisafi.Put(nodes[1:], value)
}

func (b BgpAfiSafis) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = BGP_AFISAFIS_KEY
	e.EncodeToken(start)

	for _, a := range b {
		err := e.EncodeElement(a, xml.StartElement{Name: xml.Name{Local: BGP_AFISAFI_KEY}})
		if err != nil {
			return err
		}
	}

	return e.EncodeToken(start.End())
}

//
// afi-safis/afi-safi[afi-safi-name]
//
type BgpAfiSafi struct {
	nclib.SrChanges `xml:"-"`

	AfiSafiName string            `xml:"afi-safi-name"`
	Config      *BgpAfiSafiConfig `xml:"config"`
}

func NewBgpAfiSafi(afiSafiName string) *BgpAfiSafi {
	return &BgpAfiSafi{
		SrChanges:   nclib.NewSrChanges(),
		AfiSafiName: afiSafiName,
		Config:      NewBgpAfiSafiConfig(),
	}
}

func (b *BgpAfiSafi) String() string {
	return fmt.Sprintf("%s{%s=%s, %s} %s",
		BGP_AFISAFI_KEY,
		BGP_AFISAFI_NAME_KEY, b.AfiSafiName,
		b.Config,
		b.SrChanges,
	)
}

func (b *BgpAfiSafi) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_AFISAFI_NAME_KEY:
		// b.AfiSafiName = value // set by NewBgpAfiSafi

	case OC_CONFIG_KEY:
		if err := b.Config.Put(nodes[1:], value); err != nil {
			return err
		}
	}

	b.SetChange(nodes[0].Name)
	return nil
}

//
// afi-safis/afi-safi[afi-safi-name]/config
//
type BgpAfiSafiConfig struct {
	nclib.SrChanges `xml:"-"`

	AfiSafiName BgpAfiSafiType `xml:"afi-safi-name"`
}

func NewBgpAfiSafiConfig() *BgpAfiSafiConfig {
	return &BgpAfiSafiConfig{
		SrChanges:   nclib.NewSrChanges(),
		AfiSafiName: BGP_AFI_SAFI_TYPE,
	}
}

func (b *BgpAfiSafiConfig) String() string {
	return fmt.Sprintf("%s{%s=%s} %s",
		OC_CONFIG_KEY,
		BGP_AFISAFI_NAME_KEY, b.AfiSafiName,
		b.SrChanges,
	)
}

func (b *BgpAfiSafiConfig) Put(nodes []*ncxml.XPathNode, value string) error {
	if len(nodes) == 0 {
		return nil
	}

	switch nodes[0].Name {
	case BGP_AFISAFI_NAME_KEY:
		name, err := ParseBgpAfiSafiType(value)
		if err != nil {
			return err
		}
		b.AfiSafiName = name
	}

	b.SetChange(nodes[0].Name)
	return nil
}
