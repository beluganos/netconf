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
beluganos-routing-policy module
"""

# pylint: disable=too-few-public-methods

from beluganos.yang.python import constants
from beluganos.yang.python import elements


class RoutingPolicy(elements.BaseElement):
    """
    routing-policy element
    """

    _FIELDS = ("policy_definitions",)

    def __init__(self):
        super(RoutingPolicy, self).__init__("routing-policy", space=constants.ROUTING_POLICY_NS)
        self.policy_definitions = elements.ListElement("policy-definitions")


class PolicyDefinition(elements.BaseElement):
    """
    policy-definition element
    """

    _FIELDS = ("name", "config", "statements")

    def __init__(self, name):
        super(PolicyDefinition, self).__init__("policy-definition")
        self.name = name
        self.config = PolicyDefinitionConfig(name)
        self.statements = elements.ListElement("statements")


class PolicyDefinitionConfig(elements.BaseElement):
    """
    policy-definition config element
    """

    _FIELDS = ("name",)

    def __init__(self, name):
        super(PolicyDefinitionConfig, self).__init__(constants.CONFIG)
        self.name = name
