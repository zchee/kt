// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tail

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	errors "golang.org/x/xerrors"

	cmdoptions "github.com/zchee/kt/pkg/command/options"
	"github.com/zchee/kt/pkg/controller"
	"github.com/zchee/kt/pkg/kube"
)

const (
	tailShort = "tail Kubernetes logs for a container in a pod or specified resource."
)

type Tail struct {
	ctrl        *controller.Controller
	kubeconfig  string
	concurrency int
}

func NewCmdTail(ctx context.Context, ioStreams cmdoptions.IOStreams) *cobra.Command {
	t := new(Tail)

	cmd := &cobra.Command{
		Use:   "tail [flags]",
		Short: tailShort,
	}

	f := cmd.Flags()
	f.StringVar(&t.kubeconfig, "kubeconfig", os.Getenv("KUBECONFIG"), "path to a kubeconfig.")
	f.IntVarP(&t.concurrency, "concurrency", "c", 1, "max concurrent reconciler.")

	cmd.PreRunE = func(*cobra.Command, []string) error {
		config, err := kube.RestConfig(t.kubeconfig)
		if err != nil {
			return errors.Errorf("unable create rest config: %w", err)
		}
		// TODO(zchee): inject OpenCensus RoundTripper
		// config.Transport = trace.Transport()

		t.ctrl, err = controller.NewController(ctx, config, t.concurrency)
		if err != nil {
			return errors.Errorf("unable create controller: %w", err)
		}

		return nil
	}

	cmd.RunE = func(*cobra.Command, []string) error {
		return t.RunTail(ctx, ioStreams)
	}

	return cmd
}

func (t *Tail) RunTail(ctx context.Context, ioStreams cmdoptions.IOStreams) error {
	return t.ctrl.Start(ctx.Done())
}
