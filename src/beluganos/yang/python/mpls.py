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
MPLS module
"""

from beluganos.yang.python import constants
from beluganos.yang.python import elements
from beluganos.yang.python import interface_ref

# pylint: disable=too-few-public-methods
class MpldNullLabel(elements.Element):
    """
    MPLS null label element
    """

    _NSMAP = {constants.MPLS_TYPES: constants.MPLS_TYPES_NS}

    def __init__(self, nulllabel):
        text = "{0}:{1}".format(constants.MPLS_TYPES, nulllabel)
        super(MpldNullLabel, self).__init__(constants.MPLS_NULL_LABEL, text=text, nsmap=self._NSMAP)

class Mpls(elements.BaseElement):
    """
    Mpls element
    """

    _FIELDS = ("_global", "signaling_protocols")

    def __init__(self):
        super(Mpls, self).__init__(constants.MPLS)
        self._global = MplsGlobal()
        self.signaling_protocols = MplsSignalingProtocols()


class MplsGlobal(elements.BaseElement):
    """
    Mpls global element
    """

    _FIELDS = ("config", "interface_attributes")

    def __init__(self):
        super(MplsGlobal, self).__init__(constants.MPLS_GLOBAL)
        self.config = MplsGlobalConfig()
        self.interface_attributes = elements.DictElement("interface-attributes")


class MplsGlobalConfig(elements.BaseElement):
    """
    MPLS global config element
    """

    _FIELDS = ("null_label",)

    def __init__(self):
        super(MplsGlobalConfig, self).__init__(constants.CONFIG)
        self.null_label = None

    def set_null_label(self, labe_type):
        """
        labe_type: EXPLICIT, IMPLICIT
        """
        if labe_type is None:
            self.null_label = None
        else:
            self.null_label = MpldNullLabel(labe_type)


class MplsGlobalInterface(elements.BaseElement):
    """
    MPLS global interface atribute element
    """

    _FIELDS = ("interface_id", "config", "interface_ref")

    def __init__(self, iface_id, interface, subinterface=0):
        super(MplsGlobalInterface, self).__init__(constants.INTERFACE)
        self.interface_id = iface_id
        self.config = MplsGlobalInterfaceConfig(iface_id)
        self.interface_ref = interface_ref.InterfaceRef(interface, subinterface)


class MplsGlobalInterfaceConfig(elements.BaseElement):
    """
    MPLS global interface atribute config element
    """

    _FIELDS = ("interface_id",)

    def __init__(self, iface_id):
        super(MplsGlobalInterfaceConfig, self).__init__(constants.CONFIG)
        self.interface_id = iface_id


class MplsSignalingProtocols(elements.BaseElement):
    """
    MPLS signaling protocols element
    """

    _FIELDS = ("ldp",)

    def __init__(self):
        super(MplsSignalingProtocols, self).__init__("signaling-protocols")
        self.ldp = MplsLdp()


class MplsLdp(elements.BaseElement):
    """
    MPLS LDP element
    """

    _FIELDS = ("_global",)
    def __init__(self):
        super(MplsLdp, self).__init__(constants.MPLS_LDP)
        self._global = MplsLdpGlobal()


class MplsLdpGlobal(elements.BaseElement):
    """
    MPLS LDP global element
    """

    _FIELDS = ("config", "address_families", "discovery")

    def __init__(self):
        super(MplsLdpGlobal, self).__init__(constants.MPLS_LDP_GLOBAL)
        self.config = MplsLdpGlobalConfig()
        self.address_families = MplsLdpGlobalAddrFamilies()
        self.discovery = MplsLdpGlobalDiscovery()


class MplsLdpGlobalConfig(elements.BaseElement):
    """
    MPLS LDP global config element
    """

    _FIELDS = ("lsr_id",)

    def __init__(self, lsr_id=None):
        super(MplsLdpGlobalConfig, self).__init__(constants.CONFIG)
        self.lsr_id = lsr_id


class MplsLdpGlobalAddrFamilies(elements.BaseElement):
    """
    MPLS LDP global address-families element
    """

    _FIELDS = ("ipv4",)

    def __init__(self):
        super(MplsLdpGlobalAddrFamilies, self).__init__("address-families")
        self.ipv4 = MplsLdpGlobalAddrIPv4()


class MplsLdpGlobalAddrIPv4(elements.BaseElement):
    """
    MPLS LDP global address-families(ipv4) element
    """

    _FIELDS = ("config",)

    def __init__(self):
        super(MplsLdpGlobalAddrIPv4, self).__init__("ipv4")
        self.config = MplsLdpGlobalAddrIPv4Config()


class MplsLdpGlobalAddrIPv4Config(elements.BaseElement):
    """
    MPLS LDP global address-families(ipv4) config element
    """

    _FIELDS = ("transport_address", "session_ka_holdtime", "label_policy")

    def __init__(self, transport_address=None, session_ka_holdtime=None):
        super(MplsLdpGlobalAddrIPv4Config, self).__init__(constants.CONFIG)
        self.transport_address = transport_address
        self.session_ka_holdtime = session_ka_holdtime
        self.label_policy = MplsLdpGlobalAddrIPv4ConfigLabelPolcy()


class MplsLdpGlobalAddrIPv4ConfigLabelPolcy(elements.BaseElement):
    """
    MPLS LDP global address-families(ipv4) config label-policy element
    """

    _FIELDS = ("advertise",)

    def __init__(self):
        super(MplsLdpGlobalAddrIPv4ConfigLabelPolcy, self).__init__("label-policy")
        self.advertise = MplsLdpGlobalAddrIPv4ConfigLabelPolcyAdv()


class MplsLdpGlobalAddrIPv4ConfigLabelPolcyAdv(elements.BaseElement):
    """
    MPLS LDP global address-families(ipv4) config label-policy advertise element
    """

    _FIELDS = ("egress_explicit_null",)

    def __init__(self):
        super(MplsLdpGlobalAddrIPv4ConfigLabelPolcyAdv, self).__init__("advertise")
        self.egress_explicit_null = MplsLdpGlobalAddrIPv4ConfigLabelPolcyAdvEngressExpNull()

class MplsLdpGlobalAddrIPv4ConfigLabelPolcyAdvEngressExpNull(elements.BaseElement):
    """
    MPLS LDP global address-families(ipv4) config label-policy advertise exp-null element
    """

    _FIELDS = ("enable",)

    def __init__(self, enable=True):
        super(
            MplsLdpGlobalAddrIPv4ConfigLabelPolcyAdvEngressExpNull,
            self).__init__("egress-explicit-null")
        self.enable = enable


class MplsLdpGlobalDiscovery(elements.BaseElement):
    """
    MPLS LDP global discovery element
    """

    _FIELDS = ("interfaces",)
    def __init__(self):
        super(MplsLdpGlobalDiscovery, self).__init__("discovery")
        self.interfaces = MplsLdpGlobalDiscoveryInterfaces()


class MplsLdpGlobalDiscoveryInterfaces(elements.DictElement):
    """
    MPLS LDP global discovery interfaces element
    """

    def __init__(self):
        super(MplsLdpGlobalDiscoveryInterfaces, self).__init__("interfaces")
        self.config = MplsLdpGlobalDiscoveryInterfacesConfig()

    def xml_element(self):
        elm = super(MplsLdpGlobalDiscoveryInterfaces, self).xml_element()
        elm.append(self.config.xml_element())
        return elm


class MplsLdpGlobalDiscoveryInterfacesConfig(elements.BaseElement):
    """
    MPLS LDP global discovery interfaces config element
    """

    _FIELDS = ("hello_holdtime", "hello_interval")
    def __init__(self, hello_holdtime=None, hello_interval=None):
        super(MplsLdpGlobalDiscoveryInterfacesConfig, self).__init__(constants.CONFIG)
        self.hello_holdtime = hello_holdtime
        self.hello_interval = hello_interval


class MplsLdpGlobalDiscoveryInterface(elements.BaseElement):
    """
    MPLS LDP global discovery interface element
    """

    _FIELDS = ("interface_id", "config", "address_families", "interface_ref")

    def __init__(self, iface_id, interface=None, subinterface=None):
        super(MplsLdpGlobalDiscoveryInterface, self).__init__("interface")
        self.interface_id = iface_id
        self.config = MplsLdpGlobalDiscoveryInterfaceConfig(iface_id)
        self.address_families = MplsLdpGlobalDiscoveryInterfaceAddrFamilies()
        self.interface_ref = interface_ref.InterfaceRef(interface, subinterface)


class MplsLdpGlobalDiscoveryInterfaceConfig(elements.BaseElement):
    """
    MPLS LDP global discovery interface config element
    """

    _FIELDS = ("interface_id", "hello_holdtime", "hello_interval")

    def __init__(self, iface_id, hello_holdtime=None, hello_interval=None):
        super(MplsLdpGlobalDiscoveryInterfaceConfig, self).__init__(constants.CONFIG)
        self.interface_id = iface_id
        self.hello_holdtime = hello_holdtime
        self.hello_interval = hello_interval


class MplsLdpGlobalDiscoveryInterfaceAddrFamilies(elements.BaseElement):
    """
    MPLS LDP global discovery interface address-fmilies element
    """

    _FIELDS = ("ipv4",)

    def __init__(self):
        super(MplsLdpGlobalDiscoveryInterfaceAddrFamilies, self).__init__("address-families")
        self.ipv4 = MplsLdpGlobalDiscoveryInterfaceAddrIPv4()


class MplsLdpGlobalDiscoveryInterfaceAddrIPv4(elements.BaseElement):
    """
    MPLS LDP global discovery interface address-fmilies(IPv4) element
    """

    _FIELDS = ("config",)

    def __init__(self):
        super(MplsLdpGlobalDiscoveryInterfaceAddrIPv4, self).__init__("ipv4")
        self.config = MplsLdpGlobalDiscoveryInterfaceAddrIPv4Config()


class MplsLdpGlobalDiscoveryInterfaceAddrIPv4Config(elements.BaseElement):
    """
    MPLS LDP global discovery interface address-fmilies(IPv4) config element
    """

    _FIELDS = ("enable",)

    def __init__(self, enable=True):
        super(MplsLdpGlobalDiscoveryInterfaceAddrIPv4Config, self).__init__(constants.CONFIG)
        self.enable = enable
