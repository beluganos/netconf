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
Subinterface module.
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements


# pylint: disable=too-few-public-methods
class Subinterface(elements.BaseElement):
    """
    Subinterface element.
    """

    _FIELDS = ("index", "config", "ipv4")

    def __init__(self, index):
        super(Subinterface, self).__init__(constants.SUBINTERFACE)
        self.index = index
        self.config = SubinterfaceConfig(index)
        self.ipv4 = SubinterfaceIPv4()


class SubinterfaceConfig(elements.BaseElement):
    """
    Subinterface config element.
    """

    _FIELDS = ("index", "enabled")

    def __init__(self, index):
        super(SubinterfaceConfig, self).__init__(constants.CONFIG)
        self.index = index
        self.enabled = True


class SubinterfaceIPv4(elements.BaseElement):
    """
    Subinterface/ipv4 element
    """

    _NSMAP = {None: constants.INTERFACES_IP_NS}
    _FIELDS = ("addresses", "config")

    def __init__(self):
        super(SubinterfaceIPv4, self).__init__(constants.INTERFACES_IP_V4, nsmap=self._NSMAP)
        self.addresses = elements.ListElement("addresses")
        self.config = SubinterfaceIPv4Config()


class SubinterfaceIPv4Config(elements.BaseElement):
    """
    Subinterface/ipv4/config element
    """

    _FIELDS = ("mtu",)

    def __init__(self):
        super(SubinterfaceIPv4Config, self).__init__(constants.CONFIG)
        self.mtu = None
