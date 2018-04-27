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
	"io"
	"netconf/app/ncm/cfg"
	"netconf/lib"
	"netconf/lib/openconfig"
	"netconf/lib/sysctl"
)

func cliConfig() *ncmcfg.CliConfig {
	return ncmcfg.GetConfig().Cli
}

func AddNIContainerCmd(h NICommandsHandler, name string, create bool) {
	cmd := cliConfig().LxdPath()
	arg := func(ope string) []string {
		return []string{"container", ope, name}
	}

	if create {
		h.AddCmd(
			nclib.NewShell(cmd, arg("create")...), // Do
			nclib.NewShell(cmd, arg("delete")...), // Undo
			nil, // End
		)

	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("delete")...), // Do
			nil, // Undo
			nil, // End
		)
	}
}

func AddNIContainterInitCmd(h NICommandsHandler, name string, config *openconfig.NetworkInstanceConfig) {
	cmd := cliConfig().LxcInitPath()
	lxcType, _ := getLxcType(config)
	h.AddCmd(
		nclib.NewShell(cmd, config.Name, lxcType), // Do
		nil, // Undo
		nil, // End
	)
}

func AddNIContainerInterfaceCmd(h NICommandsHandler, name string, ifname string, hwaddr string, add bool) {
	cmd := cliConfig().LxdPath()
	arg := func(ope string) []string {
		return []string{"interface", ope, name, ifname, hwaddr}
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg("add")...),    // Do
			nclib.NewShell(cmd, arg("delete")...), // Undo
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("delete")...), // Do
			nclib.NewShell(cmd, arg("add")...),    // Undo
			nil, // End
		)
	}
}

func AddNIInterfaceSysctlCmd(h NICommandsHandler, name string, ifname string, add bool) {

	AddNISysctlConfigCmd(h, name)

	cmd := cliConfig().SysPath()
	arg := func(flags ...string) []string {
		val := fmt.Sprintf("net.ipv4.conf.%s.rp_filter=0", ncsclib.FixString(ifname))
		flags = append(flags, "-H", name)
		return append([]string{"sysctl", "set", val}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...),     // Do
			nclib.NewShell(cmd, arg("-n")...), // Undo
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nclib.NewShell(cmd, arg()...),     // Undo
			nil, // End
		)
	}
}

func AddNIMplsInterfaceSysctlCmd(h NICommandsHandler, name string, ifname string, add bool) {

	AddNISysctlConfigCmd(h, name)

	cmd := cliConfig().SysPath()
	arg := func(flags ...string) []string {
		val := fmt.Sprintf("net.mpls.conf.%s.input=1", ncsclib.FixString(ifname))
		flags = append(flags, "-H", name)
		return append([]string{"sysctl", "set", val}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...),     // Do
			nclib.NewShell(cmd, arg("-n")...), // Undo
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nclib.NewShell(cmd, arg()...),     // Undo
			nil, // End
		)
	}
}

func AddNIInterfaceNetworkCmd(h NICommandsHandler, name string, device string, subif *openconfig.Subinterface, add bool) {

	AddNINetworkConfigCmd(h, name)

	cmd := cliConfig().SysPath()
	arg := func(flags ...string) []string {
		vid := fmt.Sprintf("%d", subif.Index)
		mtu := fmt.Sprintf("%d", subif.IPv4.Config.Mtu)
		flags = append(flags, "--mtu", mtu, "-H", name)
		return append([]string{"network", "set", "vlan", device, vid}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...),     // Do
			nclib.NewShell(cmd, arg("-n")...), // UnDo
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nclib.NewShell(cmd, arg()...),     // UnDo
			nil, // End
		)
	}
}

