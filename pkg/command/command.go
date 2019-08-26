// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package command provides the kt commands.
package command

import (
	"context"
	"flag"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/zchee/kt/pkg/command/completion"
	cmdoptions "github.com/zchee/kt/pkg/command/options"
	"github.com/zchee/kt/pkg/command/tail"
)

const (
	usageShort = "kt tails the Kubernetes logs for a container in a pod or specified resource."
	usageLong  = `
kt tails the Kubernetes logs for a container in a pod or specified resource.`
)

// NewCommand creates the `kt` command with arguments.
func NewCommand(ctx context.Context) *cobra.Command {
	return NewKTCommand(ctx, os.Stdin, os.Stdout, os.Stderr)
}

// NewKTCommand creates the `kt` command and its nested children.
func NewKTCommand(ctx context.Context, in io.Reader, out, err io.Writer) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "kt",
		Short: usageShort,
		Long:  usageLong,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		// Hook before and after Run initialize and write profiles to disk,
		// respectively.
		PersistentPreRunE: func(*cobra.Command, []string) error {
			return initProfiling()
		},
		PersistentPostRunE: func(*cobra.Command, []string) error {
			return flushProfiling()
		},
	}

	flags := cmds.PersistentFlags()
	addProfilingFlags(flags)
	cmds.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	ioStreams := cmdoptions.IOStreams{In: in, Out: out, ErrOut: err}

	cmds.AddCommand(tail.NewCmdTail(ctx, ioStreams))
	cmds.AddCommand(completion.NewCmdCompletion(ctx, ioStreams.Out, ""))

	return cmds
}
