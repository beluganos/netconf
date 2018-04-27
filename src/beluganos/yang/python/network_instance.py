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
Betwork instance module
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements
from beluganos.yang.python import mpls

# pylint: disable=too-few-public-methods
# pylint: disable=invalid-name

class NetworkInstanceType(elements.Element):
    """
    network-instance type element.
    """

    _NSMAP = {constants.NETWORK_INSTANCE_TYPES: constants.NETWORK_INSTANCE_TYPES_NS}

    def __init__(self, nitype):
        text = "{0}:{1}".format(constants.NETWORK_INSTANCE_TYPES, nitype)
        super(NetworkInstanceType, self).__init__("type", text=text, nsmap=self._NSMAP)


class NetworkInstance(elements.BaseElement):
    """
    Network Instance element.
    """

    _FIELDS = ("name", "config", "loopbacks", "interfaces", "mpls", "protocols")

    def __init__(self, name):
        super(NetworkInstance, self).__init__(constants.NETWORK_INSTANCE)
        self.name = name
        self.config = NetworkInstanceConfig(name)
        self.loopbacks = elements.DictElement("loopbacks")
        self.interfaces = elements.DictElement(constants.INTERFACES)
        self.mpls = mpls.Mpls()
        self.protocols = elements.DictElement(constants.NETWORK_INSTANCE_PROTOCOLS)


class NetworkInstanceLoopback(elements.BaseElement):
    """
    Network Instance loopback element.
    """

    _FIELDS = ("_id", "config", "addresses")

    def __init__(self, _id):
        super(NetworkInstanceLoopback, self).__init__("loopback")
        self._id = _id
        self.config = NetworkInstanceLoopbackConfig(_id)
        self.addresses = elements.ListElement("addresses")


class NetworkInstanceLoopbackConfig(elements.BaseElement):
    """
    Network Instance loopback config element.
    """

    _FIELDS = ("_id",)

    def __init__(self, _id):
        super(NetworkInstanceLoopbackConfig, self).__init__("config")
        self._id = _id


class NetworkInstanceLoopbackAddress(elements.BaseElement):
    """
    Network Instance loopback address element.
    """

    _FIELDS = ("index", "config")

    def __init__(self, index):
        super(NetworkInstanceLoopbackAddress, self).__init__("address")
        self.index = index
        self.config = NetworkInstanceLoopbackAddressConfig(index)


class NetworkInstanceLoopbackAddressConfig(elements.BaseElement):
    """
    Network Instance loopback address config element.
    """

    _FIELDS = ("index", "ip", "prefix_length")

    def __init__(self, index):
        super(NetworkInstanceLoopbackAddressConfig, self).__init__("config")
        self.index = index
        self.ip = None
        self.prefix_length = None


class NetworkInstanceConfig(elements.BaseElement):
    """
    Network instance config element.
    """

    _FIELDS = ("name", "nitype", "router_id", "route_distinguisher", "route_target")

    def __init__(self, name):
        super(NetworkInstanceConfig, self).__init__(constants.CONFIG)
        self.name = name
        self.nitype = NetworkInstanceType("DEFAULT_INSTANCE")
        self.router_id = None
        self.route_distinguisher = None
        self.route_target = None


class NetworkInstanceInterface(elements.BaseElement):
    """
    Network instance interface element.
    """

    _FIELDS = ("id", "config")

    def __init__(self, iface_id, interface=None, subinterface=None):
        super(NetworkInstanceInterface, self).__init__(constants.INTERFACE)
        self.id = iface_id
        self.config = NetworkInstanceInterfaceConfig(iface_id, interface, subinterface)


class NetworkInstanceInterfaceConfig(elements.BaseElement):
    """
    Network instance interface config element.
    """

    _FIELDS = ("id", "interface", "subinterface")

    def __init__(self, iface_id, interface=None, subinterface=None):
        super(NetworkInstanceInterfaceConfig, self).__init__(constants.CONFIG)
        self.id = iface_id
        self.interface = interface
        self.subinterface = subinterface


class NetworkInstanceProtocol(elements.BaseElement):
    """
    Neteork instance protocol element.
    """

    _FIELDS = ("identifier", "name", "config", "proto")

    def __init__(self, proto):
        super(NetworkInstanceProtocol, self).__init__(constants.NETWORK_INSTANCE_PROTOCOL)

        identifier = proto.get_identifier()
        name = proto.get_name()
        self.identifier = identifier
        self.name = name
        self.config = NetworkInstanceProtocolConfig(identifier, name)
        self.proto = proto


class NetworkInstanceProtocolConfig(elements.BaseElement):
    """
    Neteork instance protocol config element.
    """

    _FIELDS = ("identifier", "name")

    def __init__(self, identifier, name):
        super(NetworkInstanceProtocolConfig, self).__init__(constants.CONFIG)
        self.identifier = identifier
        self.name = name
