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

do_cfg() {
    local MODE=$1
    local INPUT=$2
    local NAME=$3

    echo "sysrepocfg ${MODE} ${INPUT} ${NAME}"
    sysrepocfg ${MODE} ${INPUT} ${NAME}
}

do_ni() {
    local MODE=$1
    local FILE=$2
    do_cfg ${MODE} ${FILE} beluganos-network-instance
}

do_iface() {
    local MODE=$1
    local FILE=$2
    do_cfg ${MODE} ${FILE} beluganos-interfaces
}

do_rpol() {
    local MODE=$1
    local FILE=$2
    do_cfg ${MODE} ${FILE} beluganos-routing-policy
}

show_cfg() {
    local NAME=$1

    echo "<!-- STA: ${NAME} -->"
    sysrepocfg -x - ${NAME}
    echo "<!-- END: ${NAME} -->"
}

show_ni() {
    show_cfg beluganos-network-instance
}

show_iface() {
    show_cfg beluganos-interfaces
}

show_rpol() {
    show_cfg beluganos-routing-policy
}

do_usage() {
    local NAME=$1
    echo "${NAME} [set/merge/show] [ni/iface/rpol] [file]"
}

case $1 in
    set)
        case $2 in
            ni)
                do_ni -i $3
                ;;
            iface)
                do_iface -i $3
                ;;
            rpol)
                do_rpol -i $3
                ;;
            *)
                do_usage $0
                ;;
        esac
        ;;

    merge)
        case $2 in
            ni)
                do_ni -m $3
                ;;
            iface)
                do_iface -m $3
                ;;
            rpol)
                do_rpol -m $3
                ;;
            *)
                do_usage $0
                ;;
        esac
        ;;

    show)
        case $2 in
            ni)
                show_ni
                ;;
            iface)
                show_iface
                ;;
            rpol)
                show_rpol
                ;;
            *)
                show_ni
                show_iface
                show_rpol
                ;;
        esac
        ;;

    *)
        do_usage $0
        ;;
esac
