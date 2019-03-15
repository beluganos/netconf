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

const NIConterinerDefaultMTU = 9000

type NIAnyHandler struct {
	*NICommands
	ev   srlib.SrNotifEvent
	oper srlib.SrChangeOper
	mtu  uint16
}

func (n *NIAnyHandler) SetOpt(key string, val interface{}) {
	switch key {
	case "mtu":
		if mtu, ok := val.(uint16); ok {
			n.mtu = mtu
		}
	case "dryrun":
		if dryRun, ok := val.(bool); ok {
			n.Cmds.DryRun = dryRun
		}
	}
}

func NewNIAnyHandler(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) NIChangeHandler {
	return newNIAnyHandler(ev, oper)
}

func newNIAnyHandler(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) *NIAnyHandler {
	return &NIAnyHandler{
		NICommands: NewNICommands(ev, oper),
		ev:         ev,
		oper:       oper,
		mtu:        NIConterinerDefaultMTU,
	}
}

func (h *NIAnyHandler) DoCmds() error {
	if err := h.Cmds.Do(); err != nil {
		log.Errorf("NI: DoCommand error. %s", err)
		h.Clear()
		return err
	}
	return nil
}

func (h *NIAnyHandler) Begin(name string, ni *openconfig.NetworkInstance) error {
	log.Debugf("NI/%s/%s/%s/BEGIN* %s", h.ev, h.oper, name, ni)
	return nil
}

func (h *NIAnyHandler) Commit() error {
	if h.NoCommit {
		log.Debugf("NI/%s/%s/COMMIT* SKIP.", h.ev, h.oper)
		return nil
	}

	log.Debugf("NI/%s/%s/COMMIT*", h.ev, h.oper)
	return h.Cmds.End()
}

func (h *NIAnyHandler) Rollback() {
	log.Debugf("NI/%s/%s/ROLLBACK*", h.ev, h.oper)
	h.Cmds.Undo()
}

