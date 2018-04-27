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
OSPF module
"""

# pylint: disable=too-few-public-methods
# pylint: disable=invalid-name

from beluganos.yang.python import constants
from beluganos.yang.python import elements
from beluganos.yang.python import policy_types
from beluganos.yang.python.interface_ref import InterfaceRef
from beluganos.yang.python.interface import InterfaceName

class Ospfv2(elements.BaseElement):
    """
    OFDPA element
    """

    _FIELDS = ("_global", "areas")

    def __init__(self, name="ospfv2", router_id=None):
        super(Ospfv2, self).__init__("ospfv2")
        self._global = Ospfv2Global(router_id)
        self.areas = elements.DictElement("areas")
        self.name = name

    @staticmethod
    def get_identifier():
        """
        network instance protocol type
        """
        return policy_types.InstallProtocolType(constants.POLICY_PROTOCOL_OSPF)


    def get_name(self):
        """
        network instance protocol name
        """
        return self.name


class Ospfv2Global(elements.BaseElement):
    """
    OSPF global element
    """

    _FIELDS = ("config",)

    def __init__(self, router_id=None):
        super(Ospfv2Global, self).__init__("global")
        self.config = Ospfv2GlobalConfig(router_id)


class Ospfv2GlobalConfig(elements.BaseElement):
    """
    OSPF global config element
    """

    _FIELDS = ("router_id",)

    def __init__(self, router_id=None):
        super(Ospfv2GlobalConfig, self).__init__(constants.CONFIG)
        self.router_id = router_id


class OspfArea(elements.BaseElement):
    """
    OSPF area element
    """

    _FIELDS = ("identifier", "config", "interfaces")

    def __init__(self, areaid):
        super(OspfArea, self).__init__("area")
        self.identifier = areaid
        self.config = OspfAreaConfig(areaid)
        self.interfaces = elements.DictElement("interfaces")


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
