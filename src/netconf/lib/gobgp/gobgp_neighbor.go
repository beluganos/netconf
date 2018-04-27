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

//
// [neighbors]
//
type Neighbor Entries

func NewNeighbor(i interface{}) Neighbor {
	return Neighbor(NewEntries(i))
}

func NewNeighbors(i interface{}) []Neighbor {
	neighs := []Neighbor{}
	switch i.(type) {
	case nil:
	default:
		for _, neigh := range i.([]interface{}) {
			neighs = append(neighs, NewNeighbor(neigh))
		}
	}
	return neighs
}

func RawNeighbors(neighs []Neighbor) interface{} {
	list := make([]interface{}, len(neighs))
	for index, neigh := range neighs {
		list[index] = Entries(neigh).Raw()
	}
	return list
}

func SelectNeighbor(i interface{}, addr string) (Neighbor, int) {
	switch i.(type) {
	case nil:
	default:
		for index, n := range i.([]interface{}) {
			if neigh := NewNeighbor(n); neigh.Config().NeighborAddress() == addr {
				return neigh, index
			}
		}
	}

	return nil, -1
}

func (n Neighbor) Config() NeighborConfig {
	return NewNeighborConfig(getValue(n, "config"))
}

func (n Neighbor) SetConfig(v NeighborConfig) {
	n["config"] = v
}

func (n Neighbor) Transport() Transport {
	return NewTransport(getValue(n, "transport"))
}

func (n Neighbor) SetTransport(v Transport) {
	n["transport"] = v
}

func (n Neighbor) AfiSafis() []AfiSafi {
	return NewAfiSafis(getValue(n, "afi-safis"))
}

func (n Neighbor) SetAfisafis(afisafis []AfiSafi) {
	n["afi-safis"] = RawAfiSafis(afisafis)
}

func (n Neighbor) ApplyPolicy() ApplyPolicy {
	return NewApplyPolicy(getValue(n, "apply-policy"))
}

//
// [neighbors.config]
//
type NeighborConfig Entries

func NewNeighborConfig(i interface{}) NeighborConfig {
	return NeighborConfig(NewEntries(i))
}

func (c NeighborConfig) NeighborAddress() string {
	return convString(c, "neighbor-address")
}

func (c NeighborConfig) SetNeighborAddress(v string) {
	c["neighbor-address"] = v
}

func (c NeighborConfig) PeerAs() uint32 {
	return uint32(convUint(c, "peer-as"))
}

func (c NeighborConfig) SetPeerAs(v uint32) {
	c["peer-as"] = v
}

func (c NeighborConfig) LocalAs() uint32 {
	return uint32(convUint(c, "local-as"))
}

func (c NeighborConfig) SetLocalAs(v uint32) {
	c["local-as"] = v
}

//
// [neighbors.apply-policy]
//
type ApplyPolicy Entries

func NewApplyPolicy(i interface{}) ApplyPolicy {
	return ApplyPolicy(NewEntries(i))
}

func (a ApplyPolicy) Config() ApplyPolicyConfig {
	return NewApplyPolicyConfig(getValue(a, "config"))
}

func (a ApplyPolicy) SetConfig(c ApplyPolicyConfig) {
	a["config"] = c
}

//
// [neighbors.apply-policy.config]
//
type ApplyPolicyConfig Entries

func NewApplyPolicyConfig(i interface{}) ApplyPolicyConfig {
	return ApplyPolicyConfig(NewEntries(i))
}

func (a ApplyPolicyConfig) DefaultExportPolicy() string {
	return convString(a, "default-export-policy")
}

func (a ApplyPolicyConfig) SetDefaultExportPolicy(policyName string) {
	a["default-export-policy"] = policyName
}

func (a ApplyPolicyConfig) DefaultImportPolicy() string {
	return convString(a, "default-import-policy")
}

func (a ApplyPolicyConfig) SetDefaultImportPolicy(policyName string) {
	a["default-import-policy"] = policyName
}

func (a ApplyPolicyConfig) ExportPolicyList() []string {
	pols := convList(a, "export-policy-list")
	list := make([]string, len(pols))
	for index, pol := range pols {
		list[index] = pol.(string)
	}
	return list
}

func (a ApplyPolicyConfig) SetExportPolicyList(policyNames []string) {
	list := make([]interface{}, len(policyNames))
	for index, policyName := range policyNames {
		list[index] = policyName
	}

	a["export-policy-list"] = list
}

func (a ApplyPolicyConfig) ImportPolicyList() []string {
	pols := convList(a, "import-policy-list")
	list := make([]string, len(pols))
	for index, pol := range pols {
		list[index] = pol.(string)
	}
	return list
}

func (a ApplyPolicyConfig) SetImportPolicyList(policyNames []string) {
	list := make([]interface{}, len(policyNames))
	for index, policyName := range policyNames {
		list[index] = policyName
	}

	a["import-policy-list"] = list
}

func (a ApplyPolicyConfig) PolicyList() []string {
	return append(a.ImportPolicyList(), a.ExportPolicyList()...)
}
