// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package commands provides the kt commands.
package commands

import (
	"context"
	"flag"
	iopkg "io"
	"os"

	"github.com/spf13/cobra"

	"github.com/zchee/kt/pkg/commands/completion"
	"github.com/zchee/kt/pkg/commands/tail"
	"github.com/zchee/kt/pkg/io"
)

const (
	usageShort = "kt tails the Kubernetes logs for a container in a pod or specified resource."
	usageLong  = `
kt tails the Kubernetes logs for a container in a pod or specified resource.`

	versionTempl = `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version: %s " .Version}}
`
)

// NewCommand creates the `kt` command with arguments.
func NewCommand(ctx context.Context) *cobra.Command {
	return NewKTCommand(ctx, os.Stdin, os.Stdout, os.Stderr)
}

// NewKTCommand creates the `kt` command and its nested children.
func NewKTCommand(ctx context.Context, in iopkg.Reader, out, errOut iopkg.Writer) *cobra.Command {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:           "kt",
		Short:         usageShort,
		Long:          usageLong,
		Version:       Version(),
		SilenceErrors: true,
		// Hook before and after Run initialize and write profiles to disk, respectively
		PersistentPreRunE:  func(*cobra.Command, []string) error { return initProfiling() },
		PersistentPostRunE: func(*cobra.Command, []string) error { return flushProfiling() },
	}
	cmds.SetVersionTemplate(versionTempl)

	cmds.Flags().BoolP("version", "v", false, "Show "+cmds.Name()+" version.") // version flag is root only

	flags := cmds.PersistentFlags()
	addProfilingFlags(flags)
	flags.AddGoFlagSet(flag.CommandLine)

	ioStreams := io.Streams{In: in, Out: out, ErrOut: errOut}

	cmdTail := tail.NewCmdTail(ctx, ioStreams)
	cmds.AddCommand(cmdTail)
	cmds.AddCommand(completion.NewCmdCompletion(ctx, ioStreams.Out, ""))

	cmds.RunE = func(cmd *cobra.Command, args []string) error {
		return cmdTail.RunE(cmd, args)
	}

	return cmds
}
