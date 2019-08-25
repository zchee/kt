// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	// Initialize all known client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/zchee/kt/pkg/command"
)

func main() {
	if err := command.NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
