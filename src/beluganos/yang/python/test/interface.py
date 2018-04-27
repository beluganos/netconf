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
from beluganos.yang.python import interface
from beluganos.yang.python import subinterface
from beluganos.yang.python import address

class InterfaceTest(unittest.TestCase):
    def test_interfaces(self):
        iface1 = interface.Interface("eth1")
        iface1.config.mtu = 100

        subif2 = subinterface.Subinterface(0)
        subif2_10 = subinterface.Subinterface(10)
        subif2_10.ipv4.addresses.append(address.Address("10.0.0.1", 24))
        iface2 = interface.Interface("eth2")
        iface2.subinterfaces.append(subif2, subif2_10)

        ifaces = elements.ListElement("interfaces")
        ifaces.append(iface1, iface2)

        print ET.tostring(ifaces.xml_element(), pretty_print=True)



if __name__ == "__main__":
    unittest.main()
