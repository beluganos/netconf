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
interface/ip module.
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements


# pylint: disable=too-few-public-methods
# pylint: disable=invalid-name
class Address(elements.BaseElement):
    """
    interface/ip address class
    """

    _FIELDS = ("ip", "config")

    def __init__(self, ipaddr, prefix_length):
        super(Address, self).__init__(constants.ADDRESS)
        self.ip = ipaddr
        self.config = AddressConfig(ipaddr, prefix_length)


class AddressConfig(elements.BaseElement):
    """
    Address config class.
    """

    _FIELDS = ("ip", "prefix_length")

    def __init__(self, ipaddr, prefix_length):
        super(AddressConfig, self).__init__(constants.CONFIG)
        self.ip = ipaddr
        self.prefix_length = prefix_length