func (h *NIAnyHandler) NetworkInstance(name string, ni *openconfig.NetworkInstance) error {
	log.Debugf("NI/%s/%s/%s* %s", h.ev, h.oper, name, ni)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceConfig(name string, config *openconfig.NetworkInstanceConfig) error {
	log.Debugf("NI/%s/%s/%s/CONF* %s", h.ev, h.oper, name, config)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceLoopback(name string, id string, lo *openconfig.NetworkInstanceLoopback) error {
	log.Debugf("NI/%s/%s/%s/%s* %s", h.ev, h.oper, name, id, lo)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceLoopbackConfig(name string, id string, config *openconfig.NetworkInstanceLoopbackConfig) error {
	log.Debugf("NI/%s/%s/%s/%s* %s", h.ev, h.oper, name, id, config)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceLoopbackAddr(name string, id string, index string, addr *openconfig.NetworkInstanceLoopbackAddr) error {
	log.Debugf("NI/%s/%s/%s/%s/%s* %s", h.ev, h.oper, name, id, index, addr)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceLoopbackAddrConfig(name string, id string, index string, config *openconfig.NetworkInstanceLoopbackAddrConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/%s* %s", h.ev, h.oper, name, id, index, config)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceInterface(name string, id string, iface *openconfig.NetworkInstanceInterface) error {
	log.Debugf("NI/%s/%s/%s/%s* %s", h.ev, h.oper, name, id, iface)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceInterfaceConfig(name string, id string, config *openconfig.NetworkInstanceInterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/CONF* %s", h.ev, h.oper, name, id, config)
	return nil
}

func (h *NIAnyHandler) MplsConfig(name string, config *openconfig.MplsConfig) error {
	log.Debugf("NI/%s/%s/%s/MPLS/CONF* %s", h.ev, h.oper, name, config)
	return nil
}

func (h *NIAnyHandler) MplsInterfaceAttr(name string, ifaceId string, attr *openconfig.MplsInterfaceAttr) error {
	log.Debugf("NI/%s/%s/%s/MPLS/%s* %s", h.ev, h.oper, name, ifaceId, attr)
	return nil
}

func (h *NIAnyHandler) MplsInterfaceAttrConfig(name string, ifaceId string, config *openconfig.MplsInterfaceAttrConfig) error {
	log.Debugf("NI/%s/%s/%s/MPLS/%s/CONF* %s", h.ev, h.oper, name, ifaceId, config)
	return nil
}

func (h *NIAnyHandler) MplsInterfaceAttrRefConfig(name string, ifaceId string, config *openconfig.InterfaceRefConfig) error {
	log.Debugf("NI/%s/%s/%s/MPLS/%s/REF* %s", h.ev, h.oper, name, ifaceId, config)
	return nil
}

func (h *NIAnyHandler) MplsLdpConfig(name string, config *openconfig.MplsLdpConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/CONF* %s", h.ev, h.oper, name, config)
	return nil
}

func (h *NIAnyHandler) MplsLdpAddressFamilyV4Config(name string, config *openconfig.MplsLdpAddressFamilyV4Config) error {
	log.Debugf("NI/%s/%s/%s/LDP/ADDR* %s", h.ev, h.oper, name, config)
	return nil
}

func (h *NIAnyHandler) MplsLdpDiscovInterfacesConfig(name string, config *openconfig.MplsLdpDiscovInterfacesConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/IFCONF* %s", h.ev, h.oper, name, config)
	return nil
}

func (h *NIAnyHandler) MplsLdpInterface(name string, ifaceId string, iface *openconfig.MplsLdpInterface) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/%s* %s", h.ev, h.oper, name, ifaceId, iface)
	return nil
}

func (h *NIAnyHandler) MplsLdpInterfaceConfig(name string, ifaceId string, config *openconfig.MplsLdpInterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISC/%s/CONF* %s", h.ev, h.oper, name, ifaceId, config)
	return nil
}

func (h *NIAnyHandler) MplsLdpInterfaceRefConfig(name string, ifaceId string, config *openconfig.InterfaceRefConfig) error {
	log.Debugf("NI/%s/%s/%s/LDP/DISCC/%s/REF* %s", h.ev, h.oper, name, ifaceId, config)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceProtocol(name string, key *openconfig.NetworkInstanceProtocolKey, proto *openconfig.NetworkInstanceProtocol) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s* %s", h.ev, h.oper, name, key, proto)
	return nil
}

func (h *NIAnyHandler) NetworkInstanceProtocolConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.NetworkInstanceProtocolConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/CONF* %s", h.ev, h.oper, name, key, config)
	return nil
}

func (h *NIAnyHandler) Ospfv2(name string, key *openconfig.NetworkInstanceProtocolKey, ospf *openconfig.Ospfv2) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s* %s", h.ev, h.oper, name, key, ospf)
	return nil
}

func (h *NIAnyHandler) Ospfv2GlobalConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.Ospfv2GlobalConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/CONF* %s", h.ev, h.oper, name, key, config)
	return nil
}

func (h *NIAnyHandler) Ospfv2Area(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, area *openconfig.Ospfv2Area) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s* %s", h.ev, h.oper, name, key, areaId, area)
	return nil
}

func (h *NIAnyHandler) Ospfv2AreaConfig(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, config *openconfig.Ospfv2AreaConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/CONF* %s", h.ev, h.oper, name, key, areaId, config)
	return nil
}

func (h *NIAnyHandler) Ospfv2Interface(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, iface *openconfig.Ospfv2Interface) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s* %s", h.ev, h.oper, name, key, areaId, ifaceId, iface)
	return nil
}

func (h *NIAnyHandler) Ospfv2InterfaceConfig(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, config *openconfig.Ospfv2InterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF* %s", h.ev, h.oper, name, key, areaId, ifaceId, config)
	return nil
}

func (h *NIAnyHandler) Ospfv2InterfaceRefConfig(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, config *openconfig.InterfaceRefConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/IFREF* %s", h.ev, h.oper, name, key, areaId, ifaceId, config)
	return nil
}

func (h *NIAnyHandler) Ospfv2InterfaceTimers(name string, key *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, timers *openconfig.Ospfv2InterfaceTimers) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/TIMERS* %s", h.ev, h.oper, name, key, areaId, ifaceId, timers)
	return nil
}

func (h *NIAnyHandler) Ospfv3(name string, nikey *openconfig.NetworkInstanceProtocolKey, ospf *openconfig.Ospfv3) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s* %s", h.ev, h.oper, name, nikey, ospf)
	return nil
}

func (h *NIAnyHandler) Ospfv3GlobalConfig(name string, nikey *openconfig.NetworkInstanceProtocolKey, config *openconfig.Ospfv3GlobalConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/CONF* %s", h.ev, h.oper, name, nikey, config)
	return nil
}

func (h *NIAnyHandler) Ospfv3Area(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, area *openconfig.Ospfv3Area) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s* %s", h.ev, h.oper, name, nikey, areaId, area)
	return nil
}

func (h *NIAnyHandler) Ospfv3AreaConfig(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, config *openconfig.Ospfv3AreaConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/CONF* %s", h.ev, h.oper, name, nikey, areaId, config)
	return nil
}

