#! /usr/bin/env python

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

import setuptools
from distutils.core import setup

setup(name='beluganos-netconf',
      version='1.0',
      description='Beluganos netconf tool.',
      packages=setuptools.find_packages(),
      entry_points={
          "console_scripts":[
              "beluganos-interfaces = beluganos.yang.app.beluganos_interfaces:_main",
              "beluganos-network-instance = beluganos.yang.app.beluganos_network_instance:_main",
              "beluganos-routing-policy = beluganos.yang.app.beluganos_routing_policy:_main",
          ],
      },
)
