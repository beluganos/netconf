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

package ncm

import (
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"

	log "github.com/sirupsen/logrus"
)

type NIModifyApplyHandler struct {
	*NIAnyHandler
}

func NewNIModifyApplyHandler(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) NIChangeHandler {
	return &NIModifyApplyHandler{
		NIAnyHandler: newNIAnyHandler(ev, oper),
	}
}

func (h *NIModifyApplyHandler) Begin(name string, ni *openconfig.NetworkInstance) error {
	log.Debugf("NI/%s/%s/%s/BEGIN: %s", h.ev, h.oper, name, ni)

	h.Clear()
	if err := openconfig.ProcessNetworkInstance(h, false, name, ni); err != nil {
		h.Clear()
		return err
	}

	return h.DoCmds()
}

func (h *NIModifyApplyHandler) Ospfv2GlobalConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.Ospfv2GlobalConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/CONF: %s", h.ev, h.oper, name, key, config)

	if config.GetChange(openconfig.OSPFV2_ROUTERID_KEY) {
		AddNIOspfRouterCmd(h, name, "router-id", config.RouterId, true)
	}

	return nil
}

func (h *NIModifyApplyHandler) Ospfv2InterfaceConfig(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, config *openconfig.Ospfv2InterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF: %s", h.ev, h.oper, name, key, areaId, ifaceId, config)

	if config.GetChange(openconfig.OSPFV2_METRIC_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf cost", config.Metric, true)
	}

	if config.GetChange(openconfig.OSPFV2_PASSIVE_KEY) {
		AddNIOspfRouterCmd(h, name, "passive-interface", ifaceId, config.Passive)
	}

	if config.GetChange(openconfig.OSPFV2_PRIORITY_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf priority", config.Priority, true)
	}

	if config.GetChange(openconfig.OSPFV2_NETWORK_TYPE_KEY) {
		n, _ := getOspfNetworkType(config)
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf network", n, true)
	}

	return nil
}

func (h *NIModifyApplyHandler) Ospfv2InterfaceTimers(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, timers *openconfig.Ospfv2InterfaceTimers) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/TIMERS: %s", h.ev, h.oper, name, key, areaId, ifaceId, timers)

	if timers.GetChange(openconfig.OSPFV2_DEAD_INTERVAL_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf dead-interval", timers.DeadInterval, true)
	}

	if timers.GetChange(openconfig.OSPFV2_HELLO_INTERVAL_KEY) {
		AddNIVtyInterfaceCmd(h, name, ifaceId, "ip ospf hello-interval", timers.HelloInterval, true)
	}

	return nil
}

func (h *NIModifyApplyHandler) MplsLdpConfig(name string, config *openconfig.MplsLdpConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/CONF: %s", h.ev, h.oper, name, config)

	if config.GetChange(openconfig.MPLS_LDP_ROUTERID_KEY) {
		AddNIMplsLdpCmd(h, name, "router-id", config.LsrId, true)
	}

	return nil
}

func (h *NIModifyApplyHandler) MplsLdpAddressFamilyV4Config(name string, config *openconfig.MplsLdpAddressFamilyV4Config) error {
	log.Debugf("NI/%s/%s/%s/LDP/ADDR: %s", h.ev, h.oper, name, config)

	if config.GetChange(openconfig.MPLS_LDP_AF_TARNSADDR_KEY) {
		AddNIMplsLdpIPv4Cmd(h, name, "discovery transport-address", config.TransportAddr, true)
	}

	if config.GetChange(openconfig.MPLS_LDP_AF_SESSION_HOLDTIME_KEY) {
		AddNIMplsLdpIPv4Cmd(h, name, "session holdtime", config.SessionHoldTime, true)
	}

	explicitNll := config.LabelPolicy.Advertise.EngressExplicitNull
	if explicitNll.GetChange(openconfig.MPLS_LDP_LABELPOLICY_ENABLE_KEY) {
		AddNIMplsLdpIPv4Cmd(h, name, "label local advertise explicit-null", "", explicitNll.Enable)
	}

	return nil
}

func (h *NIModifyApplyHandler) MplsLdpDiscovInterfacesConfig(name string, config *openconfig.MplsLdpDiscovInterfacesConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/IFCONF: %s", h.ev, h.oper, name, config)

	if config.GetChange(openconfig.MPLS_LDP_HELLO_HOLDTIME) {
		AddNIMplsLdpCmd(h, name, "discovery hello holdtime", config.HelloHoldTime, true)
	}

	if config.GetChange(openconfig.MPLS_LDP_HELLO_INTERVAL) {
		AddNIMplsLdpCmd(h, name, "discovery hello interval", config.HelloInterval, true)
	}

	return nil
}

func (h *NIModifyApplyHandler) MplsLdpInterfaceConfig(name string, ifaceId string, config *openconfig.MplsLdpInterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/%s/CONF: %s", h.ev, h.oper, name, ifaceId, config)

	if config.GetChange(openconfig.MPLS_LDP_HELLO_HOLDTIME) {
		AddNIMplsLdpIPv4IfaceCmd(h, name, ifaceId, "discovery hello holdtime", config.HelloHoldTime, true)
	}

	if config.GetChange(openconfig.MPLS_LDP_HELLO_INTERVAL) {
		AddNIMplsLdpIPv4IfaceCmd(h, name, ifaceId, "discovery hello interval", config.HelloInterval, true)
	}

	return nil
}
