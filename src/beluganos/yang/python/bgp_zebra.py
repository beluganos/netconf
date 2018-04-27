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
class BgpZebra(elements.BaseElement):
    """
    Network instance protocol(bgp/xebra) element
    """

    _FIELDS = ("config",)

    def __init__(self):
        super(BgpZebra, self).__init__(constants.BGP_ZEBRA)
        self.config = BgpZebraConfig()


class BgpZebraConfig(elements.BaseElement):
    """
    Network instance protocol(bgp/zebra/config) element
    """

    _FIELDS = ("enabled", "version", "url")

    def __init__(self):
        super(BgpZebraConfig, self).__init__(constants.CONFIG)
        self.enabled = None
        self.version = None
        self.url = None
        self.redistributes = list()


    def xml_element(self):
        """
        create XML element.
        """
        elm = super(BgpZebraConfig, self).xml_element()
        for redist in self.redistributes:
            elm.append(elements.xml_element("redistribute-routes", redist))
        return elm
