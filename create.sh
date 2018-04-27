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

. ./create.ini

SYSREPO_INSTALL=yes

set_proxy() {
    if [ "${PROXY}"x != ""x ]; then
        APT_PROXY="--env http_proxy=${PROXY}"
        HTTP_PROXY="http_proxy=${PROXY} https_proxy=${PROXY}"
        export http_proxy=${PROXY}
        export https_proxy=${PROXY}

        echo "--- Proxy settings ---"
        echo "APT_PROXY=${APT_PROXY}"
        echo "HTTP_PROXY=${HTTP_PROXY}"
    fi
}

#
# install deb packages
#
apt_install() {
    sudo ${HTTP_PROXY} apt -y install ${APT_PKGS} || { echo "apt_install error."; exit 1; }
    sudo apt -y autoremove
}

#
# install go-lang
#
golang_install() {
    local GO_FILE=go${GO_VER}.linux-amd64.tar.gz

    echo "Downloading ${GO_URL}/${GO_FILE}"
    wget -nc -P /tmp ${GO_URL}/${GO_FILE} || { echo "golang_install/wget error."; exit 1; }

    echo "Extracting /tmp/${GO_FILE}"
    sudo tar xf /tmp/${GO_FILE} -C /usr/local || { echo "golang_install/tar error."; exit 1; }
}

#
# install protobuf
#
protoc_install() {
    local PROTOC_FILE=protoc-${PROTOC_VER}-linux-x86_64.zip

    echo "Downloading ${PROTOC_URL}/${PROTOC_FILE}"
    wget -nc -P /tmp ${PROTOC_URL}/${PROTOC_FILE} || { echo "protoc_install/wget error."; exit 1; }

    echo "Extracting /tmp/${PROTOC_FILE}"
    sudo unzip -o -d /usr/local/go /tmp/${PROTOC_FILE} || { echo "protoc_install/unzip error."; exit 1; }

    sudo chmod +x /usr/local/go/bin/protoc
}

#
# install go packages
#
gopkg_install() {
    local PKG
    for PKG in ${GO_PKGS}; do
        echo "go get ${PKG}"
        go get -u ${PKG} || { echo "gopkg_install error."; exit 1; }
    done
}

sysrepo_install() {
    if [ "$SYSREPO_INSTALL" = "yes" ]; then
        pushd etc/sysrepo

        sudo HTTPS_PROXY=${PROXY} ./sr_install.sh install
        HTTPS_PROXY=${PROXY} ./sr_install.sh yang

        sudo sysrepod
        cd public/all
        ./sryang.sh install-all

        sleep 3
        sudo pkill sysrepod

        popd
    fi
}

beluganos_netconf_install() {
    # enable go-env
    . ./setenv.sh

    # install packages
    gopkg_install
    sysrepo_install

    ./bootstrap.sh
    make release
    # sudo make install-services
}

confirm() {
    local MSG=$1

    echo "$MSG [y/N]"
    read ans
    case $ans in
        [yY]) return 0;;
        *) return 1;;
    esac
}

do_all() {
    confirm "Install ALL" || exit 1

    # install base packages and tool
    apt_install
    golang_install
    protoc_install
    beluganos_netconf_install
}

do_usage() {
    echo "Usage $0 [OPTIONS]"
    echo "Options:"
    echo "  ''    : run all"
    echo "  pkg   : update apt-packages."
    echo "  golang: install or update golang and protobuf."
    echo "  gopkg : update go-packages"
    echo "  help  : show this message"
}


set_proxy
case $1 in
    pkg)
        apt_install
        ;;
    golang)
        golang_install
        protoc_install
        ;;
    gopkg)
        gopkg_install
        ;;
    beluganos-netconf)
        beluganos_netconf_install
        ;;
    sysrepo)
        sysrepo_install
        ;;
    help)
        do_usage
        ;;
    *)
        do_all
        ;;
esac
