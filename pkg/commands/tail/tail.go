// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tail

import (
	"context"
	"encoding/json"
	"os"
	"regexp"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	color "github.com/zchee/color/v2"
	errors "golang.org/x/xerrors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	"github.com/zchee/kt/pkg/controller"
	"github.com/zchee/kt/pkg/internal/unsafes"
	"github.com/zchee/kt/pkg/io"
	"github.com/zchee/kt/pkg/manager"
	"github.com/zchee/kt/pkg/options"
)

const (
	tailShort = "tail Kubernetes logs for a container in a pod or specified resource."
)

type Tail struct {
	ctrl *controller.Controller
	mgr  *manager.Manager
}

func NewCmdTail(ctx context.Context, ioStreams io.Streams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tail pod-query [flags]",
		Short: tailShort,
	}

	opts := &options.Options{
		Container:      ".*",
		ContainerState: "running",
		Since:          48 * time.Hour,
		Concurrency:    1,
		UseColor:       "auto",
		Format:         "",
		Output:         "default",
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.Help, "help", "h", opts.Help, "Show help")

	// kubeconfig and context
	f.StringVar(&opts.KubeConfig, "kubeconfig", opts.KubeConfig, "Path to kubeconfig file to use")
	f.StringVar(&opts.KubeContext, "context", opts.KubeContext, "Kubernetes context to use. Default to current context configured in kubeconfig.")

	// global filters
	f.StringSliceVarP(&opts.Exclude, "exclude", "e", opts.Exclude, "Regex of log lines to exclude")
	f.StringSliceVarP(&opts.Include, "include", "i", opts.Include, "Regex of log lines to include")

	// pod filters
	f.StringVarP(&opts.Container, "container", "c", opts.Container, "Container name when multiple containers in pod")
	f.StringVar(&opts.ContainerState, "container-state", opts.ContainerState, "If present, tail containers with status in running, waiting or terminated. Default to running.")
	f.StringVarP(&opts.ExcludeContainer, "exclude-container", "E", opts.ExcludeContainer, "Exclude a Container name")
	f.StringSliceVarP(&opts.Namespaces, "namespaces", "n", opts.Namespaces, "Kubernetes namespace to use. Default to namespace configured in Kubernetes context. can set command separated multiple namespaces.")
	f.BoolVar(&opts.AllNamespaces, "all-namespaces", opts.AllNamespaces, "If present, tail across all namespaces. A specific namespace is ignored even if specified with --namespace.")
	f.StringVarP(&opts.Selector, "selector", "l", opts.Selector, "Selector (label query) to filter on. If present, default to \".*\" for the pod-query.")
	f.BoolVarP(&opts.Timestamps, "timestamps", "t", opts.Timestamps, "Print timestamps")
	f.DurationVarP(&opts.Since, "since", "s", opts.Since, "Return logs newer than a relative duration like 5s, 2m, or 3h.")
	f.IntVar(&opts.Concurrency, "concurrency", opts.Concurrency, "max concurrent reconciler.")

	// misc options
	f.Int64Var(&opts.Lines, "tail", opts.Lines, "The number of lines from the end of the logs to show. Defaults to -1, showing all logs.")
	f.StringVar(&opts.UseColor, "color", opts.UseColor, "Color output. Can be 'always', 'never', or 'auto'")
	f.StringVarP(&opts.Format, "format", "f", opts.Format, "Template to use for log lines, leave empty to use --output flag")
	f.StringVarP(&opts.Output, "output", "o", opts.Output, "Specify predefined template. Currently support: [default, raw, json]")

	cmd.PreRunE = func(*cobra.Command, []string) error {
		if opts.Help {
			return cmd.Usage()
		}
		return nil
	}

	cmd.RunE = func(_ *cobra.Command, args []string) error {
		podQuery := ".*"
		if len(args) == 1 {
			podQuery = args[0]
		}
		pQuery, err := regexp.Compile(podQuery)
		if err != nil {
			return errors.Errorf("failed to compile regular expression from query: %w", err)
		}
		opts.PodQuery = pQuery

		cQuery, err := regexp.Compile(opts.Container)
		if err != nil {
			return errors.Errorf("failed to compile regular expression for container query: %w", err)
		}
		opts.ContainerQuery = cQuery

		if opts.KubeConfig == "" {
			opts.KubeConfig = os.Getenv("KUBECONFIG")
			if opts.KubeConfig == "" {
				opts.KubeConfig = clientcmd.RecommendedHomeFile
			}
		}

		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{
				ExplicitPath: opts.KubeConfig,
			},
			&clientcmd.ConfigOverrides{
				CurrentContext: opts.KubeContext,
			},
		)

		config, err := clientConfig.ClientConfig()
		if err != nil {
			return errors.Errorf("unable create rest config: %w", err)
		}
		// TODO(zchee): inject OpenCensus RoundTripper
		// config.Transport = trace.Transport()

		var mgrOpts manager.Options
		switch {
		case opts.AllNamespaces:
			mgrOpts.Namespace = metav1.NamespaceAll
		case len(opts.Namespaces) >= 1:
			if len(opts.Namespaces) == 1 {
				mgrOpts.Namespace = opts.Namespaces[0]
			} else {
				mgrOpts.NewCache = cache.MultiNamespacedCacheBuilder(opts.Namespaces)
			}
		default:
			rawConfig, err := clientConfig.RawConfig()
			if err != nil {
				return errors.Errorf("unable get raw config: %w", err)
			}
			if currentNamespace := rawConfig.Contexts[rawConfig.CurrentContext].Namespace; currentNamespace != "" {
				mgrOpts.Namespace = currentNamespace
			}
		}

		t := &Tail{}
		t.mgr, err = manager.New(config, mgrOpts)
		if err != nil {
			return errors.Errorf("unable create manager: %w", err)
		}

		switch opts.UseColor {
		case "auto":
			// nothig to do
		case "always":
			color.NoColor = false
		case "never":
			color.NoColor = true
		default:
			return errors.New("color flag should be one of 'always', 'never', or 'auto'")
		}

		if format := opts.Format; format == "" {
			switch opts.Output {
			case "default":
				if color.NoColor {
					format = "{{.PodName}} {{.ContainerName}} {{.Message}}\n"
					if opts.AllNamespaces {
						format = "{{.Namespace}} " + format
					}
				} else {
					format = "{{color .PodColor .PodName}} {{color .ContainerColor .ContainerName}} {{.Message}}\n"
					if opts.AllNamespaces {
						format = "{{color .PodColor .Namespace}} " + format
					}

				}
			case "raw":
				format = "{{.Message}}"
			case "json":
				format = "{{json .}}\n"
			}
			opts.Format = format
		}

		tmplFuncs := map[string]interface{}{
			"json": func(in interface{}) (string, error) {
				b, err := json.Marshal(in)
				if err != nil {
					return "", err
				}
				return unsafes.String(b), nil
			},
			"color": func(c color.Color, text string) string {
				return c.SprintFunc()(text)
			},
		}
		opts.Template = template.Must(template.New("log").Funcs(tmplFuncs).Parse(opts.Format))

		t.ctrl, err = controller.New(ctx, ioStreams, t.mgr, opts)
		t.ctrl.Clientset, err = kubernetes.NewForConfig(t.mgr.GetConfig())
		if err != nil {
			return errors.Errorf("failed to new clientset: %w", err)
		}

		return t.RunTail(ctx)
	}

	return cmd
}

func (t *Tail) RunTail(ctx context.Context) error {
	return t.mgr.Start(ctx.Done())
}
