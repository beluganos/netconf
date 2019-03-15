#! /usr/bin/env python
# -*- coding: utf-8 -*-

# Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import collections

"""
XML for beluganos-network-instance builder.
"""

# pylint: disable=invalid-name

from lxml import etree as ET
from netaddr import IPNetwork

from beluganos.yang.python import constants
from beluganos.yang.python.elements import ListElement
from beluganos.yang.python.interface import InterfaceName
from beluganos.yang.python.network_instance import (NetworkInstance,
                                                    NetworkInstanceType,
                                                    NetworkInstanceInterface,
                                                    NetworkInstanceLoopback,
                                                    NetworkInstanceLoopbackAddress,
                                                    NetworkInstanceProtocol)
from beluganos.yang.python.mpls import (Mpls,
                                        MplsGlobalInterface,
                                        MplsLdpGlobalDiscoveryInterface)
from beluganos.yang.python.static_route import (StaticRoute,
                                                StaticRoutes,
                                                StaticRouteNexthop)
from beluganos.yang.python.ospf import OspfArea, OspfAreaInterface, OspfAreaRange
from beluganos.yang.python.ospfv2 import Ospfv2
from beluganos.yang.python.ospfv3 import Ospfv3
from beluganos.yang.python.bgp import Bgp
from beluganos.yang.python.bgp_neigh import BgpNeighbor, BgpNeighborAfiSafi

_ZEBRA_URL = "unix:/var/run/frr/zserv.api"
_ZEBRA_VER = 5

def get_sub_dict(d, *names):
    """
    get sub dictionary recursive.
    """
    for name in names:
        d = d.get(name, None)
        if d is None:
            return dict()
    return d


def new_bgp_neigh(neighaddr, neighcfg):
    neigh = BgpNeighbor(neighaddr)
    neigh.config.peer_as = neighcfg.get("peer-as", None)
    neigh.config.local_as = neighcfg.get("local-as", None)

    timercfg = get_sub_dict(neighcfg, "timers")
    if timercfg:
        neigh.timers.config.hold_time = timercfg.get("hold-time", None)
        neigh.timers.config.keepalive_interval = timercfg.get("keepalive-interval", None)

    neigh.transport.config.local_address = neighcfg.get("transport-local-address", None)

    policycfg = get_sub_dict(neighcfg, "apply-policy")
    if policycfg:
        neigh.apply_policy.config.import_policy = policycfg.get("import-policy", list())
        neigh.apply_policy.config.export_policy = policycfg.get("export-policy", list())
        neigh.apply_policy.config.default_import_policy = policycfg.get("default-import-policy", list())
        neigh.apply_policy.config.default_export_policy = policycfg.get("default-export-policy", list())

    afisaficfg = get_sub_dict(neighcfg, "afi-safis")
    if afisaficfg:
        for afisafi_name in afisaficfg:
            neigh.afi_safis.append(BgpNeighborAfiSafi(afisafi_name))

    return neigh


def new_bgp(cfg):
    """
    create protocols(bgp)
    """
    bgpcfg = get_sub_dict(cfg, "bgp")
    if not bgpcfg:
        return None

    bgp = Bgp("BGP")
    bgp._global.config._as = bgpcfg.get("as", None)
    bgp._global.config.router_id = bgpcfg.get("router-id", None)

    zebracfg = get_sub_dict(bgpcfg, "zebra")
    if zebracfg:
        bgp.zebra.config.enabled = zebracfg.get("enabled", None)
        bgp.zebra.config.redistributes = zebracfg.get("redistribute-routes", list())
        if bgp.zebra.config.enabled:
            bgp.zebra.config.version = zebracfg.get("version", _ZEBRA_VER)
            bgp.zebra.config.url = zebracfg.get("url", _ZEBRA_URL)

    neighscfg = get_sub_dict(bgpcfg, "neighbors")
    for neighaddr, neighcfg in neighscfg.items():
        bgp.neighbors.append(neighaddr, new_bgp_neigh(neighaddr, neighcfg))

    return NetworkInstanceProtocol(bgp)


