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
Policy types module.
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements

# pylint: disable=too-few-public-methods
class InstallProtocolType(elements.Element):
    """
    Network instance protocol identifier.
    """

    _NSMAP = {
        constants.NETWORK_INSTANCE_PROTOCOL_TYPES: constants.NETWORK_INSTANCE_PROTOCOL_TYPES_NS,
    }

    def __init__(self, identifier):
        text = "{0}:{1}".format(constants.NETWORK_INSTANCE_PROTOCOL_TYPES, identifier)
        super(InstallProtocolType, self).__init__(
            constants.NETWORK_INSTANCE_PROTOCOL_IDENT,
            text=text,
            nsmap=self._NSMAP,
        )
