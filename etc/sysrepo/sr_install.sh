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

WORK_HOME=`pwd`

# CMAKE_BUILD="-DCMAKE_BUILD_TYPE=Debug"
CMAKE_BUILD="-DCMAKE_BUILD_TYPE=Release"
CMAKE_PREFIX="-DCMAKE_INSTALL_PREFIX:PATH=/usr"
CMAKE_OPTS="${CMAKE_BUILD} ${CMAKE_PREFIX}"

LIBYANG_VER=v0.13-r2
SYSREPO_VER=v0.7.1
LIBNC2_VER=v0.9-r1
NETOP2_VER=v0.4-r1

MAKEALL_OPTS="-j5"

apt_install() {
    sudo -E apt -y install git cmake build-essential bison flex libpcre3-dev libev-dev libavl-dev libprotobuf-c-dev protobuf-c-compiler libpcre3-dev python-dev lua5.2 swig python-setuptools libssh-dev
}

install_cmocka() {
    pushd ${WORK_HOME}

    if [ ! -e cmocka ]; then
        git clone https://git.cryptomilk.org/projects/cmocka.git
    fi

    cd cmocka
    mkdir build; cd build
    cmake ${CMAKE_OPTS} ..
    make clean
    make ${MAKEALL_OPTS}
    sudo make install

    popd
}

install_libyang() {
    pushd ${WORK_HOME}

    if [ ! -e libyang ]; then
        git clone https://github.com/CESNET/libyang.git
        pushd libyang
        git checkout -b ${LIBYANG_VER} ${LIBYANG_VER}
        popd
    fi

    cd libyang
    mkdir build; cd build
    cmake ${CMAKE_OPTS} ..
    make clean
    make ${MAKEALL_OPTS}
    sudo make install

    popd
}

install_sysrepo() {
    pushd ${WORK_HOME}

    if [ ! -e sysrepo ]; then
        git clone https://github.com/sysrepo/sysrepo.git
        pushd sysrepo
        git checkout -b ${SYSREPO_VER} ${SYSREPO_VER}
        popd
    fi

    cd sysrepo
    mkdir build; cd build
    cmake ${CMAKE_OPTS} -DWITH_SYSTEMD=yes -DREPOSITORY_LOC:PATH=/etc/sysrepo ..
    make clean
    make ${MAKEALL_OPTS}
    sudo make install

    popd
}

install_libnetconf2() {

    pushd ${WORK_HOME}

    if [ ! -e libnetconf2 ]; then
        git clone https://github.com/CESNET/libnetconf2.git
        pushd libnetconf2
        git checkout -b ${LIBNC2_VER} ${LIBNC2_VER}
        popd
    fi

    cd libnetconf2
    mkdir build; cd build
    cmake ${CMAKE_OPTS} ..
    make clean
    make ${MAKEALL_OPTS}
    sudo make install

    popd
}

install_netop2() {
    pushd ${WORK_HOME}

    if [ ! -e Netopper2 ]; then
        git clone https://github.com/CESNET/Netopeer2.git
        pushd Netopeer2
        git checkout -b ${NETOP2_VER} ${NETOP2_VER}
        popd
    fi

    pushd Netopeer2/server
    mkdir build; cd build
    cmake ${CMAKE_OPTS} ..
    make clean
    make ${MAKEALL_OPTS}
    sudo make install
    popd

    pushd Netopeer2/cli
    mkdir build; cd build
    cmake ${CMAKE_OPTS} ..
    make clean
    make ${MAKEALL_OPTS}
    sudo make install
    popd

    popd

    sudo install -v -C netopeer2-server.service /lib/systemd/system/
}

apply_system() {
    sudo ldconfig
    sudo systemctl daemon-reload
}

get_openconfig() {
    pushd ${WORK_HOME}

    if [ ! -e public ]; then
        git clone https://github.com/openconfig/public.git
        cd public
        mkdir all; cd all
        ln -s ../release/models/*.yang .
        ln -s ../release/models/*/*.yang .
    fi

    popd
}

get_beluganos_yang() {
    install -v -C ../openconfig/reload.sh ${WORK_HOME}/public/all/sryang.sh
    install -v -C ../openconfig/beluganos-*.yang ${WORK_HOME}/public/all/
    install -v -C ../openconfig/ietf-*.yang ${WORK_HOME}/public/all/
    install -v -C ../openconfig/iana-*.yang ${WORK_HOME}/public/all/
}

get_yangmodel() {
    mkdir -p ${WORK_HOME}
    pushd ${WORK_HOME}

    if [ ! -e yang ]; then
        git clone https://github.com/YangModels/yang.git
    fi

    popd
}

do_install() {
    apt_install
    install_cmocka
    install_libyang
    install_sysrepo
    install_libnetconf2
    install_netop2
    apply_system
}

do_yang() {
    get_openconfig
    get_beluganos_yang
    # get_yangmodel
}

do_usage() {
    echo "test"
}

case $1 in
    install)
        do_install
        ;;
    libyang)
        install_libyang
        ;;
    sysrepo)
        install_sysrepo
        ;;
    netconf)
        install_libnetconf2
        ;;
    netop2)
        install_netop2
        ;;
    yang)
        do_yang
        ;;
    *)
        do_usage
        ;;
esac
