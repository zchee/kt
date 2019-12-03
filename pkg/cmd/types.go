// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "github.com/spf13/cobra"

type cobraRunEFunc func(cmd *cobra.Command, args []string) error
