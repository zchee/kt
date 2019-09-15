// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	// initialize all known client auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/zchee/kt/pkg/commands"
	"github.com/zchee/kt/pkg/signalcontext"
)

func main() {
	ctx, cancel := context.WithCancel(signalcontext.NewContext())
	defer cancel()

	if err := commands.NewCommand(ctx).Execute(); err != nil {
		os.Exit(1)
	}
}