def new_ospf_iface(ifname, ifcfg):
    iface = OspfAreaInterface(ifname)
    iface.config.metric = ifcfg.get("metric", None)
    iface.config.passive = ifcfg.get("passive", None)
    iface.timers.config.hello_interval = ifcfg.get("hello-interval", None)
    iface.timers.config.dead_interval = ifcfg.get("dead-interval", None)
    return iface


def new_ospf_area(areaid, areacfg):
    area = OspfArea(areaid)
    for ifname, ifcfg in sorted(areacfg.items()):
        area.interfaces.append(ifname, new_ospf_iface(ifname, ifcfg))

    return area


def new_ospf_area_range(rangeid, rangecfg):
    ipnet = IPNetwork(rangeid)
    rng = OspfAreaRange(ipnet.ip, ipnet.prefixlen)
    return rng


def new_ospf(cfg):
    """
    create protocols(ospf)
    """
    ospfcfg = get_sub_dict(cfg, "ospfv2")
    if not ospfcfg:
        return None

    ospf = Ospfv2(router_id=ospfcfg["router-id"])

    areascfg = get_sub_dict(ospfcfg, "areas")
    if areascfg:
        for areaid, areacfg in sorted(areascfg.items()):
            ospf.areas.append(areaid, new_ospf_area(areaid, areacfg))

    return NetworkInstanceProtocol(ospf)


def new_ospfv3(cfg):
    """
    create protocols(ospfv3)
    """
    ospfcfg = get_sub_dict(cfg, "ospfv3")
    if not ospfcfg:
        return None

    ospfv3 = Ospfv3(router_id=ospfcfg["router-id"])

    areascfg = get_sub_dict(ospfcfg, "areas")
    if areascfg:
        for areaid, areacfg in sorted(areascfg.items()):
            ospfv3.areas.append(areaid, new_ospf_area(areaid, areacfg))

    rangescfg = get_sub_dict(ospfcfg, "ranges")
    if rangescfg:
        for areaid, rangecfg in sorted(rangescfg.items()):
            area = ospfv3.areas.get(areaid)
            if not area:
                continue

            for rangeid in sorted(rangecfg):
                area.ranges.append(rangeid, new_ospf_area_range(rangeid, rangecfg))

    return NetworkInstanceProtocol(ospfv3)


def new_route(cfg):
    """
    create protocols(static-route)
    """
    routescfg = get_sub_dict(cfg, "routes")
    if not routescfg:
        return None

    routes = StaticRoutes()

    for dest, destcfgs in sorted(routescfg.items()):
        ipnet = IPNetwork(dest)
        route = StaticRoute(ipnet.ip, ipnet.prefixlen)

        for index, destcfg in enumerate(destcfgs):
            nexthop = StaticRouteNexthop(index, **destcfg)
            route.next_hops.append(nexthop)

        routes.append(dest, route)

    return NetworkInstanceProtocol(routes)


def new_mpls(cfg):
    """
    create mpls
    """
    ldpcfg = get_sub_dict(cfg, "mpls", "ldp")
    if not ldpcfg:
        return None

    mpls = Mpls()
    mpls._global.config.set_null_label(ldpcfg.get("null-label", None))

    ldp_global = mpls.signaling_protocols.ldp._global
    ldp_global.config.lsr_id = ldpcfg.get("router-id", None)

    ipv4cfg = ldpcfg.get("ipv4", None)
    if ipv4cfg:
        ipv4 = ldp_global.address_families.ipv4
        ipv4.config.transport_address = ipv4cfg.get("transport-address")
        ipv4.config.session_ka_holdtime = ipv4cfg.get("session-holdtime", None)
        ipv4.config.label_policy.advertise.egress_explicit_null.enable = ipv4cfg.get("egress-explicit-null", None)

    timcfg = ldpcfg.get("timers", None)
    if timcfg:
        ldp_global.discovery.interfaces.config.hello_holdtime = timcfg.get("hello-holdtime", None)
        ldp_global.discovery.interfaces.config.hello_interval = timcfg.get("hello-interval", None)

    for ifname, ifcfg in sorted(get_sub_dict(ldpcfg, "interfaces").items()):
        ifname = InterfaceName(ifname)

        iface = MplsGlobalInterface(str(ifname), ifname.iface, ifname.index)
        mpls._global.interface_attributes.append(ifname, iface)

        iface = MplsLdpGlobalDiscoveryInterface(str(ifname), ifname.iface, ifname.index)
        if ifcfg:
            iface.config.hello_holdtime = ifcfg.get("hello-holdtime", None)
            iface.config.hello_interval = ifcfg.get("hello-interval", None)
        ldp_global.discovery.interfaces.append(ifname, iface)

    return mpls


