// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package signalcontext

import (
	"os"

	"golang.org/x/sys/unix"
)

var shutdownSignals = []os.Signal{os.Interrupt, unix.SIGTERM}
