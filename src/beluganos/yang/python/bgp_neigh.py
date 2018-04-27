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
Network Instance Protocols(BGP) subsubmodule
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements

# pylint: disable=too-few-public-methods
class BgpNeighbor(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor) element.
    """

    _FIELDS = ("neighbor_address", "config", "timers", "transport", "afi_safis", "apply_policy")

    def __init__(self, address, peer_as=None, local_as=None):
        super(BgpNeighbor, self).__init__(constants.BGP_NEIGHBOR)
        self.neighbor_address = address
        self.config = BgpNeighborConfig(address, peer_as, local_as)
        self.timers = BgpNeighborTimers()
        self.transport = BgpNeighborTransport()
        self.afi_safis = elements.ListElement("afi-safis")
        self.apply_policy = BgpNeighborApplyPolicy()


class BgpNeighborConfig(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/config) element.
    """

    _FIELDS = ("neighbor_address", "peer_as", "local_as")

    def __init__(self, address, peer_as=None, local_as=None):
        super(BgpNeighborConfig, self).__init__(constants.CONFIG)
        self.neighbor_address = address
        self.peer_as = peer_as
        self.local_as = local_as


class BgpNeighborTimers(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/timers) element.
    """

    _FIELDS = ("config",)

    def __init__(self, hold_time=None, keepalive_interval=None):
        super(BgpNeighborTimers, self).__init__("timers")
        self.config = BgpNeighborTimersConfig(hold_time, keepalive_interval)


class BgpNeighborTimersConfig(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/timers/config) element.
    """

    _FIELDS = ("hold_time", "keepalive_interval")

    def __init__(self, hold_time=None, keepalive_interval=None):
        super(BgpNeighborTimersConfig, self).__init__(constants.CONFIG)
        self.hold_time = hold_time
        self.keepalive_interval = keepalive_interval


class BgpNeighborTransport(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/transport) element.
    """

    _FIELDS = ("config",)

    def __init__(self):
        super(BgpNeighborTransport, self).__init__("transport")
        self.config = BgpNeighborTransportConfig()


class BgpNeighborTransportConfig(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/transport/config) element.
    """

    _FIELDS = ("local_address",)

    def __init__(self, local_address=None):
        super(BgpNeighborTransportConfig, self).__init__(constants.CONFIG)
        self.local_address = local_address


class BgpNeighborApplyPolicy(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/apply-policy) element.
    """

    _FIELDS = ("config",)

    def __init__(self):
        super(BgpNeighborApplyPolicy, self).__init__("apply-policy")
        self.config = BgpNeighborApplyPolicyConfig()


class BgpNeighborApplyPolicyConfig(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/apply-policy/config) element.
    """

    _FIELDS = (
        "default_import_policy",
        "default_export_policy"
    )

    def __init__(self,
                 default_import_policy=None,
                 default_export_policy=None):
        super(BgpNeighborApplyPolicyConfig, self).__init__(constants.CONFIG)
        self.import_policy = list()
        self.export_policy = list()
        self.default_import_policy = default_import_policy
        self.default_export_policy = default_export_policy


    def xml_element(self):
        """
        create XML element.
        """
        elm = super(BgpNeighborApplyPolicyConfig, self).xml_element()
        for name in self.import_policy:
            elm.append(elements.xml_element("import-policy", name))
        for name in self.export_policy:
            elm.append(elements.xml_element("export-policy", name))
        return elm

class BgpNeighborAfiSafi(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/afi-safis/afi-safi) element.
    """

    _FIELDS = ("afi_safi_name", "config")

    def __init__(self, name):
        super(BgpNeighborAfiSafi, self).__init__("afi-safi")
        self.afi_safi_name = BgpNeighborAfiSafiName(name)
        self.config = BgpNeighborAfiSafiConfig(name)


class BgpNeighborAfiSafiConfig(elements.BaseElement):
    """
    Neteork instance protocol(BGP/neighbor/afi-safis/afi-safi/config) element.
    """

    _FIELDS = ("afi_safi_name",)

    def __init__(self, name):
        super(BgpNeighborAfiSafiConfig, self).__init__(constants.CONFIG)
        self.afi_safi_name = BgpNeighborAfiSafiName(name)


class BgpNeighborAfiSafiName(elements.Element):

    _NSMAP = {
        constants.BGP_TYPES: constants.BGP_TYPES_NS,
    }

    def __init__(self, name):
        super(BgpNeighborAfiSafiName, self).__init__(
            "afi-safi-name",
            text="{0}:{1}".format(constants.BGP_TYPES, name),
            nsmap=self._NSMAP,
        )
