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

package main

import (
	"netconf/lib/sysrepo"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ConfigServer struct {
	Viper *viper.Viper
	DpIds map[DatapathId]struct{}
}

func NewConfigServer(dpids map[DatapathId]struct{}) *ConfigServer {
	return &ConfigServer{
		Viper: viper.New(),
		DpIds: dpids,
	}
}

func (s *ConfigServer) Serve(cfgPath string, cfgType string) error {
	s.Viper.SetConfigFile(cfgPath)
	s.Viper.SetConfigType(cfgType)
	if err := s.Viper.ReadInConfig(); err != nil {
		return err
	}

	ch := make(chan struct{})

	s.Viper.WatchConfig()
	s.Viper.OnConfigChange(func(fsnotify.Event) {
		ch <- struct{}{}
	})

	for {
		s.Update()
		<-ch
	}
}

func (s *ConfigServer) HasDpid(dpid DatapathId) bool {
	if len(s.DpIds) == 0 {
		return true
	}

	_, ok := s.DpIds[dpid]
	return ok
}

func (s *ConfigServer) Update() {
	log.Debugf("Update")

	conn := srlib.NewSrConnection()
	if err := conn.Connect("ncmi"); err != nil {
		log.Errorf("%s", err)
		return
	}
	defer conn.Disconnect()

	session := srlib.NewSrSession(nil)
	if err := session.Start(conn, srlib.SR_DS_RUNNING); err != nil {
		log.Errorf("%s", err)
		return
	}
	defer session.Stop()

	cfg := &Config{}
	s.Viper.Unmarshal(&cfg)
	for dpid, ports := range cfg.Dps {
		if ok := s.HasDpid(dpid); !ok {
			log.Debugf("dpid(%d) ignored.", dpid)
			continue
		}

		for _, p := range ports.Ports {
			log.Infof("update interfaces. dpid=%d, %v", dpid, p)
		}

		vals, err := NewInterfaceFactory("eth%d").NewInterfacesVals(ports.Ports)
		if err != nil {
			log.Warnf("NewInterfacesVals error. %s", err)
			continue
		}

		for _, val := range vals {
			log.Debugf("Set Item %s", val)
			if err := session.SetItem(val, srlib.SR_EDIT_DEFAULT); err != nil {
				log.Errorf("SetItem error. %s", err)
				return
			}
		}
	}

	if err := session.Commit(); err != nil {
		log.Errorf("Commit error. %s", err)
	}
}
