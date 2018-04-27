// -*- coding: utf-8 -*-

// Copyright (C) 2018 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#include "helper.h"

sr_val_t *get_val(sr_val_t *val, size_t i) {
  return &val[i];
}

extern int Go_module_change_cb(sr_session_ctx_t*, const char*, sr_notif_event_t, void*);

int module_change_cb(sr_session_ctx_t* s, const char* module_name, sr_notif_event_t e, void* p) {
  return Go_module_change_cb(s, module_name, e, p);
}


extern int Go_subtree_change_cb(sr_session_ctx_t*, const char*, sr_notif_event_t, void*);

int subtree_change_cb(sr_session_ctx_t* s, const char* xpath, sr_notif_event_t e, void* p) {
  return Go_subtree_change_cb(s, xpath, e, p);
}

extern void Go_log_cb(sr_log_level_t, const char*);

void log_cb(sr_log_level_t level, const char* message) {
  Go_log_cb(level, message);
}
