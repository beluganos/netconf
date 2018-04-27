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
BGP policy action module
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements


class BgpActions(elements.BaseElement):
    """
    bgp-actions
    """

    _FIELDS = ("config", )

    def __init__(self):
        super(BgpActions, self).__init__("bgp-actions", space=constants.BGP_POLICY_NS)
        self.config = BgpActionsConfig()


class BgpActionsConfig(elements.BaseElement):
    """
    bgp-actions config
    """

    _FIELDS = ("set_local_pref", "set_next_hop")

    def __init__(self):
        super(BgpActionsConfig, self).__init__(constants.CONFIG)
        self.set_local_pref = None
        self.set_next_hop = None
