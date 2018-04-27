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

package srocgobgp

import (
	"bytes"
	"container/list"
	"fmt"
	"netconf/lib/openconfig"
)

type ConfigProcessor struct {
	items *list.List
}

func NewConfigProcessor() *ConfigProcessor {
	return &ConfigProcessor{
		items: list.New(),
	}
}

func (p *ConfigProcessor) Iterate(f func(string) error) error {
	for e := p.items.Front(); e != nil; e = e.Next() {
		if err := f(e.Value.(string)); err != nil {
			return err
		}
	}
	return nil
}

func (b *ConfigProcessor) Len() int {
	return b.items.Len()
}

func (p *ConfigProcessor) Items() []string {
	items := []string{}
	p.Iterate(func(s string) error {
		items = append(items, s)
		return nil
	})
	return items
}

func (p *ConfigProcessor) Bytes() *bytes.Buffer {
	var b bytes.Buffer
	p.Iterate(func(s string) error {
		b.WriteString(s)
		b.WriteString("\n")
		return nil
	})
	return &b
}

func (p *ConfigProcessor) Clear() {
	p.items = list.New()
}

func (p *ConfigProcessor) addList(name string) {
	p.items.PushBack(fmt.Sprintf("[[%s]]", name))
}

func (p *ConfigProcessor) addNode(name string) {
	p.items.PushBack(fmt.Sprintf("[%s]", name))
}

func (p *ConfigProcessor) addItem(name string, value interface{}) {
	p.items.PushBack(fmt.Sprintf("%s = %v", name, value))
}

func (p *ConfigProcessor) Bgp(name string, key *openconfig.NetworkInstanceProtocolKey, bgp *openconfig.Bgp) error {
	return nil
}

func (p *ConfigProcessor) BgpGlobalConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.BgpGlobalConfig) error {

	p.addNode("global.config")

	if config.GetChange(openconfig.BGP_AS_KEY) {
		p.addItem("as", config.As)
	}

	if config.GetChange(openconfig.BGP_ROUTERID_KEY) {
		p.addItem("router-id", QString(config.RouterId))
	}

	return nil
}

func (p *ConfigProcessor) BgpZebraConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.BgpZebraConfig) error {

	p.addNode("zebra.config")

	if config.GetChange(openconfig.OC_ENABLED_KEY) {
		p.addItem("enabled", config.Enabled)
	}

	if config.GetChange(openconfig.BGP_ZEBRA_VERSION_KEY) {
		p.addItem("version", config.Version)
	}

	if config.GetChange(openconfig.BGP_ZEBRA_URL_KEY) {
		p.addItem("url", QString(config.Url))
	}

	if config.GetChange(openconfig.BGP_ZEBRA_REDISTROUTES_KEY) {
		p.addItem("redistribute-route-type-list", InstallProtocolTypes(config.RedistRoutes))
	}

	return nil
}

func (p *ConfigProcessor) BgpNeighbor(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, neighbor *openconfig.BgpNeighbor) error {
	p.addList("neighbors")
	return nil
}

func (p *ConfigProcessor) BgpNeighborConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.BgpNeighborConfig) error {
	p.addNode("neighbors.config")

	if config.GetChange(openconfig.BGP_NEIGHBOR_ADDR_KEY) {
		p.addItem("neighbor-address", QString(config.Address))
	}

	if config.GetChange(openconfig.BGP_PEERAS_KEY) {
		p.addItem("peer-as", config.PeerAs)
	}

	if config.GetChange(openconfig.BGP_LOCALAS_KEY) {
		p.addItem("local-as", config.LocalAs)
	}

	return nil
}

func (p *ConfigProcessor) BgpNeighborTimersConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.BgpNeighborTimersConfig) error {

	p.addNode("neighbors.timers.config")

	if config.GetChange(openconfig.BGP_HOLDTIME_KEY) {
		p.addItem("hold-time", config.HoldTime)
	}

	if config.GetChange(openconfig.BGP_KEEPALIVE_INTERVAL_KEY) {
		p.addItem("keepalive-interval", config.KeepAlive)
	}

	return nil
}

func (p *ConfigProcessor) BgpNeighborTransportConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.BgpNeighborTransportConfig) error {

	p.addNode("neighbors.transport.config")

	if config.GetChange(openconfig.BGP_LOCAL_ADDR_KEY) {
		p.addItem("local-address", QString(config.LocalAddr))
	}

	return nil
}

func (p *ConfigProcessor) BgpNeighborAfiSafi(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, afiSafiName string, afisafi *openconfig.BgpAfiSafi) error {

	p.addList("neighbors.afi-safis")

	return nil
}

