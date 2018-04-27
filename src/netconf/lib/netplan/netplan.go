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

package ncnplib

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func ListConfigFiles(path string) []string {
	pattern := fmt.Sprintf("%s/*.yaml", path)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return []string{}
	}

	return matches
}

func ReadConfigDir(dir string) (*Config, error) {
	return ReadConfigFiles(ListConfigFiles(dir))
}

func ReadConfigFiles(paths []string) (*Config, error) {
	c := NewConfig()
	for _, path := range paths {
		cfg, err := ReadConfigFile(path)
		if err != nil {
			return nil, err
		}

		c.Merge(cfg)
	}

	return c, nil
}

func ReadConfigFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadConfig(f)
}

func ReadConfig(r io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func WriteConfigFile(path string, c *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return WriteConfig(f, c)
}

func WriteConfig(w io.Writer, c *Config) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

type BondMode string

const (
	BOND_MODE_BALANCE_RR    BondMode = "balance-rr"
	BOND_MODE_ACTIVE_BACKUP BondMode = "active-backup"
	BOND_MODE_BALANCE_XOR   BondMode = "balance-xor"
	BOND_MODE_BROADCAST     BondMode = "broadcast"
	BOND_MODE_802_3_AD      BondMode = "802.3ad"
	BOND_MODE_BALANCE_TLB   BondMode = "balance-tlb"
	BOND_MODE_BALANCE_ALB   BondMode = "balance-alb"
)

var bondModeNames = map[BondMode]interface{}{
	BOND_MODE_BALANCE_RR:    nil,
	BOND_MODE_ACTIVE_BACKUP: nil,
	BOND_MODE_BALANCE_XOR:   nil,
	BOND_MODE_BROADCAST:     nil,
	BOND_MODE_802_3_AD:      nil,
	BOND_MODE_BALANCE_TLB:   nil,
	BOND_MODE_BALANCE_ALB:   nil,
}

func ParseBondMode(s string) (BondMode, error) {
	v := BondMode(s)
	if _, ok := bondModeNames[v]; ok {
		return v, nil
	}
	return "", fmt.Errorf("Invalid BondMode. %s", s)
}

type BondTransmitHashPolicy string

const (
	BOND_TRANS_HASH_POLICY_L2    BondTransmitHashPolicy = "layer2"
	BOND_TRAMS_HASH_POLICY_L23   BondTransmitHashPolicy = "layer2+3"
	BOND_TRAMS_HASH_POLICY_L34   BondTransmitHashPolicy = "layer3+4"
	BOND_TRAMS_HASH_POLICY_ENC23 BondTransmitHashPolicy = "encap2+3"
	BOND_TRANS_HASH_POLICY_ENC34 BondTransmitHashPolicy = "encap3+4"
)

var bondTransmitHashPolicyNames = map[BondTransmitHashPolicy]interface{}{
	BOND_TRANS_HASH_POLICY_L2:    nil,
	BOND_TRAMS_HASH_POLICY_L23:   nil,
	BOND_TRAMS_HASH_POLICY_L34:   nil,
	BOND_TRAMS_HASH_POLICY_ENC23: nil,
	BOND_TRANS_HASH_POLICY_ENC34: nil,
}

func ParseBondTransmitHashPolicy(s string) (BondTransmitHashPolicy, error) {
	v := BondTransmitHashPolicy(s)
	if _, ok := bondTransmitHashPolicyNames[v]; ok {
		return v, nil
	}
	return "", fmt.Errorf("Invalid BondTransmitHashPolicy. %s", s)
}

type BondAdSelect string

const (
	BOND_AD_SELECT_STABLE   BondAdSelect = "stable"
	BOND_AD_SELECT_BANDWITH BondAdSelect = "bandwith"
	BOND_AD_SELECT_COUNT    BondAdSelect = "count"
)

var bondAdSelectNames = map[BondAdSelect]interface{}{
	BOND_AD_SELECT_STABLE:   nil,
	BOND_AD_SELECT_BANDWITH: nil,
	BOND_AD_SELECT_COUNT:    nil,
}

func ParseBondAdSelect(s string) (BondAdSelect, error) {
	v := BondAdSelect(s)
	if _, ok := bondAdSelectNames[v]; ok {
		return v, nil
	}
	return "", fmt.Errorf("Invalid BondAdSelect. %s", s)
}

type BondARPValidate string

const (
	BOND_ARP_VALIDATE_NONE   BondARPValidate = "none"
	BOND_ARP_VALIDATE_ACTIVE BondARPValidate = "active"
	BOND_ARP_VALIDATE_BACKUP BondARPValidate = "backup"
	BOND_ARP_VALIDATE_ALL    BondARPValidate = "all"
)

var bondARPValidateNames = map[BondARPValidate]interface{}{
	BOND_ARP_VALIDATE_NONE:   nil,
	BOND_ARP_VALIDATE_ACTIVE: nil,
	BOND_ARP_VALIDATE_BACKUP: nil,
	BOND_ARP_VALIDATE_ALL:    nil,
}

func ParseBondARPValidate(s string) (BondARPValidate, error) {
	v := BondARPValidate(s)
	if _, ok := bondARPValidateNames[v]; ok {
		return v, nil
	}
	return "", fmt.Errorf("Invalid BondARPValidate. %s", s)
}

type BondARPAllTargets string

const (
	BOND_ARP_ALLTARGETS_ANY BondARPAllTargets = "any"
	BOND_ARP_ALLTARGETS_ALL BondARPAllTargets = "all"
)

var bondARPAllTargetsNames = map[BondARPAllTargets]interface{}{
	BOND_ARP_ALLTARGETS_ANY: nil,
	BOND_ARP_ALLTARGETS_ALL: nil,
}

func ParseBondARPAllTargets(s string) (BondARPAllTargets, error) {
	v := BondARPAllTargets(s)
	if _, ok := bondARPAllTargetsNames[v]; ok {
		return v, nil
	}
	return "", fmt.Errorf("Invalid BondARPAllTargets. %s", s)
}

type BondFailOverMACPolicy string

const (
	BOND_FO_MAC_POLICY_NONE   BondFailOverMACPolicy = "none"
	BOND_FO_MAC_POLICY_ACTIVE BondFailOverMACPolicy = "active"
	BOND_FO_MAC_POLICY_FOLLOW BondFailOverMACPolicy = "follow"
)

var bondFailOverMACPolicyName = map[BondFailOverMACPolicy]interface{}{
	BOND_FO_MAC_POLICY_NONE:   nil,
	BOND_FO_MAC_POLICY_ACTIVE: nil,
	BOND_FO_MAC_POLICY_FOLLOW: nil,
}

func ParseBondFailOverMACPolicy(s string) (BondFailOverMACPolicy, error) {
	v := BondFailOverMACPolicy(s)
	if _, ok := bondFailOverMACPolicyName[v]; ok {
		return v, nil
	}
	return "", fmt.Errorf("Invalid BondFailOverMACPolicy. %s", s)
}

type BondPrimaryReselectPolicy string

const (
	BOND_PRI_RESELECT_POLICY_ALYAWS  BondPrimaryReselectPolicy = "always"
	BOND_PRI_RESELECT_POLICY_BETTER  BondPrimaryReselectPolicy = "better"
	BOND_PRI_RESELECT_POLICY_FAILURE BondPrimaryReselectPolicy = "failure"
)

var bondPrimaryReselectPolicy = map[BondPrimaryReselectPolicy]interface{}{
	BOND_PRI_RESELECT_POLICY_ALYAWS:  nil,
	BOND_PRI_RESELECT_POLICY_BETTER:  nil,
	BOND_PRI_RESELECT_POLICY_FAILURE: nil,
}

func ParseBondPrimaryReselectPolicy(s string) (BondPrimaryReselectPolicy, error) {
	v := BondPrimaryReselectPolicy(s)
	if _, ok := bondPrimaryReselectPolicy[v]; ok {
		return v, nil
	}
	return "", fmt.Errorf("Invalid BondPrimaryReselectPolicy. %s", s)
}
