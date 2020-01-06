// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kttesting

import (
	"os"
	"strings"
)

// InTest reports whether the current state is testing.
func InTest() bool {
	return len(os.Args) > 0 && strings.HasSuffix(strings.TrimSuffix(os.Args[0], ".exe"), ".test")
}
