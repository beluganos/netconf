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
	"fmt"
	"netconf/lib/openconfig"
	srlib "netconf/lib/sysrepo"

	log "github.com/sirupsen/logrus"
)

type NIDeleteVerifyHandler struct {
	*NIAnyHandler
}

func NewNIDeleteVerifyHandler(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) NIChangeHandler {
	return &NIDeleteVerifyHandler{
		NIAnyHandler: newNIAnyHandler(ev, oper),
	}
}

func (h *NIDeleteVerifyHandler) Begin(name string, ni *openconfig.NetworkInstance) error {
	log.Debugf("NI/%s/%s/%s/BEGIN; %s", h.ev, h.oper, name, ni)
	h.Clear()
	return openconfig.ProcessNetworkInstance(h, false, name, ni)
}

func (h *NIDeleteVerifyHandler) NetworkInstanceLoopbackAddrConfig(name string, id string, index string, config *openconfig.NetworkInstanceLoopbackAddrConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, index, config)

	if err := VerifyNILoopbackAddrConfig(config); err != nil {
		log.Errorf("NI/%s/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, index, err)
		return err
	}

	return nil
}

//
// /network-instances/network-instance[name]/interfaces/interface[id]
//
func (h *NIDeleteVerifyHandler) NetworkInstanceInterface(name string, id string, iface *openconfig.NetworkInstanceInterface) error {
	log.Debugf("NI/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, iface)

	if err := VerifyNIInterface(iface); err != nil {
		log.Errorf("NI/%s/%s/%s/%s: %s", h.ev, h.oper, name, id, err)
		return err
	}

	log.Debugf("NI/%s/%s/%s/%s: VERIFY OK.", h.ev, h.oper, name, id)
	return nil
}

//
// /network-instances/network-instance[name]/interfaces/interface[id]/config
//
func (h *NIDeleteVerifyHandler) NetworkInstanceInterfaceConfig(name string, id string, config *openconfig.NetworkInstanceInterfaceConfig) error {
	log.Debugf("NI/%s/%s/%s/%s/CONF: %s", h.ev, h.oper, name, id, config)

	if err := VerifyNIInterfaceConfig(id, config); err != nil {
		log.Errorf("NI/%s/%s/%s/%s/CONF: %s", h.ev, h.oper, name, id, err)
		return err
	}

	log.Debugf("NI/%s/%s/%s/%s: OK", h.ev, h.oper, name, id)
	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/static-routes/static[rtkey]
//
func (h *NIDeleteVerifyHandler) StaticRoute(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, route *openconfig.StaticRoute) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s: %s", h.ev, h.oper, name, prkey, rtkey, route)

	if prkey.Ident != openconfig.INSTALL_PROTOCOL_STATIC {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s: Invalid Protocol identifier.", h.ev, h.oper, name, prkey, rtkey)
	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/static-routes/static[rtkey]/next-hops/next-hop[index]
//
func (h *NIDeleteVerifyHandler) StaticRouteNexthop(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, nexthop *openconfig.StaticRouteNexthop) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s: %s", h.ev, h.oper, name, prkey, rtkey, index, nexthop)

	if chg := nexthop.GetChange(openconfig.OC_CONFIG_KEY); !chg {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s: do not change next-hop patrially.", h.ev, h.oper, name, prkey, rtkey, index)
	}

	_, nhType, err := nexthop.Config.GetNexthop()
	if err != nil {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s: invalid next-hop. %s", h.ev, h.oper, name, prkey, rtkey, index, err)
	}

	if nhType == openconfig.LOCAL_DEFINED_NEXT_HOP_LOCAL_LINK {
		if err := VerifyNIInterfaceRefConfig(nexthop.IfaceRef.Config); err != nil {
			return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s: %s", h.ev, h.oper, name, prkey, rtkey, index, err)
		}
	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/static-routes/static[rtkey]/next-hops/next-hop[index]/config
//
func (h *NIDeleteVerifyHandler) StaticRouteNexthopConfig(name string, prkey *openconfig.NetworkInstanceProtocolKey, rtkey *openconfig.StaticRouteKey, index string, config *openconfig.StaticRouteNexthopConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF: %s", h.ev, h.oper, name, prkey, rtkey, index, config)

	if chg := config.GetChange(openconfig.STATICROUTE_NEXTHOP_KEY); !chg {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s/%s/%s/CONF: next-hop not specified.", h.ev, h.oper, name, prkey, rtkey, index)
	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/ospfv2
//
func (h *NIDeleteVerifyHandler) Ospfv2(name string, key *openconfig.NetworkInstanceProtocolKey, ospf *openconfig.Ospfv2) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s: %s", h.ev, h.oper, name, key, ospf)

	if key.Ident != openconfig.INSTALL_PROTOCOL_OSPF {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s: Invalid Protocol identifier.", h.ev, h.oper, name, key)

	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/bgp
//
func (h *NIDeleteVerifyHandler) Bgp(name string, key *openconfig.NetworkInstanceProtocolKey, bgp *openconfig.Bgp) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s: %s", h.ev, h.oper, name, key, bgp)

	if key.Ident != openconfig.INSTALL_PROTOCOL_BGP {
		return fmt.Errorf("NI/%s/%s/%s/PROTOS/%s: Invalid Protocol identifier.", h.ev, h.oper, name, key)

	}

	return nil
}

//
// /network-instances/network-instance[name]/protocols/protocol[prkey]/bgp/neighbors/neighbor[addr]/apply-policy/config
//
func (h *NIDeleteVerifyHandler) BgpNeighborApplyPolicyConfig(name string, key *openconfig.NetworkInstanceProtocolKey, addr string, config *openconfig.PolicyApplyConfig) error {
	log.Debugf("NI/%s/%s/%s/PROTOS/%s/%s/APPLYPOL: %s", h.ev, h.oper, name, key, addr, config)

	// always success.

	return nil
}
