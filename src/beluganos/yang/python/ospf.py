# -*- coding: utf-8 -*-

# Copyright (C) 2019 Nippon Telegraph and Telephone Corporation.
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
OSPF module
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements
from beluganos.yang.python.interface_ref import InterfaceRef
from beluganos.yang.python.interface import InterfaceName


class OspfArea(elements.BaseElement):
    """
    OSPF area element
    """

    _FIELDS = ("identifier", "config", "interfaces", "ranges")

    def __init__(self, areaid):
        super(OspfArea, self).__init__("area")
        self.identifier = areaid
        self.config = OspfAreaConfig(areaid)
        self.interfaces = elements.DictElement("interfaces")
        self.ranges = elements.DictElement("ranges")


class OspfAreaConfig(elements.BaseElement):
    """
    OSPF area config element
    """

    _FIELDS = ("identifier",)

    def __init__(self, areaid):
        super(OspfAreaConfig, self).__init__(constants.CONFIG)
        self.identifier = areaid


class OspfAreaInterface(elements.BaseElement):
    """
    OSPF area  interface element
    """

    _FIELDS = ("id", "config", "interface_ref", "timers")

    def __init__(self, iface_id):
        super(OspfAreaInterface, self).__init__("interface")
        ifname = InterfaceName(iface_id)
        self.id = iface_id
        self.config = OspfAreaInterfaceConfig(iface_id)
        self.interface_ref = InterfaceRef(ifname.iface, ifname.index)
        self.timers = OspfAreaInterfaceTimers()


class OspfAreaInterfaceConfig(elements.BaseElement):
    """
    OSPF area interface config element
    """

    _FIELDS = ("id", "metric", "passive")

    def __init__(self, iface_id, metric=None, passive=None):
        super(OspfAreaInterfaceConfig, self).__init__(constants.CONFIG)
        self.id = iface_id
        self.metric = metric
        self.passive = passive


class OspfAreaInterfaceTimers(elements.BaseElement):
    """
    OSPF area interface timers element
    """

    _FIELDS = ("config",)

    def __init__(self, hello_interval=None, dead_interval=None):
        super(OspfAreaInterfaceTimers, self).__init__("timers")
        self.config = OspfAreaInterfaceTimersConfig(hello_interval, dead_interval)


class OspfAreaInterfaceTimersConfig(elements.BaseElement):
    """
    OSPF area interface timers config element
    """

    _FIELDS = ("hello_interval", "dead_interval")

    def __init__(self, hello_interval=None, dead_interval=None):
        super(OspfAreaInterfaceTimersConfig, self).__init__(constants.CONFIG)
        self.dead_interval = dead_interval
        self.hello_interval = hello_interval


class OspfAreaRange(elements.BaseElement):
    """
    OSPF area range element
    """

    _FIELDS = ("ip", "prefix_length", "config")

    def __init__(self, ip, prefix_length):
        super(OspfAreaRange, self).__init__("range")
        self.ip = ip
        self.prefix_length = prefix_length
        self.config = OspfAreaRangeConfig(ip, prefix_length)


class OspfAreaRangeConfig(elements.BaseElement):
    """
    OSPF area range config element
    """

    _FIELDS = ("ip", "prefix_length")

    def __init__(self, ip, prefix_length):
        super(OspfAreaRangeConfig, self).__init__(constants.CONFIG)
        self.ip = ip
        self.prefix_length = prefix_length