func (p *ConfigProcessor) BgpNeighborAfiSafiConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, afiSafiName string, config *openconfig.BgpAfiSafiConfig) error {

	p.addNode("neighbors.afi-safis.config")

	if config.GetChange(openconfig.BGP_AFISAFI_NAME_KEY) {
		p.addItem("afi-safi-name", QString(BgpAfiSafiType(config.AfiSafiName)))
	}

	return nil
}

func (p *ConfigProcessor) BgpNeighborApplyPolicyConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.PolicyApplyConfig) error {

	p.addNode("neighbors.apply-policy.config")

	if config.GetChange(openconfig.POLICYAPPLY_IMPORT_KEY) {
		p.addItem("import-policy-list", QStringList(config.ImportPolicy))
	}

	if config.GetChange(openconfig.POLICYAPPLY_EXPORT_KEY) {
		p.addItem("export-policy-list", QStringList(config.ExportPolicy))
	}

	if config.GetChange(openconfig.POLICYAPPLY_IMPORT_DEF_KEY) {
		p.addItem("default-import-policy", QString(PolicyDefaultType(config.ImportDefault)))
	}

	if config.GetChange(openconfig.POLICYAPPLY_EXPORT_DEF_KEY) {
		p.addItem("default-export-policy", QString(PolicyDefaultType(config.ExportDefault)))
	}

	return nil
}

func (p *ConfigProcessor) PolicyDefinition(polName string, pol *openconfig.PolicyDefinition) error {
	p.addList("policy-definitions")
	p.addItem("name", QString(polName))
	return nil
}

func (p *ConfigProcessor) PolicyDefinitionConfig(polName string, config *openconfig.PolicyDefinitionConfig) error {
	return nil
}

func (p *ConfigProcessor) PolicyStatement(polName string, stmtName string, stmt *openconfig.PolicyStatement) error {
	p.addList("policy-definitions.statements")
	p.addItem("name", QString(stmtName))
	return nil
}

func (p *ConfigProcessor) PolicyStatementConfig(polName string, stmtName string, config *openconfig.PolicyStatementConfig) error {
	return nil
}

func (p *ConfigProcessor) PolicyStatementActionsConfig(polName string, stmtName string, config *openconfig.PolicyStatementActionsConfig) error {
	p.addNode("policy-definitions.statements.actions")

	if config.GetChange(openconfig.POLICYDEF_ACTS_POLRESULT_KEY) {
		p.addItem("route-disposition", QString(PolicyResultType(config.PolicyResult)))
	}

	return nil
}

func (p *ConfigProcessor) PolicyBgpActionsConfig(polName string, stmtName string, config *openconfig.PolicyBgpActionsConfig) error {
	p.addNode("policy-definitions.statements.actions.bgp-actions")

	if config.GetChange(openconfig.BGP_ACTIONS_SET_NEXTHOP_KEY) {
		nhIP, _, err := config.SetNexthop.Values()
		if err != nil {
			return err
		}

		nh := func() string {
			if nhIP != nil {
				return nhIP.String()
			} else {
				return "self"
			}
		}()

		p.addItem("set-next-hop", QString(nh))
	}

	if config.GetChange(openconfig.BGP_ACTIONS_SET_LOCALPREF_KEY) {
		p.addItem("set-local-pref", config.SetLocalPref)
	}

	return nil
}

func (p *ConfigProcessor) PolicyNeighborSet(polName string, neighSet *openconfig.PolicyNeighborSet) error {
	return nil
}

func (p *ConfigProcessor) PolicyNeighborSetConfig(polName string, config *openconfig.PolicyNeighborSetConfig) error {
	return nil
}

func (p *ConfigProcessor) PolicyPrefixSet(polName string, pfxSet *openconfig.PolicyPrefixSet) error {
	return nil
}

func (p *ConfigProcessor) PolicyPrefixSetConfig(polName string, config *openconfig.PolicyPrefixSetConfig) error {
	return nil
}

func (p *ConfigProcessor) PolicyPrefixSetPrefix(polName string, pfxKey *openconfig.PolicyPrefixSetPrefixKey, prefix *openconfig.PolicyPrefixSetPrefix) error {
	return nil
}

func (p *ConfigProcessor) PolicyPrefixSetPrefixConfig(polName string, pfxKey *openconfig.PolicyPrefixSetPrefixKey, config *openconfig.PolicyPrefixSetPrefixConfig) error {
	return nil
}

func (p *ConfigProcessor) PolicyTagSet(polName string, tagSet *openconfig.PolicyTagSet) error {
	return nil
}

func (p *ConfigProcessor) PolicyTagSetConfig(polName string, config *openconfig.PolicyTagSetConfig) error {
	return nil
}
