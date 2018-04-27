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

package ncproplib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Property interface {
	Set(string, string) error
	Get(func(string) error) error
}

func ReadFile(path string, pr Property) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return Read(f, pr)
}

func Read(rd io.Reader, pr Property) error {
	br := bufio.NewReader(rd)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		k, v, err := ParseLine(line)
		if err != nil {
			return err
		}

		if err := pr.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

func ParseLine(line string) (string, string, error) {
	items := strings.SplitN(line, "=", 2)
	if len(items) != 2 {
		return "", "", fmt.Errorf("Invalid format. '%s'", line)
	}
	return strings.TrimSpace(items[0]), strings.TrimSpace(items[1]), nil
}

func WriteFile(path string, pr Property) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return Write(f, pr)
}

func Write(wr io.Writer, pr Property) error {
	bw := bufio.NewWriter(wr)
	defer bw.Flush()

	return pr.Get(func(line string) error {
		_, err := bw.WriteString(line + "\n")
		return err
	})
}
