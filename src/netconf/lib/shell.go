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

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Shell struct {
	cmd  string
	args []string
	In   io.Reader
}

func NewShell(cmd string, args ...string) *Shell {
	return NewShellIn(cmd, nil, args...)
}

func NewShellIn(cmd string, in io.Reader, args ...string) *Shell {
	return &Shell{
		cmd:  cmd,
		args: args,
		In:   in,
	}
}

func (s *Shell) String() string {
	if s == nil {
		return "-"
	}
	return fmt.Sprintf("%s %s", s.cmd, strings.Join(s.args, " "))
}

func (s *Shell) Exec() ([]byte, error) {
	if s == nil {
		return []byte{}, nil
	}
	cmd := exec.Command(s.cmd, s.args...)
	cmd.Stdin = s.In
	return cmd.CombinedOutput()
}

type ShellCommand struct {
	cmds map[CommandAction]*Shell
}

func NewShellCommand(doCmd, undoCmd, endCmd *Shell) *ShellCommand {
	return &ShellCommand{
		cmds: map[CommandAction]*Shell{
			CommandActionDo:   doCmd,
			CommandActionUndo: undoCmd,
			CommandActionEnd:  endCmd,
		},
	}
}

func (s *ShellCommand) String() string {
	return fmt.Sprintf("%s", s.cmds)
}

func (s *ShellCommand) Line(action CommandAction) string {
	return fmt.Sprintf("%s", s.cmds[action])
}

func (s *ShellCommand) DoCommand() ([]byte, error) {
	return s.cmds[CommandActionDo].Exec()
}

func (s *ShellCommand) EndCommand() ([]byte, error) {
	return s.cmds[CommandActionEnd].Exec()
}

func (s *ShellCommand) UndoCommand() ([]byte, error) {
	return s.cmds[CommandActionUndo].Exec()
}
