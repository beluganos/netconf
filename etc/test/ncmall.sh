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

do_kill() {
    local pid=$1

    local child
    local children=`ps -o pid --no-headers --ppid ${pid}`
    for child in ${children}; do
        do_kill $child
    done

    # echo "kill pid=${pid}"
    kill -9 ${pid} &> /dev/null
}

do_stop() {
    do_kill $$
}

do_start() {
    ./ncmd.sh $* &
    ./ncms.sh $* &
    ./ncmi.sh $* &
}

trap "do_stop; exit 0" INT
do_start $*
wait
