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

package nclib

type CommandAction string

const (
	CommandActionDo   = CommandAction("Do")
	CommandActionUndo = CommandAction("Undo")
	CommandActionEnd  = CommandAction("End")
)

type CommandMon func(CommandAction, Command, []byte)

func comandDefaultMon(CommandAction, Command, []byte) {}

type Command interface {
	DoCommand() ([]byte, error)
	EndCommand() ([]byte, error)
	UndoCommand() ([]byte, error)
	Line(CommandAction) string
}

func DoCommands(mon CommandMon, dryRun bool, cmds ...Command) error {
	if dryRun {
		for _, cmd := range cmds {
			mon(CommandActionDo, cmd, nil)
		}
		return nil
	}

	for index, cmd := range cmds {
		b, err := cmd.DoCommand()
		mon(CommandActionDo, cmd, b)
		if err != nil {
			UndoCommands(mon, dryRun, cmds[:index]...)
			return err
		}
	}

	return nil
}

func EndCommands(mon CommandMon, dryRun bool, cmds ...Command) error {
	if dryRun {
		for _, cmd := range cmds {
			mon(CommandActionEnd, cmd, nil)
		}
		return nil
	}

	for _, cmd := range cmds {
		b, err := cmd.EndCommand()
		mon(CommandActionEnd, cmd, b)
		if err != nil {
			return err
		}
	}
	return nil
}

func UndoCommands(mon CommandMon, dryRun bool, cmds ...Command) {
	if dryRun {
		for index := len(cmds) - 1; index >= 0; index-- {
			mon(CommandActionUndo, cmds[index], nil)
		}
	}

	for index := len(cmds) - 1; index >= 0; index-- {
		b, _ := cmds[index].UndoCommand()
		mon(CommandActionUndo, cmds[index], b)
	}
}

type Commands struct {
	cmds   []Command
	mon    CommandMon
	DryRun bool
}

func NewCommands(mon CommandMon) *Commands {
	return &Commands{
		cmds:   []Command{},
		mon:    mon,
		DryRun: false,
	}
}

func (c *Commands) Add(cmd Command) *Commands {
	c.cmds = append(c.cmds, cmd)
	return c
}

func (c *Commands) Set(cmd Command, pos int) *Commands {
	c.cmds[pos] = cmd
	return c
}

func (c *Commands) Clear() {
	c.cmds = []Command{}
}

func (c *Commands) Do() error {
	return DoCommands(c.mon, c.DryRun, c.cmds...)
}

func (c *Commands) Undo() {
	UndoCommands(c.mon, c.DryRun, c.cmds...)
}

func (c *Commands) End() error {
	return EndCommands(c.mon, c.DryRun, c.cmds...)
}

func (c *Commands) Size() int {
	return len(c.cmds)
}
