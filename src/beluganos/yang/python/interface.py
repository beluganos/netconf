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
Interface module.
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements

# pylint: disable=too-few-public-methods

def parse_ifname(ifname):
    """
    parse interface name(<name>.<index>)
    """
    items = ifname.split(".", 2)
    if len(items) == 1:
        return ifname, "0"

    return items[0], items[1]


class InterfaceName(object):
    """
    Interface name
    """
    def __init__(self, ifname):
        self.iface, self.index = parse_ifname(ifname)
        self.ifname = ifname

    def __str__(self):
        return self.ifname


class InterfaceType(elements.Element):
    """
    Inteface type element.
    """

    _NSMAP = {constants.IANA_IFTYPE: constants.IANA_IFTYPE_NS}

    def __init__(self, iftype):
        text = "{0}:{1}".format(constants.IANA_IFTYPE, iftype)
        super(InterfaceType, self).__init__("type", text=text, nsmap=self._NSMAP)


class Interface(elements.BaseElement):
    """
    Interface element.
    """

    _FIELDS = ("name", "config", "subinterfaces", "ethernet")

    def __init__(self, name):
        super(Interface, self).__init__(constants.INTERFACE)
        self.name = name
        self.config = InterfaceConfig(name, "ethernetCsmacd")
        self.subinterfaces = elements.ListElement("subinterfaces")
        self.ethernet = InterfaceEthernet()


class InterfaceConfig(elements.BaseElement):
    """
    Interface config element.
    """

    _FIELDS = ("name", "iftype", "enabled", "mtu")

    def __init__(self, name, iftype):
        super(InterfaceConfig, self).__init__(constants.CONFIG)
        self.name = name
        self.iftype = InterfaceType(iftype)
        self.enabled = True
        self.mtu = None


class InterfaceEthernet(elements.BaseElement):
    """
    Interface/ethernet element
    """

    _NSMAP = {None: constants.INTERFACES_ETH_NS}
    _FIELDS = ("config",)

    def __init__(self):
        super(InterfaceEthernet, self).__init__(constants.INTERFACES_ETH, nsmap=self._NSMAP)
        self.config = InterfaceEthernetConfig()


class InterfaceEthernetConfig(elements.BaseElement):
    """
    Subinterface/ethernet/config element
    """

    _FIELDS = ("mac_address",)

    def __init__(self, mac_address=None):
        super(InterfaceEthernetConfig, self).__init__(constants.CONFIG)
        self.mac_address = mac_address