func (h *NIAnyHandler) Ospfv3AreaRange(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, rngkey *openconfig.Ospfv3AreaRangeKey, rng *openconfig.Ospfv3AreaRange) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s* %s", h.ev, h.oper, name, nikey, areaId, rngkey, rng)
	return nil
}

func (h *NIAnyHandler) Ospfv3AreaRangeConfig(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, rngkey *openconfig.Ospfv3AreaRangeKey, config *openconfig.Ospfv3AreaRangeConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF* %s", h.ev, h.oper, name, nikey, areaId, rngkey, config)
	return nil
}

func (h *NIAnyHandler) Ospfv3Interface(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, iface *openconfig.Ospfv3Interface) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s* %s", h.ev, h.oper, name, nikey, areaId, ifaceId, iface)
	return nil
}

func (h *NIAnyHandler) Ospfv3InterfaceConfig(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, config *openconfig.Ospfv3InterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF* %s", h.ev, h.oper, name, nikey, areaId, ifaceId, config)

	return nil
}

func (h *NIAnyHandler) Ospfv3InterfaceRefConfig(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, config *openconfig.InterfaceRefConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/IFREF* %s", h.ev, h.oper, name, nikey, areaId, ifaceId, config)
	return nil
}

func (h *NIAnyHandler) Ospfv3InterfaceTimers(name string, nikey *openconfig.NetworkInstanceProtocolKey, areaId string, ifaceId string, timers *openconfig.Ospfv3InterfaceTimers) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/TIMERS* %s", h.ev, h.oper, name, nikey, areaId, ifaceId, timers)
	return nil
}

func (h *NIAnyHandler) StaticRoute(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, route *openconfig.StaticRoute) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s* %s", h.ev, h.oper, name, prkey, rtkey, route)
	return nil
}

func (h *NIAnyHandler) StaticRouteConfig(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, config *openconfig.StaticRouteConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/CONF* %s", h.ev, h.oper, name, prkey, rtkey, config)
	return nil
}

func (h *NIAnyHandler) StaticRouteNexthop(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, nexthop *openconfig.StaticRouteNexthop) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s* %s", h.ev, h.oper, name, prkey, rtkey, index, nexthop)
	return nil
}

func (h *NIAnyHandler) StaticRouteNexthopConfig(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, config *openconfig.StaticRouteNexthopConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF* %s", h.ev, h.oper, name, prkey, rtkey, index, config)
	return nil
}

func (h *NIAnyHandler) StaticRouteNexthopIfaceRefConfig(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, config *openconfig.InterfaceRefConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/IFREF* %s", h.ev, h.oper, name, prkey, rtkey, index, config)
	return nil
}

func (h *NIAnyHandler) Bgp(name string, key *openconfig.NetworkInstanceProtocolKey, bgp *openconfig.Bgp) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s* %s", h.ev, h.oper, name, key, bgp)
	return nil
}

func (h *NIAnyHandler) BgpGlobalConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.BgpGlobalConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/GLOBAL/CONF* %s", h.ev, h.oper, name, key, config)
	return nil
}

func (h *NIAnyHandler) BgpZebraConfig(name string, key *openconfig.NetworkInstanceProtocolKey, config *openconfig.BgpZebraConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/ZEBRA/CONF* %s", h.ev, h.oper, name, key, config)
	return nil
}

func (h *NIAnyHandler) BgpNeighbor(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, neigh *openconfig.BgpNeighbor) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s* %s", h.ev, h.oper, name, key, addr, neigh)
	return nil
}

func (h *NIAnyHandler) BgpNeighborConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.BgpNeighborConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/CONF* %s", h.ev, h.oper, name, key, addr, config)
	return nil
}

func (h *NIAnyHandler) BgpNeighborTimersConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.BgpNeighborTimersConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/TIMERS* %s", h.ev, h.oper, name, key, addr, config)
	return nil
}

func (h *NIAnyHandler) BgpNeighborTransportConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.BgpNeighborTransportConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/TRANS* %s", h.ev, h.oper, name, key, addr, config)
	return nil
}

func (h *NIAnyHandler) BgpNeighborApplyPolicyConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.PolicyApplyConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/APPLYPOL* %s", h.ev, h.oper, name, key, addr, config)
	return nil
}

func (h *NIAnyHandler) BgpNeighborAfiSafi(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, AfiSafiName string, afiSafi *openconfig.BgpAfiSafi) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s* %s", h.ev, h.oper, name, key, addr, AfiSafiName, afiSafi)
	return nil
}

func (h *NIAnyHandler) BgpNeighborAfiSafiConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, AfiSafiName string, config *openconfig.BgpAfiSafiConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF* %s", h.ev, h.oper, name, key, addr, AfiSafiName, config)
	return nil
}
