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
contant values.
"""

NETWORK_INSTANCE_NS = "https://github.com/beluganos/beluganos/yang/network-instance"
NETWORK_INSTANCES = "network-instances"
NETWORK_INSTANCE = "network-instance"

NETWORK_INSTANCE_TYPES_NS = "http://openconfig.net/yang/network-instance-types"
NETWORK_INSTANCE_TYPES = "oc-ni-types"

NETWORK_INSTANCE_PROTOCOLS = "protocols"
NETWORK_INSTANCE_PROTOCOL = "protocol"
NETWORK_INSTANCE_PROTOCOL_IDENT = "identifier"
NETWORK_INSTANCE_PROTOCOL_TYPES = "oc-pol-types"
NETWORK_INSTANCE_PROTOCOL_TYPES_NS = "http://openconfig.net/yang/policy-types"

INTERFACES_NS = "https://github.com/beluganos/beluganos/yang/interfaces"
INTERFACES = "interfaces"
INTERFACE = "interface"

INTERFACES_ETH_NS = "https://github.com/beluganos/beluganos/yang/interfaces/ethernet"
INTERFACES_ETH = "ethernet"
INTERFACES_IP_NS = "https://github.com/beluganos/beluganos/yang/interfaces/ip"
INTERFACES_IP_V4 = "ipv4"
INTERFACES_IP_V6 = "ipv6"

SUBINTERFACES = "subinterfaces"
SUBINTERFACE = "subinterface"

INTERFACE_REF = "interface-ref"

ADDRESSES = "addresses"
ADDRESS = "address"

BGP = "bgp"
BGP_GLOBAL = "global"
BGP_NEIGHBORS = "neighbors"
BGP_NEIGHBOR = "neighbor"
BGP_ZEBRA = "zebra"
BGP_TYPES = "oc-bgp-types"
BGP_TYPES_NS = "http://openconfig.net/yang/bgp-types"
BGP_POLICY_NS = "https://github.com/beluganos/beluganos/yang/bgp-policy"

MPLS = "mpls"
MPLS_GLOBAL = "global"
MPLS_GLOBAL_IFATTRS = "interface-attributes"
MPLS_LDP = "ldp"
MPLS_LDP_GLOBAL = "global"

MPLS_TYPES_NS = "http://openconfig.net/yang/mpls-types"
MPLS_TYPES = "oc-mpls-types"
MPLS_NULL_LABEL = "null-label"

CONFIG = "config"

POLICY_PROTOCOL_BGP = "BGP"
POLICY_PROTOCOL_STATIC = "STATIC"
POLICY_PROTOCOL_OSPF = "OSPF"
POLICY_PROTOCOL_OSPF6 = "OSPF6"

ROUTING_POLICY_NS = "https://github.com/beluganos/beluganos/yang/routing-policy"

IANA_IFTYPE_NS = "urn:ietf:params:xml:ns:yang:iana-if-type"
IANA_IFTYPE = "iana-if-type"
