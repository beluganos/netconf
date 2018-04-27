#! /bin/bash
# -*- coding: utf-8; mode: shell-script -*-

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


if [ -z "${NC_HOME}" ]; then
    NC_BIN=/usr/bin
else
    NC_BIN=$NC_HOME/bin
fi

BIN_DIR="${NC_BIN}"
SRC_DIR="${NC_HOME}/etc/lxcinit"
DST_DIR="/tmp"

BIN_FILES="cfgd cfgcp cfgnet cfgfrr cfgsysctl cfgbgp netplan+"

do_usage() {
    echo "$0 <containe name> <continer type>"
}

err_exit() {
    local MSG=$1
    echo ${MSG}
    logger ${MSG}
    exit 1
}

do_init() {
    local LXC_NAME=$1
    local LXC_TYPE=$2
    local LXC_SRC="${SRC_DIR}/${LXC_TYPE}"
    local LXC_DST="${DST_DIR}/${LXC_TYPE}"

    # copy BIN_FILES to container.
    local BIN_NAME
    for BIN_NAME in ${BIN_FILES}; do
        echo "'${BIN_DIR}/${BIN_NAME}' -> '${LXC_NAME}/usr/bin/'"
        lxc file push ${BIN_DIR}/${BIN_NAME}  ${LXC_NAME}/usr/bin/ || err_exit "[${LXC_NAME}] copy bin error."
    done

    # run lxcinit.sh on local.
    ${LXC_SRC}/lxcinit.sh ${LXC_NAME} ${LXC_SRC} "local"

    # copy LXC_SRC/* to container.
    lxc file push -r -p ${LXC_SRC} ${LXC_NAME}${DST_DIR} || err_exit "[${LXC_NAME}] copy files error."

    # run lxcinit.sh on container.
    lxc exec ${LXC_NAME} ${LXC_DST}/lxcinit.sh ${LXC_NAME} ${LXC_DST} || err_exit "[${LXC_NAME}] lxcinit error."

}

if [ $# -ne 2 ]; then
    do_usage
    exit 1
fi

do_init $1 $2

exit 0
