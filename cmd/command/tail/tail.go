// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tail

import (
	"context"
	"os"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	errors "golang.org/x/xerrors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/zchee/kt/pkg/cmdoptions"
	"github.com/zchee/kt/pkg/controller"
	"github.com/zchee/kt/pkg/manager"
)

const (
	tailShort = "tail Kubernetes logs for a container in a pod or specified resource."
)

type Tail struct {
	ctrl *controller.Controller
	mgr  *manager.Manager
}

type Options struct {
	// kubeconfig and context
	kubeConfig  string
	kubeContext string

	// global filters
	exclude []string
	include []string

	// pod filters
	container        string
	containerState   string
	excludeContainer string
	namespace        string
	allNamespaces    bool
	selector         string
	timestamps       bool
	since            time.Duration
	concurrency      int

	// misc options
	lines      int64
	color      string
	tmplString string
	output     string
	help       bool
}

func NewCmdTail(ctx context.Context, ioStreams cmdoptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tail pod-query [flags]",
		Short: tailShort,
	}

	opts := &Options{
		container:      ".*",
		containerState: "running",
		since:          48 * time.Hour,
		concurrency:    1,
		lines:          -1,
		color:          "auto",
		tmplString:     "",
		output:         "default",
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.help, "help", "h", opts.help, "Show help")

	// kubeconfig and context
	f.StringVar(&opts.kubeConfig, "kubeconfig", opts.kubeConfig, "Path to kubeconfig file to use")
	f.StringVar(&opts.kubeContext, "context", opts.kubeContext, "Kubernetes context to use. Default to current context configured in kubeconfig.")

	// global filters
	f.StringSliceVarP(&opts.exclude, "exclude", "e", opts.exclude, "Regex of log lines to exclude")
	f.StringSliceVarP(&opts.include, "include", "i", opts.include, "Regex of log lines to include")

	// pod filters
	f.StringVarP(&opts.container, "container", "c", opts.container, "Container name when multiple containers in pod")
	f.StringVar(&opts.containerState, "container-state", opts.containerState, "If present, tail containers with status in running, waiting or terminated. Default to running.")
	f.StringVarP(&opts.excludeContainer, "exclude-container", "E", opts.excludeContainer, "Exclude a Container name")
	f.StringVarP(&opts.namespace, "namespace", "n", opts.namespace, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context")
	f.BoolVar(&opts.allNamespaces, "all-namespaces", opts.allNamespaces, "If present, tail across all namespaces. A specific namespace is ignored even if specified with --namespace.")
	f.StringVarP(&opts.selector, "selector", "l", opts.selector, "Selector (label query) to filter on. If present, default to \".*\" for the pod-query.")
	f.BoolVarP(&opts.timestamps, "timestamps", "t", opts.timestamps, "Print timestamps")
	f.DurationVarP(&opts.since, "since", "s", opts.since, "Return logs newer than a relative duration like 5s, 2m, or 3h.")
	f.IntVar(&opts.concurrency, "concurrency", opts.concurrency, "max concurrent reconciler.")

	// misc options
	f.Int64Var(&opts.lines, "tail", opts.lines, "The number of lines from the end of the logs to show. Defaults to -1, showing all logs.")
	f.StringVar(&opts.color, "color", opts.color, "Color output. Can be 'always', 'never', or 'auto'")
	f.StringVar(&opts.tmplString, "template", opts.tmplString, "Template to use for log lines, leave empty to use --output flag")
	f.StringVarP(&opts.output, "output", "o", opts.output, "Specify predefined template. Currently support: [default, raw, json]")

	cmd.PreRunE = func(*cobra.Command, []string) error {
		if opts.help {
			return cmd.Usage()
		}

		return nil
	}

	cmd.RunE = func(*cobra.Command, []string) error {
		t := &Tail{}

		if opts.kubeConfig == "" {
			opts.kubeConfig = os.Getenv("KUBECONFIG")
			if opts.kubeConfig == "" {
				opts.kubeConfig = clientcmd.RecommendedHomeFile
			}
		}

		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{
				ExplicitPath: opts.kubeConfig,
			},
			&clientcmd.ConfigOverrides{
				CurrentContext: opts.kubeContext,
			},
		)

		config, err := clientConfig.ClientConfig()
		if err != nil {
			return errors.Errorf("unable create rest config: %w", err)
		}
		// TODO(zchee): inject OpenCensus RoundTripper
		// config.Transport = trace.Transport()

		mgrOpts := manager.Options{
			Namespace: metav1.NamespaceAll,
		}
		switch {
		case opts.allNamespaces:
			// already set
		case opts.namespace != "":
			mgrOpts.Namespace = opts.namespace
		default:
			rawConfig, err := clientConfig.RawConfig()
			if err != nil {
				return errors.Errorf("unable get raw config: %w", err)
			}
			if currentNamespace := rawConfig.Contexts[rawConfig.CurrentContext].Namespace; currentNamespace != "" {
				mgrOpts.Namespace = currentNamespace
			}
		}

		if opts.tmplString == "" {
			if opts.output == "raw" {
				opts.tmplString = `{{.Message}}`
			} else {
				opts.tmplString = `{{.Pod.Name}} {{.Container.Name}} {{.Message}}`
				if opts.allNamespaces {
					opts.tmplString = `{{.Pod.Namespace}}/` + opts.tmplString
				}
			}
			if opts.timestamps {
				opts.tmplString = `{{.Timestamp}} ` + opts.tmplString
			}
		}
		opts.tmplString += "\n"

		tmpl, err := template.New("line").Parse(opts.tmplString)
		if err != nil {
			return errors.Errorf("invalid template: %w", err)
		}

		t.mgr, err = manager.NewManager(config, mgrOpts)
		if err != nil {
			return errors.Errorf("unable create manager: %w", err)
		}
		ctrlOpts := []controller.Options{
			controller.WithIOStearms(ioStreams),
			controller.WithTemplate(tmpl),
		}
		if opts.concurrency > 1 {
			ctrlOpts = append(ctrlOpts, controller.WithConcurrency(opts.concurrency))
		}
		t.ctrl, err = controller.NewController(ctx, t.mgr, ctrlOpts...)
		if err != nil {
			return errors.Errorf("unable create controller: %w", err)
		}

		return t.RunTail(ctx, ioStreams)
	}

	return cmd
}

func (t *Tail) RunTail(ctx context.Context, ioStreams cmdoptions.IOStreams) error {
	return t.mgr.Start(ctx.Done())
}
