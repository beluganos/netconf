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

package openconfig

import (
	"net"
	srlib "netconf/lib/sysrepo"
	"testing"
)

func makeStaticRoutes(datas [][2]string) StaticRoutes {
	sr := NewStaticRoutes()
	for _, data := range datas {
		xpath, value := data[0], data[1]
		nodes := srlib.ParseXPath(xpath)
		if err := sr.Put(nodes[1:], value); err != nil {
			panic(err)
		}
	}
	return sr
}

func TestStaticRoutes_Put_empty(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes", ""},
	})

	t.Log(routes)

	if v := len(routes); v != 0 {
		t.Errorf("StaticRoutes.Put unmatch. len=%d", v)
	}
}

func TestStaticRoutes_Put(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']", ""},
	})

	t.Log(routes)

	if v := len(routes); v != 1 {
		t.Errorf("StaticRoutes.Put unmatch. len=%d", v)
	}

	route, ok := routes[*NewStaticRouteKey("192.168.122.0", 24)]
	if !ok {
		t.Errorf("StaticRoutes.Put not found. %t.", ok)
	}
	if v := route.IP; v != "192.168.122.0" {
		t.Errorf("StaticRoutes.Put unmatch. ip=%s", v)
	}
	if v := route.PrefixLen; v != 24 {
		t.Errorf("StaticRoutes.Put unmatch. plen=%d", v)
	}
}

func TestStaticRoutes_Put_multi(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']", ""},
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='25']", ""},
		{"/static-routes/static[ip='192.168.123.0'][prefix-length='24']", ""},
		{"/static-routes/static[ip='192.168.123.0'][prefix-length='25']", ""},
	})

	t.Log(routes)

	if v := len(routes); v != 4 {
		t.Errorf("StaticRoutes.Put unmatch. len=%d", v)
	}

	if _, ok := routes[*NewStaticRouteKey("192.168.122.0", 24)]; !ok {
		t.Errorf("StaticRoutes.Put unmatch.")
	}
	if _, ok := routes[*NewStaticRouteKey("192.168.122.0", 25)]; !ok {
		t.Errorf("StaticRoutes.Put unmatch.")
	}
	if _, ok := routes[*NewStaticRouteKey("192.168.123.0", 24)]; !ok {
		t.Errorf("StaticRoutes.Put unmatch.")
	}
	if _, ok := routes[*NewStaticRouteKey("192.168.123.0", 25)]; !ok {
		t.Errorf("StaticRoutes.Put unmatch.")
	}
}

func TestStaticRoutes_Put_dup(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes", ""},
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/name", ""},
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']", ""},
	})

	t.Log(routes)

	if v := len(routes); v != 1 {
		t.Errorf("StaticRoutes.Put unmatch. len=%d", v)
	}

	if _, ok := routes[*NewStaticRouteKey("192.168.122.0", 24)]; !ok {
		t.Errorf("StaticRoutes.Put unmatch.")
	}
}

func TestStaticRoutes_Put_err(t *testing.T) {
	xpaths := []string{
		"/static-routes/static[ip='192.168.122.0']",
		"/static-routes/static[ip='192.168.122.0'][prefix-leng='24']",
		"/static-routes/static[prefix-length='24']",
		"/static-routes/static[ipp='192.168.122.0'][prefix-length='24']",
	}

	for _, xpath := range xpaths {
		routes := NewStaticRoutes()
		nodes := srlib.ParseXPath(xpath)
		if err := routes.Put(nodes[1:], ""); err == nil {
			t.Errorf("StaticRoutes.Put must be error. %s", xpath)
		}
	}
}

func TestNewStaticRouteKey(t *testing.T) {
	ip := "192.168.122.1"
	plen := uint8(24)

	key := NewStaticRouteKey(ip, plen)

	if v := key.String(); v != "192.168.122.1/24" {
		t.Errorf("NewStaticRouteKey unmatch. %s", key)
	}
}

