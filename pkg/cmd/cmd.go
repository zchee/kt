// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	iopkg "io"
	"os"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	color "github.com/zchee/color/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"github.com/zchee/kt/pkg/controller"
	regexp "github.com/zchee/kt/pkg/internal/lazyregexp"
	"github.com/zchee/kt/pkg/internal/unsafes"
	"github.com/zchee/kt/pkg/io"
	"github.com/zchee/kt/pkg/manager"
	"github.com/zchee/kt/pkg/options"
)

const (
	usageShort = "kt tails the Kubernetes logs for a container in a pod or specified resource."
	usageLong  = `
kt tails the Kubernetes logs for a container in a pod or specified resource.`

	versionTempl = `{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version: %s " .Version}}
`
)

// NewCommand creates the kt command with arguments.
func NewCommand() *cobra.Command {
	return NewKTCommand(os.Stdin, os.Stdout, os.Stderr)
}

type kt struct {
	ctrl *controller.Controller
	mgr  *manager.Manager

	ioStreams  io.Streams
	completion string
	opts       *options.Options
}

// NewKTCommand creates the `kt` command and its nested children.
func NewKTCommand(in iopkg.Reader, out, errOut iopkg.Writer) *cobra.Command {
	kt := &kt{}

	kt.ioStreams = io.Streams{In: in, Out: out, ErrOut: errOut}

	// set default options.Options.
	kt.opts = &options.Options{
		Container:      ".*",
		ContainerState: "running",
		Since:          48 * time.Hour,
		Concurrency:    10,
		UseColor:       "auto",
		Format:         "",
		Output:         "default",
	}

	cmd := &cobra.Command{
		Use:                "kt [pod-query]",
		Short:              usageShort,
		Long:               usageLong,
		Version:            Version(),
		PersistentPreRunE:  initProfiling(),  // Hook before and after Run initialize and write profiles to disk, respectively
		PersistentPostRunE: flushProfiling(), // Hook before and after Run initialize and write profiles to disk, respectively
	}

	addVersionFlag(cmd) // version flag is root only

	f := cmd.Flags()
	addProfilingFlags(f)
	f.AddGoFlagSet(flag.CommandLine)

	// kubeconfig and context
	f.StringVar(&kt.opts.KubeConfig, "kubeconfig", kt.opts.KubeConfig, "Path to kubeconfig file to use")
	f.StringVar(&kt.opts.KubeContext, "context", kt.opts.KubeContext, "Kubernetes context to use. Default to current context configured in kubeconfig.")

	// global filters
	f.StringSliceVarP(&kt.opts.Exclude, "exclude", "e", kt.opts.Exclude, "Regex of log lines to exclude")
	f.StringSliceVarP(&kt.opts.Include, "include", "i", kt.opts.Include, "Regex of log lines to include")

	// pod filters
	f.StringVarP(&kt.opts.Container, "container", "c", kt.opts.Container, "Container name when multiple containers in pod")
	f.StringVar(&kt.opts.ContainerState, "container-state", kt.opts.ContainerState, "If present, tail containers with status in running, waiting or terminated. Default to running.")
	f.StringVarP(&kt.opts.ExcludeContainer, "exclude-container", "E", kt.opts.ExcludeContainer, "Exclude a Container name")
	f.StringSliceVarP(&kt.opts.Namespaces, "namespaces", "n", kt.opts.Namespaces, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context. can set command separated multiple namespaces.")
	f.BoolVar(&kt.opts.AllNamespaces, "all-namespaces", kt.opts.AllNamespaces, "If present, tail across all namespaces. A specific namespace is ignored even if specified with --namespace.")
	f.StringVarP(&kt.opts.Selector, "selector", "l", kt.opts.Selector, "Selector (label query) to filter on. If present, default to \".*\" for the pod-query.")
	f.BoolVarP(&kt.opts.Timestamps, "timestamps", "t", kt.opts.Timestamps, "Print timestamps")
	f.DurationVarP(&kt.opts.Since, "since", "s", kt.opts.Since, "Return logs newer than a relative duration like 5s, 2m, or 3h.")
	f.IntVar(&kt.opts.Concurrency, "concurrency", kt.opts.Concurrency, "max concurrent reconciler.")

	// misc options
	f.BoolVarP(&kt.opts.Debug, "debug", "d", false, "debug mode.")
	f.Int64Var(&kt.opts.Lines, "tail", kt.opts.Lines, "The number of lines from the end of the logs to show. Defaults to -1, showing all logs.")
	f.StringVar(&kt.opts.UseColor, "color", kt.opts.UseColor, "Color output. Can be 'always', 'never', or 'auto'")
	f.StringVarP(&kt.opts.Format, "format", "f", kt.opts.Format, "Template to use for log lines, leave empty to use --output flag")
	f.StringVarP(&kt.opts.Output, "output", "o", kt.opts.Output, "Specify predefined template. Currently support: [default, raw, json]")

	// completions
	cmd.Flags().StringVar(&kt.completion, "completion", kt.completion, "Outputs kt command-line completion code for the specified shell. Can be 'bash' or 'zsh'")

	cmd.RunE = kt.Run(context.Background())

	return cmd
}

// Run runs the tail command.
func (kt *kt) Run(ctx context.Context) cobraRunEFunc {
	return func(cmd *cobra.Command, args []string) (err error) {
		if kt.completion != "" {
			return RunCompletion(kt.ioStreams.Out, kt.completion, cmd)
		}

		if kt.opts.KubeConfig == "" {
			kt.opts.KubeConfig = os.Getenv("KUBECONFIG")
			if kt.opts.KubeConfig == "" {
				kt.opts.KubeConfig = clientcmd.RecommendedHomeFile
			}
		}

		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{
				ExplicitPath: kt.opts.KubeConfig,
			},
			&clientcmd.ConfigOverrides{
				CurrentContext: kt.opts.KubeContext,
			},
		)

		config, err := clientConfig.ClientConfig()
		if err != nil {
			return fmt.Errorf("unable create rest config: %w", err)
		}
		// TODO(zchee): inject OpenCensus RoundTripper
		// config.Transport = trace.Transport()

		var mgrOpts manager.Options
		switch {
		case kt.opts.AllNamespaces:
			mgrOpts.Namespace = metav1.NamespaceAll
		case len(kt.opts.Namespaces) >= 1:
			if len(kt.opts.Namespaces) == 1 {
				mgrOpts.Namespace = kt.opts.Namespaces[0]
			} else {
				mgrOpts.NewCache = cache.MultiNamespacedCacheBuilder(kt.opts.Namespaces)
			}
		default:
			rawConfig, err := clientConfig.RawConfig()
			if err != nil {
				return fmt.Errorf("unable get raw config: %w", err)
			}
			if currentNamespace := rawConfig.Contexts[rawConfig.CurrentContext].Namespace; currentNamespace != "" {
				mgrOpts.Namespace = currentNamespace
			}
		}

		kt.mgr, err = manager.New(config, mgrOpts)
		if err != nil {
			return fmt.Errorf("unable create manager: %w", err)
		}

		switch kt.opts.UseColor {
		case "auto":
			// nothig to do
		case "always":
			color.NoColor = false
		case "never":
			color.NoColor = true
		default:
			return errors.New("color flag should be one of 'always', 'never', or 'auto'")
		}

		if kt.opts.Format == "" {
			var format string
			switch kt.opts.Output {
			case "default":
				if color.NoColor {
					format = "{{.PodName}} {{.ContainerName}} {{.Message}}\n"
					if kt.opts.AllNamespaces {
						format = "{{.Namespace}} " + format
					}
				} else {
					format = "{{color .PodColor .PodName}} {{color .ContainerColor .ContainerName}} {{.Message}}\n"
					if kt.opts.AllNamespaces {
						format = "{{color .PodColor .Namespace}} " + format
					}

				}
			case "raw":
				format = "{{.Message}}"
			case "json":
				format = "{{json .}}\n"
			}

			kt.opts.Format = format
		}

		tmplFuncs := map[string]interface{}{
			"json": func(v interface{}) (string, error) {
				b, err := json.Marshal(v)
				if err != nil {
					return "", err
				}
				return unsafes.String(b), nil
			},
			"color": func(c color.Color, text string) string {
				return c.SprintFunc()(text)
			},
		}
		kt.opts.Template = template.Must(template.New("log").Funcs(tmplFuncs).Parse(kt.opts.Format))

		query := new(options.Query)

		podQuery := ".*"
		if len(args) == 1 {
			podQuery = args[0]
		}
		query.PodQuery = regexp.New(podQuery)
		query.ContainerQuery = regexp.New(kt.opts.Container)

		query.ContainerState, err = options.NewContainerState(kt.opts.ContainerState)
		if err != nil {
			return err
		}

		kt.opts.Query = query

		kt.ctrl, err = controller.New(ctx, kt.ioStreams, kt.mgr, kt.opts)
		if err != nil {
			return fmt.Errorf("failed to create controller: %w", err)
		}
		defer kt.ctrl.Close()

		return kt.mgr.Start(ctx.Done())
	}
}
