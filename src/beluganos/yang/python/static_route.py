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
Static route submodule
"""

# pylint: disable=too-few-public-methods
# pylint: disable=invalid-name

from netaddr import IPAddress

from beluganos.yang.python import constants
from beluganos.yang.python import elements
from beluganos.yang.python import policy_types
from beluganos.yang.python.interface_ref import InterfaceRef
from beluganos.yang.python.interface import InterfaceName

class StaticRoutes(elements.DictElement):
    """
    statis routes element.
    """
    def __init__(self, name="static"):
        super(StaticRoutes, self).__init__("static-routes")
        self.name = name


    @staticmethod
    def get_identifier():
        """
        network instance protocol type
        """
        return policy_types.InstallProtocolType(constants.POLICY_PROTOCOL_STATIC)


    def get_name(self):
        """
        network instance protocol name
        """
        return self.name


class StaticRoute(elements.BaseElement):
    """
    Static route element
    """

    _FIELDS = ("ip", "prefix_length", "config", "next_hops")

    def __init__(self, ip, prefix_length):
        super(StaticRoute, self).__init__("static")
        self.ip = ip
        self.prefix_length = prefix_length
        self.config = StaticRouteConfig(ip, prefix_length)
        self.next_hops = elements.ListElement("next-hops")


class StaticRouteConfig(elements.BaseElement):
    """
    Static route config element
    """

    _FIELDS = ("ip", "prefix_length")

    def __init__(self, ip, prefix_length):
        super(StaticRouteConfig, self).__init__(constants.CONFIG)
        self.ip	= ip
        self.prefix_length = prefix_length


class StaticRouteNexthop(elements.BaseElement):
    """
    Static route nexthop element.
    """

    _FIELDS = ("index", "config", "interface_ref")

    def __init__(self, index, via=None, dev=None, drop=None):
        super(StaticRouteNexthop, self).__init__("next-hop")
        self.index = index
        self.interface_ref = None

        if via is not None:
            nexthop = IPAddress(via)
            self.config = StaticRouteNexthopConfig(index, nexthop)

        elif dev is not None:
            ifname = InterfaceName(dev)
            self.config = StaticRouteNexthopConfig(index, "LOCAL_LINK")
            self.interface_ref = InterfaceRef(ifname.iface, ifname.index)

        elif drop:
            self.config = StaticRouteNexthopConfig(index, "DROP")


class StaticRouteNexthopConfig(elements.BaseElement):
    """
    Static route nexthop config element.
    """

    _FIELDS = ("index", "next_hop")

    def __init__(self, index, next_hop):
        super(StaticRouteNexthopConfig, self).__init__(constants.CONFIG)
        self.index = index
        self.next_hop = next_hop
