#! /bin/bash
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

MODULES=(
    ietf-interfaces
    ietf-ip
    iana-if-type
    openconfig-extensions
    openconfig-types
    openconfig-yang-types
    openconfig-inet-types
    openconfig-mpls-types
    openconfig-bgp-types
    openconfig-ospf-types
    openconfig-policy-types
    openconfig-network-instance-types
    beluganos-interfaces
    beluganos-if-ip
    beluganos-if-ethernet
    beluganos-mpls-ldp
    beluganos-mpls
    beluganos-bgp
    beluganos-ospfv2
    beluganos-network-instance
    beluganos-bgp-policy
)

do_show() {
    sysrepoctl -l
}

do_install() {
    local MODULE=$1
    sudo sysrepoctl -i -g ${MODULE}.yang
}

do_uninstall() {
    local MODULE=$1
    sudo sysrepoctl -u -m ${MODULE}
}

do_check() {
    local MODULE=$1
    pyang ${MODULE}.yang
}

do_sample() {
    local MODULE=$1
    pyang ${MODULE}.yang -f tree > ${MODULE}.txt
    pyang ${MODULE}.yang -f sample-xml-skeleton > ${MODULE}.xml
}

do_install_all() {
    for module in ${MODULES[@]}; do
        do_install $module
    done
}

do_uninstall_all() {
    for (( i=${#MODULES[@]}-1 ; i>=0 ; i--)) ; do
        do_uninstall ${MODULES[$i]}
    done
}

case $1 in
    install)
        do_install $2
        do_show
        ;;
    uninstall)
        do_uninstall $2
        do_show
        ;;
    reinstall)
        do_uninstall $2
        do_install $2
        do_show
        ;;
    check)
        do_check $2
        ;;
    sample)
        do_sample $2
        ;;
    show)
        do_show
        ;;
    install-all)
        do_install_all
        do_show
        ;;
    uninstall-all)
        do_uninstall_all
        do_show
        ;;
    reinstall-all)
        do_uninstall_all
        do_install_all
        do_show
        ;;
esac
