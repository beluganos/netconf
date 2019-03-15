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
