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

"""
XML for beluganos-routing-policy builder.
"""

from lxml import etree as ET

from beluganos.yang.python.elements import ListElement, Element
from beluganos.yang.python import constants
from beluganos.yang.python.rpol import RoutingPolicy
from beluganos.yang.python.rpol import PolicyDefinition
from beluganos.yang.python.rpol_stmt import PolicyStatement
from beluganos.yang.python.bgp_policy import BgpActions

def _get_sub_dict(d, *names):
    """
    get sub dictionary recursive.
    """
    for name in names:
        d = d.get(name, None)
        if d is None:
            return dict()
    return d


def new_rpol_stmt(stmtname, stmtcfg):
    stmt = PolicyStatement(stmtname)

    actcfg = stmtcfg.get("actions", None)
    if actcfg:
        stmt.actions.config.policy_result = actcfg.get("policy-result", None)

        bgpcfg = actcfg.get("bgp", None)
        if bgpcfg:
            act = BgpActions()
            act.config.set_local_pref = bgpcfg.get("set-local-pref")
            act.config.set_next_hop = bgpcfg.get("set-next-hop")
            stmt.actions.action = act

    return stmt

def new_rpol_def(pdefname, pdefcfg):
    defs = PolicyDefinition(pdefname)
    stmtcfgs = pdefcfg.get("stmts", dict())
    for stmtname, stmtcfg in stmtcfgs.items():
        defs.statements.append(new_rpol_stmt(stmtname, stmtcfg))
    return defs

def new_rpol(cfg):
    rpol = RoutingPolicy()

    pdefcfgs = cfg.get("pdefs", dict())
    for pdefname, pdefcfg in pdefcfgs.items():
        rpol.policy_definitions.append(new_rpol_def(pdefname, pdefcfg))

    return rpol


def read_config(path):
    """
    Read yaml file and return dict.
    """
    import yaml
    with open(path, "r") as cfg:
        return yaml.load(cfg)


def read_configs(paths):
    """
    ead yaml files and merge to dict.
    """
    pdefs = dict()
    for path in paths:
        cfg = read_config(path)
        pdefscfg = cfg.get("policy-definitions", dict())

        for defname, defcfg in pdefscfg.items():
            pdefs[defname] = defcfg

    return dict(pdefs=pdefs)


def _getopts():
    from optparse import OptionParser
    parser = OptionParser()
    parser.add_option("-v", "--verbose", dest="verbose", action="store_true",
                      help="show detail messages.")
    return parser.parse_args()


def _main():
    _, args = _getopts()

    cfg = read_configs(args)
    rpols = new_rpol(cfg)
    print ET.tostring(rpols.xml_element(), pretty_print=True)


if __name__ == "__main__":
    _main()
