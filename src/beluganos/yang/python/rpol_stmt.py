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
XML policy-definition statement element
"""

# pylint: disable=too-few-public-methods

from beluganos.yang.python import constants
from beluganos.yang.python import elements


class PolicyStatement(elements.BaseElement):
    """
    policy-definition statement
    """

    _FIELDS = ("name", "config", "actions")

    def __init__(self, name):
        super(PolicyStatement, self).__init__("statement")
        self.name = name
        self.config = PolicyStatementConfig(name)
        self.actions = PolicyActions()


class PolicyStatementConfig(elements.BaseElement):
    """
    policy-definition statement config
    """

    _FIELDS = ("name",)

    def __init__(self, name):
        super(PolicyStatementConfig, self).__init__(constants.CONFIG)
        self.name = name


class PolicyActions(elements.BaseElement):
    """
    policy-definition actions element
    """

    _FIELDS = ("config", "action")

    def __init__(self, policy_result="ACCEPT_ROUTE"):
        super(PolicyActions, self).__init__("actions")
        self.config = PolicyActionsConfig(policy_result)
        self.action = None


class PolicyActionsConfig(elements.BaseElement):
    """
    policy-definition actions config element
    """

    _FIELDS = ("policy_result",)

    def __init__(self, policy_result="ACCEPT_ROUTE"):
        super(PolicyActionsConfig, self).__init__(constants.CONFIG)
        self.policy_result = policy_result
