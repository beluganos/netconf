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

import (
	"io"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

func NewConfig(v *viper.Viper) *Config {
	if v == nil {
		v = viper.New()
	}
	return &Config{
		Viper: v,
	}
}

func NewConfigFile(confPath string, confType string) *Config {
	v := viper.New()
	v.SetConfigFile(confPath)
	v.SetConfigType(confType)
	return NewConfig(v)
}

func ReadConfigFile(confPath string, confType string) (*Config, error) {
	c := NewConfigFile(confPath, confType)
	if err := c.ReadInConfig(); err != nil {
		return nil, err
	}

	return c, nil
}

func ReadConfig(r io.Reader, confType string) (*Config, error) {
	v := viper.New()
	v.SetConfigType(confType)
	if err := v.ReadConfig(r); err != nil {
		return nil, err
	}

	return NewConfig(v), nil
}

func (c *Config) HasGlobal() bool {
	return c.InConfig("global")
}

func (c *Config) Global() Global {
	return NewGlobal(c.Get("global"))
}

func (c *Config) SetGlobal(g Global) {
	c.Set("global", g)
}

func (c *Config) HasZebra() bool {
	return c.InConfig("zebra")
}

func (c *Config) Zebra() Zebra {
	return NewZebra(c.Get("zebra"))
}

func (c *Config) SetZebra(z Zebra) {
	c.Set("zebra", z)
}

func (c *Config) Neighbor(addr string) (Neighbor, int) {
	return SelectNeighbor(c.Get("neighbors"), addr)
}

func (c *Config) HasNeighbors() bool {
	return c.InConfig("neighbors")
}

func (c *Config) Neighbors() []Neighbor {
	return NewNeighbors(c.Get("neighbors"))
}

func (c *Config) SetNeighbors(neighbors []Neighbor) {
	c.Set("neighbors", RawNeighbors(neighbors))
}

func (c *Config) PolicyDefinition(name string) (PolicyDefinition, int) {
	return SelectPolicyDefinition(c.Get("policy-definitions"), name)
}

func (c *Config) HasPolicyDefinitions() bool {
	return c.InConfig("policy-definitions")
}

func (c *Config) PolicyDefinitions() []PolicyDefinition {
	return NewPolicyDefinitions(c.Get("policy-definitions"))
}

func (c *Config) SetPolicyDefinitions(defs []PolicyDefinition) {
	c.Set("policy-definitions", RawPolicyDefinitions(defs))
}

func (c *Config) DeleteUnusedPolicyDefinition() {
	polNames := make(map[string]struct{})
	for _, neigh := range c.Neighbors() {
		for _, name := range neigh.ApplyPolicy().Config().PolicyList() {
			polNames[name] = struct{}{}
		}
	}

	pols := []PolicyDefinition{}
	for _, pol := range c.PolicyDefinitions() {
		if _, ok := polNames[pol.Name()]; ok {
			pols = append(pols, pol)
		}
	}
	c.SetPolicyDefinitions(pols)
}

func (c *Config) DeleteDuplicatePolicyDefinition() {
	polNames := make(map[string]struct{})
	pols := []PolicyDefinition{}
	for _, pol := range c.PolicyDefinitions() {
		if _, ok := polNames[pol.Name()]; !ok {
			pols = append(pols, pol)
			polNames[pol.Name()] = struct{}{}
		}
	}
	c.SetPolicyDefinitions(pols)
}

func (c *Config) Merge(src *Config) {
	if src.HasGlobal() {
		c.SetGlobal(src.Global())
	}

	if src.HasZebra() {
		c.SetZebra(src.Zebra())
	}

	neighs := c.Neighbors()
	for _, sNeigh := range src.Neighbors() {
		addr := sNeigh.Config().NeighborAddress()
		if _, index := c.Neighbor(addr); index < 0 {
			neighs = append(neighs, sNeigh)
		} else {
			neighs[index] = sNeigh
		}
	}
	c.SetNeighbors(neighs)

	pols := c.PolicyDefinitions()
	for _, sPol := range src.PolicyDefinitions() {
		name := sPol.Name()
		if _, index := c.PolicyDefinition(name); index < 0 {
			pols = append(pols, sPol)
		} else {
			pols[index] = sPol
		}
	}
	c.SetPolicyDefinitions(pols)
	c.DeleteDuplicatePolicyDefinition()
}

func (c *Config) Delete(src *Config) {
	if src.HasGlobal() {
		c.SetGlobal(NewGlobal(nil))
	}

	if src.HasZebra() {
		c.SetZebra(NewZebra(nil))
	}

	neighs := []Neighbor{}
	for _, neigh := range c.Neighbors() {
		addr := neigh.Config().NeighborAddress()
		if _, index := src.Neighbor(addr); index < 0 {
			neighs = append(neighs, neigh)
		}
	}
	c.SetNeighbors(neighs)
	c.DeleteUnusedPolicyDefinition()
}
