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

package ncm

import (
	"fmt"
	"netconf/lib"
	"netconf/lib/gobgp/openconfig"
	"netconf/lib/net"
	"netconf/lib/openconfig"
	"netconf/lib/sysrepo"

	log "github.com/sirupsen/logrus"
)

func copyConfigAsync(session *srlib.SrSession, module string) {
	go func() {
		err := session.CopyConfig(module, srlib.SR_DS_RUNNING, srlib.SR_DS_STARTUP)
		if err != nil {
			log.Errorf("CopyConfig(%s) error. %s", module, err)
		} else {
			log.Infof("CopyConfig(%s) success.", module)
		}
	}()
}

func getLxcType(config *openconfig.NetworkInstanceConfig) (string, error) {
	std_or_vpn := func() string {
		if config.RT.Type() == ncnet.RD_TYPE_NONE {
			return "std"
		}
		return "vpn"
	}()

	switch config.Type {
	case openconfig.NETWORK_INSTANCE_DEFAULT:
		return fmt.Sprintf("%s_mic", std_or_vpn), nil

	case openconfig.NETWORK_INSTANCE_L3VRF:
		return fmt.Sprintf("%s_ric", std_or_vpn), nil

	default:
		return "", fmt.Errorf("Unsupported network-instance-type %s", config.Type)
	}
}

type NetworkInstancesSet map[srlib.SrChangeOper]openconfig.NetworkInstances

func NewNetworkInstancesSet() NetworkInstancesSet {
	return NetworkInstancesSet{
		srlib.SR_OP_CREATED:  openconfig.NewNetworkInstances(),
		srlib.SR_OP_MODIFIED: openconfig.NewNetworkInstances(),
		srlib.SR_OP_DELETED:  openconfig.NewNetworkInstances(),
	}
}

func (s NetworkInstancesSet) Unmarshall(cv *srlib.SrChangeVal) error {
	return cv.Dispatch(
		s[srlib.SR_OP_CREATED],
		s[srlib.SR_OP_MODIFIED],
		s[srlib.SR_OP_DELETED],
	)
}

func (s NetworkInstancesSet) Walk(oper srlib.SrChangeOper, f func(string, *openconfig.NetworkInstance) error) error {
	if nis, ok := s[oper]; ok {
		for name, ni := range nis {
			if err := f(name, ni); err != nil {
				return err
			}
		}
	}

	return nil
}

type NIUpdateType int

const (
	NI_UPDATE_TYPE NIUpdateType = iota
	NI_UPDATE_VTY
	NI_UPDATE_SYSCTL
	NI_UPDATE_SYSVRF
	NI_UPDATE_NETWORK
	NI_UPDATE_GOBGP
	NI_UPDATE_GOBGP_CFG
)

type NIUpdates map[NIUpdateType]int

func (n NIUpdates) Clear() {
	for t, _ := range n {
		n.SetUpdate(t, -1)
	}
}

func (n NIUpdates) SetUpdate(t NIUpdateType, pos int) {
	n[t] = pos
}

func (n NIUpdates) GetUpdate(t NIUpdateType) int {
	if b, ok := n[t]; ok {
		return b
	}
	return -1
}

type NICommands struct {
	Cmds     *nclib.Commands
	Upds     NIUpdates
	Bgps     *srocgobgp.ConfigProcessor
	NoCommit bool
}

func NewNICommands(ev srlib.SrNotifEvent, oper srlib.SrChangeOper) *NICommands {
	cmds := nclib.NewCommands(func(act nclib.CommandAction, cmd nclib.Command, ret []byte) {
		log.Debugf("NI/%s/%s/%s %s", ev, oper, act, cmd.Line(act))
		log.Debugf("NI/%s/%s/%s %s", ev, oper, act, string(ret))
	})

	return &NICommands{
		Cmds:     cmds,
		Upds:     NIUpdates{},
		Bgps:     srocgobgp.NewConfigProcessor(),
		NoCommit: false,
	}
}

type NICommandsHandler interface {
	AddCmd(*nclib.Shell, *nclib.Shell, *nclib.Shell)
	OnceCmd(NIUpdateType, *nclib.Shell, *nclib.Shell, *nclib.Shell)
	SetCmd(NIUpdateType, *nclib.Shell, *nclib.Shell, *nclib.Shell)
}

func (n *NICommands) Clear() {
	n.Cmds.Clear()
	n.Upds.Clear()
	n.Bgps.Clear()
	n.NoCommit = false
}

func (n *NICommands) AddCmd(do, undo, end *nclib.Shell) {
	n.Cmds.Add(nclib.NewShellCommand(do, undo, end))
}

func (n *NICommands) OnceCmd(up NIUpdateType, do, undo, end *nclib.Shell) {
	if n.Upds.GetUpdate(up) < 0 {
		n.Upds.SetUpdate(up, n.Cmds.Size())
		n.Cmds.Add(nclib.NewShellCommand(do, undo, end))
	}
}

func (n *NICommands) SetCmd(t NIUpdateType, do, undo, end *nclib.Shell) {
	cmd := nclib.NewShellCommand(do, undo, end)
	if pos := n.Upds.GetUpdate(t); pos < 0 {
		n.Upds.SetUpdate(t, n.Cmds.Size())
		n.Cmds.Add(cmd)
	} else {
		n.Cmds.Set(cmd, pos)
	}
}
