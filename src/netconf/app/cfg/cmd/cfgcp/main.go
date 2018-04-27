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
	"flag"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Args struct {
	Move    bool
	Force   bool
	Verbose bool
	Args    []string
}

func (a *Args) Parse() {
	flag.BoolVar(&a.Move, "m", false, "move mode.")
	flag.BoolVar(&a.Force, "f", false, "no error if src not exist.")
	flag.BoolVar(&a.Verbose, "v", false, "show detail message.")
	flag.Parse()
	a.Args = flag.Args()

	if len(a.Args) != 2 {
		os.Exit(1)
	}
}

func (a *Args) Src() string {
	return a.Args[0]
}

func (a *Args) Dst() string {
	return a.Args[1]
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func touchFile(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func run(a *Args) error {
	if _, err := os.Stat(a.Src()); err != nil {
		if a.Force {
			log.Warnf("%s", err)
			if err := touchFile(a.Dst()); err != nil {
				log.Warnf("%s", err)
			}
			return nil
		}

		return err
	}

	if a.Move {
		return os.Rename(a.Src(), a.Dst())
	} else {
		return copyFile(a.Src(), a.Dst())
	}
}

func main() {
	a := Args{}
	a.Parse()

	if a.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	if err := run(&a); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}

	os.Exit(0)
}