func AddNIRouterIdCmd(h NICommandsHandler, name string, routerId string, add bool) {

	AddNIVtyConfigCmd(h, name)

	cmd := cliConfig().VtyPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"global", "router-id", routerId}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIVrfCmd(h NICommandsHandler, name string, rd string, rt string, add bool) {

	AddNIVrfConfigCmd(h, name)

	cmd := cliConfig().SysPath()
	arg := func(flags ...string) []string {
		args := []string{"vrf", "set"}
		if len(rd) != 0 {
			args = append(args, fmt.Sprintf("RD=%s", rd))
		}
		if len(rt) != 0 {
			args = append(args, fmt.Sprintf("RT=%s", rt))
		}
		flags = append(flags, "-H", name)
		return append(args, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIVtyInterfaceCmd(h NICommandsHandler, name string, ifname string, key string, val interface{}, add bool) {

	AddNIVtyConfigCmd(h, name)

	cmd := cliConfig().VtyPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"interface", ifname, key, fmt.Sprintf("%v", val)}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIOspfRouterCmd(h NICommandsHandler, name string, key string, val interface{}, add bool) {
	AddNIVtyConfigCmd(h, name)

	cmd := cliConfig().VtyPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"ospf", key, fmt.Sprintf("%v", val)}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIMplsLdpCmd(h NICommandsHandler, name string, key string, val interface{}, add bool) {

	AddNIVtyConfigCmd(h, name)

	cmd := cliConfig().VtyPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"mpls", "ldp", key, fmt.Sprintf("%v", val)}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do,
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do,
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIMplsLdpIPv4Cmd(h NICommandsHandler, name string, key string, val interface{}, add bool) {
	AddNIVtyConfigCmd(h, name)

	cmd := cliConfig().VtyPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"mpls", "ldp", "address-family", "ipv4", key, fmt.Sprintf("%v", val)}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do,
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do,
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIMplsLdpIPv4IfaceCmd(h NICommandsHandler, name string, ifname string, key string, val interface{}, add bool) {

	AddNIVtyConfigCmd(h, name)

	cmd := cliConfig().VtyPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"mpls", "ldp", "address-family", "ipv4", "interface", ifname, key, fmt.Sprintf("%v", val)}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do,
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do,
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIStaticRouteCmd(h NICommandsHandler, name string, dest string, nexthop string, add bool) {

	AddNIVtyConfigCmd(h, name)

	cmd := cliConfig().VtyPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"ip", "route", dest, nexthop}, flags...)
	}

	if add {
		h.AddCmd(
			nclib.NewShell(cmd, arg()...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	} else {
		h.AddCmd(
			nclib.NewShell(cmd, arg("-n")...), // Do
			nil, // Undo (restart frr if failed.)
			nil, // End
		)
	}
}

func AddNIBgpConfigCmd(h NICommandsHandler, name string, cfgs io.Reader, restart bool, add bool) {
	cmd := cliConfig().GoBgpPath()
	arg := func(flags ...string) []string {
		flags = append(flags, "-H", name)
		return append([]string{"config", "set", "-"}, flags...)
	}

	if restart {
		h.SetCmd(NI_UPDATE_GOBGP,
			nclib.NewShell(cmd, "config", "backup", "-H", name),   // Do
			nclib.NewShell(cmd, "config", "rollback", "-H", name), // Undo
			nclib.NewShell(cmd, "config", "restart", "-H", name),  // End
		)
	} else {
		h.OnceCmd(NI_UPDATE_GOBGP,
			nclib.NewShell(cmd, "config", "backup", "-H", name),   // Do
			nclib.NewShell(cmd, "config", "rollback", "-H", name), // Undo
			nclib.NewShell(cmd, "config", "commit", "-H", name),   // End
		)
	}

	if add {
		h.SetCmd(
			NI_UPDATE_GOBGP_CFG,
			nclib.NewShellIn(cmd, cfgs, arg()...), // Do
			nil, // Undo
			nil, // End
		)
	} else {
		h.SetCmd(
			NI_UPDATE_GOBGP_CFG,
			nclib.NewShellIn(cmd, cfgs, arg("-n")...), // Do
			nil, // Undo
			nil, // End
		)
	}
}

func AddNIVtyConfigCmd(h NICommandsHandler, name string) {
	vtycmd := cliConfig().VtyPath()

	undoCmd := func() *nclib.Shell {
		syscmd := cliConfig().SysPath()
		switch ncmcfg.GetConfig().Frr.AutoRstart {
		case "restart":
			return nclib.NewShell(syscmd, "systemctl", "restart", "frr", "-H", name)
		case "reload":
			return nclib.NewShell(syscmd, "systemctl", "reload", "frr", "-H", name)
		default:
			return nil
		}
	}()

	h.OnceCmd(NI_UPDATE_VTY,
		nclib.NewShell(vtycmd, "config", "backup", "-H", name), // Do
		undoCmd, // Undo
		nclib.NewShell(vtycmd, "config", "save", "-H", name), // End
	)
}

func AddNISysctlConfigCmd(h NICommandsHandler, name string) {
	cmd := cliConfig().SysPath()

	h.OnceCmd(NI_UPDATE_SYSCTL,
		nclib.NewShell(cmd, "sysctl", "backup", "-H", name),   // Do
		nclib.NewShell(cmd, "sysctl", "rollback", "-H", name), // Undo
		nclib.NewShell(cmd, "sysctl", "load", "-H", name),     // End
	)
}

func AddNIVrfConfigCmd(h NICommandsHandler, name string) {
	cmd := cliConfig().SysPath()

	h.OnceCmd(NI_UPDATE_SYSVRF,
		nclib.NewShell(cmd, "vrf", "backup", "-H", name),   // Do
		nclib.NewShell(cmd, "vrf", "rollback", "-H", name), // Undo
		nclib.NewShell(cmd, "vrf", "load", "-H", name),     // End
	)

}

func AddNINetworkConfigCmd(h NICommandsHandler, name string) {
	cmd := cliConfig().SysPath()

	h.OnceCmd(NI_UPDATE_NETWORK,
		nclib.NewShell(cmd, "network", "backup", "-H", name),   // Do
		nclib.NewShell(cmd, "network", "rollback", "-H", name), // Undo
		nclib.NewShell(cmd, "network", "load", "-H", name),     // End
	)
}
