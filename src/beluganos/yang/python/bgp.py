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
from beluganos.yang.python.bgp_global import BgpGlobal
from beluganos.yang.python.bgp_zebra import BgpZebra
from beluganos.yang.python.policy_types import InstallProtocolType

# pylint: disable=too-few-public-methods
class Bgp(elements.BaseElement):
    """
    Neteork instance protocol(BGP) element.
    """

    _FIELDS = ("_global", "zebra", "neighbors")
    def __init__(self, name):
        super(Bgp, self).__init__(constants.BGP)
        self.name = name
        self._global = BgpGlobal()
        self.zebra = BgpZebra()
        self.neighbors = elements.DictElement("neighbors")

    @staticmethod
    def get_identifier():
        """
        network instance protocol type
        """
        return InstallProtocolType(constants.POLICY_PROTOCOL_BGP)

    def get_name(self):
        """
        network instance protocol name
        """
        return self.name