def new_ni_iface(iface_id, cfg):
    """
    create interface
    """
    ifname = InterfaceName(iface_id)
    iface = NetworkInstanceInterface(iface_id)
    if "both" in cfg:
        iface.config.interface = ifname.iface
        iface.config.subinterface = ifname.index
    if "iface" in cfg:
        iface.config.interface = ifname.iface
    if "subif" in cfg:
        iface.config.subinterface = ifname.index
    return iface


def new_ni_loopback(cfg):
    locfgs = cfg.get("loopbacks", None)
    if not locfgs:
        return dict()

    loopbacks = dict()
    for loname, locfg in locfgs.items():
        lo = NetworkInstanceLoopback(loname)
        for addr in locfg:
            addr = IPNetwork(addr)
            loaddr = NetworkInstanceLoopbackAddress(len(lo.addresses))
            loaddr.config.ip = addr.ip
            loaddr.config.prefix_length = addr.prefixlen
            lo.addresses.append(loaddr)

        loopbacks[loname] = lo

    return loopbacks


def new_ni(name, cfg):
    """
    create netwotk-instance
    """
    ni = NetworkInstance(name)
    ni.config.nitype = NetworkInstanceType(cfg.get("type", "DEFAULT_INSTANCE"))
    ni.config.router_id = cfg.get("router-id", None)
    ni.config.route_distinguisher = cfg.get("route-distinguisher", None)
    ni.config.route_target = cfg.get("route-target", None)

    ifacecfgs = cfg.get("interfaces", None)
    if ifacecfgs:
        for ifname, ifcfg in sorted(ifacecfgs.items()):
            ni.interfaces.append(ifname, new_ni_iface(ifname, ifcfg))

    for loname, lo in new_ni_loopback(cfg).items():
        ni.loopbacks.append(loname, lo)

    ni.mpls = new_mpls(cfg)

    static_routes = new_route(cfg)
    if static_routes:
        ni.protocols.append("static", static_routes)

    ospfv2 = new_ospf(cfg)
    if ospfv2:
        ni.protocols.append("ospfv2", ospfv2)

    ospfv3 = new_ospfv3(cfg)
    if ospfv3:
        ni.protocols.append("ospfv3", ospfv3)

    bgp = new_bgp(cfg)
    if bgp:
        ni.protocols.append("bgp", bgp)

    return ni


def new_top(nis):
    """
    create network-instances
    """
    top = ListElement("network-instances", space=constants.NETWORK_INSTANCE_NS)
    for name, ni in nis.items():
        top.append(ni)

    return top


def read_config(path):
    """
    Read yaml file and return dict.
    """
    import yaml
    with open(path, "r") as cfg:
        return yaml.load(cfg).get("network-instances", dict())


def read_configs(paths):
    """
    Read yaml files and merged dict.
    """
    nis = collections.OrderedDict()
    for path in paths:
        cfg = read_config(path)
        for niname, nicfg in cfg.items():
            ni = new_ni(niname, nicfg)
            if niname in nis:
                nis[niname].merge_element(ni)
            else:
                nis[niname] = ni

    return nis


def _getopts():
    from optparse import OptionParser
    parser = OptionParser()
    parser.add_option("-v", "--verbose", dest="verbose", action="store_true",
                      help="show detail messages.")
    return parser.parse_args()


def _main():
    _, args = _getopts()

    nis = read_configs(args)
    top = new_top(nis)
    print ET.tostring(top.xml_element(), pretty_print=True)


if __name__ == "__main__":
    _main()