func TestStaticRoute_ip(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/ip", "192.168.122.0"},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]

	if v := route.Compare(STATICROUTE_IP_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := route.IP; v != "192.168.122.0" {
		t.Errorf("StaticRoutes.Put not found. ip=%s", v)
	}
	if v := route.PrefixLen; v != 24 {
		t.Errorf("StaticRoutes.Put not found. plen=%d", v)
	}
}

func TestStaticRoute_prefixlen(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/prefix-length", "24"},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]

	if v := route.Compare(STATICROUTE_PREFIXLEN_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := route.IP; v != "192.168.122.0" {
		t.Errorf("StaticRoutes.Put not found. ip=%s", v)
	}
	if v := route.PrefixLen; v != 24 {
		t.Errorf("StaticRoutes.Put not found. plen=%d", v)
	}
}

func TestStaticRoute_config(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/config", ""},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]

	if v := route.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := route.Config.Compare(); !v {
		t.Errorf("StaticRoutes.Put unmatch. config.compare=%t", v)
	}
}

func TestStaticRoute_config_ip(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/config/ip", "192.168.122.0"},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]

	if v := route.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := route.Config.Compare(STATICROUTE_IP_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := route.Config.Ip; v.String() != "192.168.122.0" {
		t.Errorf("StaticRoutes.Put unmatch. config.ip=%s", v)
	}
}

func TestStaticRoute_config_plen(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/config/prefix-length", "24"},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]

	if v := route.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := route.Config.Compare(STATICROUTE_PREFIXLEN_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := route.Config.PrefixLen; v != 24 {
		t.Errorf("StaticRoutes.Put unmatch. config.plen=%d", v)
	}
}

func TestStaticRoute_config_error(t *testing.T) {
	xpaths := []string{
		"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/config/ip",
		"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/config/prefix-length",
	}

	for _, xpath := range xpaths {
		nodes := srlib.ParseXPath(xpath)
		routes := NewStaticRoutes()
		if err := routes.Put(nodes[1:], "ERROR"); err == nil {
			t.Errorf("StaticRoutes.Put must be error. %s %s", err, xpath)
		}
	}
}

func TestStaticRoute_nexthops(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops", ""},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]

	if v := route.Compare(STATICROUTE_NEXTHOPS_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. compare=%t", v)
	}
	if v := len(route.Nexthops); v != 0 {
		t.Errorf("StaticRoutes.Put unmatch. #nexthops=%d", v)
	}
}

func TestStaticRoute_nexthops_error(t *testing.T) {
	xpaths := []string{
		"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop",
	}

	for _, xpath := range xpaths {
		nodes := srlib.ParseXPath(xpath)
		routes := NewStaticRoutes()
		if err := routes.Put(nodes[1:], ""); err == nil {
			t.Errorf("StaticRoutes.Put must be error. %s, %s", err, xpath)
		}

	}
}

func TestStaticRoute_nexthop(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']", ""},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]

	if v := len(route.Nexthops); v != 1 {
		t.Errorf("StaticRoutes.Put unmatch. #nexthops=%d", v)
	}

	nh, ok := route.Nexthops["TEST1"]

	if !ok {
		t.Errorf("StaticRoutes.Put unmatch. nexthop[TEST1]=%t", ok)
	}
	if v := nh.Compare(); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
}

func TestStaticRoute_nexthop_multi(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']", ""},
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST2']", ""},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]
	if v := len(route.Nexthops); v != 2 {
		t.Errorf("StaticRoutes.Put unmatch. #nexthops=%d", v)
	}

	nh1, ok := route.Nexthops["TEST1"]

	if !ok {
		t.Errorf("StaticRoutes.Put unmatch. nexthop[TEST1]=%t", ok)
	}
	if v := nh1.Compare(); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
	if v := nh1.Index; v != "TEST1" {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.Index=%s", v)
	}

	nh2, ok := route.Nexthops["TEST2"]

	if !ok {
		t.Errorf("StaticRoutes.Put unmatch. nexthop[TEST2]=%t", ok)
	}
	if v := nh2.Compare(); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
	if v := nh2.Index; v != "TEST2" {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.Index=%s", v)
	}
}

func TestStaticRoute_nexthop_index(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/index", "TEST1"},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]
	nh := route.Nexthops["TEST1"]

	if v := nh.Compare(OC_INDEX_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
	if v := nh.Index; v != "TEST1" {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.Index=%s", v)
	}
}

func TestStaticRoute_nexthop_config(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/config", ""},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]
	nh := route.Nexthops["TEST1"]

	if v := nh.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
	if v := nh.Config.Compare(); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.config.compare=%t", v)
	}
}

func TestStaticRoute_nexthop_config_index(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/config/index", "TEST1"},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]
	nh := route.Nexthops["TEST1"]

	if v := nh.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
	if v := nh.Config.Compare(OC_INDEX_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.config.compare=%t", v)
	}
	if v := nh.Config.Index; v != "TEST1" {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.config.index=%s", v)
	}
}

func TestStaticRoute_nexthop_config_nexthop(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/config/next-hop", "LOCAL_LINK"},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]
	nh := route.Nexthops["TEST1"]

	if v := nh.Compare(OC_CONFIG_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
	if v := nh.Config.Compare(STATICROUTE_NEXTHOP_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.config.compare=%t", v)
	}
	if v := nh.Config.Nexthop; v != "LOCAL_LINK" {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.config.nexthop=%s", v)
	}
}

func TestStaticRoute_nexthop_ifref(t *testing.T) {
	routes := makeStaticRoutes([][2]string{
		{"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/interface-ref", ""},
	})

	t.Log(routes)

	route := routes[*NewStaticRouteKey("192.168.122.0", 24)]
	nh := route.Nexthops["TEST1"]

	if v := nh.Compare(INTERFACE_REF_KEY); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.compare=%t", v)
	}
	if v := nh.IfaceRef.Compare(); !v {
		t.Errorf("StaticRoutes.Put unmatch. nexthop.ifref.compare=%t", v)
	}
}

func TestStaticRoute_nexthop_error(t *testing.T) {
	xpaths := []string{
		"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/config/index",
		"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/config/next-hop",
		"/static-routes/static[ip='192.168.122.0'][prefix-length='24']/next-hops/next-hop[index='TEST1']/interface-ref/config/subinterface",
	}

	for _, xpath := range xpaths {
		routes := NewStaticRoutes()
		nodes := srlib.ParseXPath(xpath)
		if err := routes.Put(nodes[1:], ""); err == nil {
			t.Errorf("StaticRoutes.Put must be error. %s %s", err, xpath)
		}
	}
}

func TestStaticRoute_nexthop_config_getnexthop(t *testing.T) {
	c := NewStaticRouteNexthopConfig()

	var ip net.IP
	var nh LocalDefinedNexthop
	var err error

	c.Nexthop = "DROP"
	ip, nh, err = c.GetNexthop()
	if err != nil {
		t.Errorf("StaticRouteNexthopConfig.GetNexthop error. %s", err)
	}
	if ip != nil {
		t.Errorf("StaticRouteNexthopConfig.GetNexthop unmatch. ip=%s", ip)
	}
	if nh != LOCAL_DEFINED_NEXT_HOP_DROP {
		t.Errorf("StaticRouteNexthopConfig.GetNexthop unmatch. nexhop=%s", nh)
	}

	c.Nexthop = "192.168.1.1"
	ip, nh, err = c.GetNexthop()
	if err != nil {
		t.Errorf("StaticRouteNexthopConfig.GetNexthop error. %s", err)
	}
	if ip.String() != "192.168.1.1" {
		t.Errorf("StaticRouteNexthopConfig.GetNexthop unmatch. ip=%s", ip)
	}
	if nh != LOCAL_DEFINED_NEXT_HOP {
		t.Errorf("StaticRouteNexthopConfig.GetNexthop unmatch. nexhop=%s", nh)
	}
}

func TestStaticRoute_nexthop_config_getnexthop_error(t *testing.T) {
	c := NewStaticRouteNexthopConfig()

	if _, _, err := c.GetNexthop(); err == nil {
		t.Errorf("StaticRouteNexthopConfig.GetNexthop must be error. %s", err)
	}
}
