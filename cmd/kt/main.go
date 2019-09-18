// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	// initialize all known client auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/zchee/kt/pkg/commands"
)

func main() {
	if err := commands.NewCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
