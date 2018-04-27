#! /usr/bin/env python
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


import unittest
from lxml import etree as ET
from beluganos.yang.python import elements
from beluganos.yang.python import network_instance as ni
from beluganos.yang.python import mpls
from beluganos.yang.python import bgp

class TestNetworkInstance(unittest.TestCase):
    def test_ni(self):
        ni1 = ni.NetworkInstance("pe1")
        ni1.config.router_id = "10.0.0.1"
        ni1.config.route_distinguisher = "10:100"
        ni1.config.route_target = "100:100"

        ni1.interfaces.append(
            ni.NetworkInstanceInterface("wth1.10", "eth1", 10),
            ni.NetworkInstanceInterface("wth2.10"),
        )

        ni1.mpls._global.interface_attributes.append(
            mpls.MplsGlobalInterface("eth2.10", "eth2", 10),
        )

        ni1.mpls.signaling_protocols.ldp._global.discovery.interfaces.append(
            mpls.MplsLdpGlobalDiscoveryInterface("eth1.10", "eth1", 10),
        )

        _bgp = bgp.Bgp("test-bgp")
        _bgp._global.config.as_number = 65000
        _bgp._global.config.router_id = "10.0.1.1"
        _bgp.neighbors.append(
            bgp.BgpNeighbor("10.0.0.1"),
        )
        ni1.protocols.append(ni.NetworkInstanceProtocol(_bgp))

        nis = elements.ListElement("network-instances")
        nis.append(ni1)

        print ET.tostring(nis.xml_element(), pretty_print=True)


if __name__ == "__main__":
    unittest.main()
