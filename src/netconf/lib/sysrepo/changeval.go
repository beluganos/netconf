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

package srlib

import (
	"fmt"
	"netconf/lib/xml"
)

type SrChangeVal struct {
	Oper   SrChangeOper
	OldVal *SrVal
	NewVal *SrVal
}

type SrChangeVals interface {
	Recv() <-chan *SrChangeVal
}

type SrChangeValHandler interface {
	Put([]*ncxml.XPathNode, string) error
}

func (v *SrChangeVal) String() string {
	val := func() string {
		if v.OldVal == nil {
			return v.NewVal.String()
		} else {
			if v.NewVal == nil {
				return v.OldVal.String()
			} else {
				return fmt.Sprintf("%s -> %s", v.OldVal, v.NewVal)
			}
		}
	}()
	return fmt.Sprintf("%s: %s", v.Oper, val)
}

func (c *SrChangeVal) Dispatch(cre, mod, del SrChangeValHandler) error {
	switch c.Oper {
	case SR_OP_CREATED:
		xpath, data := c.NewVal.Xpath, c.NewVal.Data
		nodes := ParseXPath(xpath)
		if err := cre.Put(nodes[1:], data); err != nil {
			return err
		}

	case SR_OP_MODIFIED:
		xpath, data := c.NewVal.Xpath, c.NewVal.Data
		nodes := ParseXPath(xpath)
		if err := mod.Put(nodes[1:], data); err != nil {
			return nil
		}

	case SR_OP_DELETED:
		xpath, data := c.OldVal.Xpath, c.OldVal.Data
		nodes := ParseXPath(xpath)
		if err := del.Put(nodes[1:], data); err != nil {
			return err
		}

	default:

	}

	return nil
}
