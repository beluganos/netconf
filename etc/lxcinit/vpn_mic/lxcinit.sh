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

set_mic_name() {
    if [ -z "${NC_HOME}" ]; then
	MICNAME_DIR=/etc/beluganos
    else
	MICNAME_DIR=$NC_HOME/etc/test
    fi

    MICNAME_FILE=$MICNAME_DIR/mic_name
}

save_mic_name() {
    mkdir -v -p $MICNAME_DIR
    echo "$LXC_NAME" > $MICNAME_FILE

    echo "save mic-name '$LXC_NAME' to $MICNAME_FILE"
}

fix_mic_name() {
    local FILE_NAME=$1
    local TEMP_NAME=$FILE_NAME.bak
    local MIC_NAME=$LXC_NAME

    cp -v $FILE_NAME $TEMP_NAME

    sed s/%MIC_NAME%/$MIC_NAME/g $TEMP_NAME > $FILE_NAME
    echo "mic-name $MIC_NAME $FILE_NAME"

    rm -v -f $TEMP_NAME
}

do_init() {
    local BEL_USER="beluganos"
    local FRR_USER="frr"
    local SERVICES="beluganos.service nlad.service ribcd.service ribpd.service beluganos.target gobgpd.service cfgd.service netplan-ext.service vrf.service"

    # add user and create directory for beluganos.
    adduser --system --no-create-home --group ${BEL_USER}
    mkdir -p /etc/beluganos
    chown ${BEL_USER}:${BEL_USER} /etc/beluganos

    # copy config files
    install -v -m 0644 ./conf/sysctl.conf /etc/sysctl.d/30-beluganos.conf
    install -v -m 0644 -o ${FRR_USER} -g ${FRR_USER} ./conf/daemons     /etc/frr/daemons
    install -v -m 0644 -o ${FRR_USER} -g ${FRR_USER} ./conf/gobgpd.conf /etc/frr/gobgpd.toml
    install -v -m 0644 -o ${FRR_USER} -g ${FRR_USER} ./conf/gobgp.conf  /etc/frr/gobgp.conf
    install -v -m 0644 -o ${BEL_USER} -g ${BEL_USER} ./conf/ribxd.conf  /etc/beluganos/ribxd.conf
    install -v -m 0644 -o ${BEL_USER} -g ${BEL_USER} ./conf/vrf.conf    /etc/vrf.conf

    fix_mic_name /etc/beluganos/ribxd.conf

    # create frr.conf and restart frr
    touch /etc/frr/frr.conf
    systemctl restart frr

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
    local BEL_BIN_HOME
    if [ -z "${NC_HOME}" ]; then
        BEL_BIN_HOME=/usr/bin
    else
        BEL_BIN_HOME=$HOME/go/bin
    fi

    local BEL_BINS="nlad nlac ribcd ribpd ribsd ribsdmp gobgpd gobgp"
    local BEL_BIN
    for BEL_BIN in ${BEL_BINS}; do
        echo "'${BEL_BIN_HOME}/${BEL_BIN}' -> '${LXC_NAME}/usr/bin/'"
        lxc file push ${BEL_BIN_HOME}/${BEL_BIN} ${LXC_NAME}/usr/bin/
    done

    save_mic_name
}

_main() {
    echo "[lxcinit] START: $LXC_NAME/$WORK_DIR $RUN_MODE"
    cd $WORK_DIR

    set_mic_name

    if [ "$RUN_MODE" = "local" ]; then
        do_local
    else
        do_init
    fi

    exit 0
}

_main
