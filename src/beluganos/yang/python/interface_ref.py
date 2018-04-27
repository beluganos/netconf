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
Interface ref module
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements

# pylint: disable=too-few-public-methods
class InterfaceRef(elements.BaseElement):
    """
    Interface ref element
    """

    _FIELDS = ("config",)

    def __init__(self, interface, subinterface=0):
        super(InterfaceRef, self).__init__(constants.INTERFACE_REF)
        self.config = InterfaceRefConfig(interface, subinterface)


class InterfaceRefConfig(elements.BaseElement):
    """
    Interface ref config element
    """

    _FIELDS = ("interface", "subinterface")

    def __init__(self, interface, subinterface=0):
        super(InterfaceRefConfig, self).__init__(constants.CONFIG)
        self.interface = interface
        self.subinterface = subinterface
