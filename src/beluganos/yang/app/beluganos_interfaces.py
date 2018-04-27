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

"""
XML for beluganos-interfaces builder.
"""

# pylint: disable=no-member

from lxml import etree as ET

from netaddr import IPNetwork
from beluganos.yang.python import constants
from beluganos.yang.python.elements import ListElement
from beluganos.yang.python.address import Address
from beluganos.yang.python.interface import Interface
from beluganos.yang.python.subinterface import Subinterface


def new_subiface(index, device):
    """
    Create Subinterface.
    """
    subiface = Subinterface(index)

    if "addresses" in device:
        for addr in device["addresses"]:
            addr = IPNetwork(addr)
            subiface.ipv4.addresses.append(Address(addr.ip, addr.prefixlen))

    if "mtu" in device:
        subiface.ipv4.config.mtu = device["mtu"]

    return subiface


def new_iface(ifname, device):
    """
    Create Interface
    """
    iface = Interface(ifname)

    if "mtu" in device:
        iface.config.mtu = device["mtu"]

    if "mac-address" in device:
        iface.ethernet.config.mac_address = device["mac-address"]

    return iface


def new_ifaces(cfg):
    """
    Create Interfaces
    """
    ifaces = ListElement("interfaces", space=constants.INTERFACES_NS)
    ifaces_dict = dict()

    for ifname, eth in cfg["ethernets"].items():
        iface = new_iface(ifname, eth)
        iface.subinterfaces.append(new_subiface(0, eth))

        ifaces_dict[ifname] = iface
        ifaces.append(iface)

    if "vlans" in cfg:
        for ifname, vlan in cfg["vlans"].items():
            vlanif = new_subiface(vlan["id"], vlan)
            iface = ifaces_dict[vlan["link"]]
            iface.subinterfaces.append(vlanif)

    return ifaces


def read_config(path):
    """
    Read yaml file and return dict.
    """
    import yaml
    with open(path, "r") as cfg:
        return yaml.load(cfg)


def read_configs(paths):
    """
    Read yaml files and merged dict.
    """
    eths = dict()
    vlans = dict()
    bonds = dict()
    for path in paths:
        cfg = read_config(path)
        ifaces = cfg.get("network", dict())

        if "ethernets" in ifaces:
            eths.update(ifaces["ethernets"])
        if "vlans" in ifaces:
            vlans.update(ifaces["vlans"])
        if "bonds" in ifaces:
            bonds.update(ifaces["bonds"])

    return dict(
        ethernets=eths,
        vlans=vlans,
        bonds=bonds
    )


def _getopts():
    from optparse import OptionParser
    parser = OptionParser()
    parser.add_option("-v", "--verbose", dest="verbose", action="store_true",
                      help="show detail messages.")
    return parser.parse_args()


def _main():
    _, args = _getopts()

    cfg = read_configs(args)
    ifaces = new_ifaces(cfg)
    print ET.tostring(ifaces.xml_element(), pretty_print=True)


if __name__ == "__main__":
    _main()
