#! /bin/bash
#! -*- coding: utf-8 -*-

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

. ./ncm.conf

echo "[ncmd]: ${NC_BINS}/ncmd"
echo "[ncmd]:   NC_HOME=${NC_HOME}"
echo "[ncmd]:   -c ${NCMD_CONF}"
echo "[ncmd]:   dry-run=${NCMD_DRYRUN}"
echo "[ncmd]:   $*"

export NC_HOME
${NC_BINS}/ncmd -c ${NCMD_CONF} ${NCMD_DRYRUN} $* 2>&1 | tee ${NC_LOGS}/ncmd.log
