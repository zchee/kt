// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

var (
	// version indicates which version of the binary is running.
	version = "dev"

	// GitCommit indicates which git hash the binary was built off of.
	gitCommit = ""
)

// Version is the current spinctl version.
func Version() string {
	if gitCommit != "" {
		version += "@" + gitCommit
	}

	return version
}
