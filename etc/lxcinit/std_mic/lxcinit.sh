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

LXC_NAME=$1
WORK_DIR=$2
RUN_MODE=$3

do_init() {
    local BEL_USER="beluganos"
    local FRR_USER="frr"
    local SERVICES="cfgd.service netplan-ext.service"

    # copy config files
    install -v -m 0644 -o ${FRR_USER} -g ${FRR_USER} ./conf/daemons /etc/frr/daemons

    # copy service files.
    local SERVICE
    for SERVICE in ${SERVICES}; do
        install -v -m 0644 ./service/${SERVICE} /etc/systemd/system/${SERVICE}
    done

    # enable and start services.
    systemctl daemon-reload
    for SERVICE in ${SERVICES}; do
        systemctl enable ${SERVICE}
        systemctl start  ${SERVICE}
        echo "${SERVICE} started."
    done
}

do_local() {
    echo "success"
}

echo "[lxcinit] START: $LXC_NAME/$WORK_DIR $RUN_MODE"
cd $WORK_DIR

if [ "$RUN_MODE" = "local" ]; then
    do_local
else
    do_init
fi

exit 0
