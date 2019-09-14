// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"bufio"
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	xxhash "github.com/cespare/xxhash/v2"
	"github.com/go-logr/logr"
	color "github.com/zchee/color/v2"
	"go.uber.org/zap"
	errors "golang.org/x/xerrors"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	ctrlbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/zchee/kt/internal/unsafes"
	"github.com/zchee/kt/pkg/io"
	"github.com/zchee/kt/pkg/options"
)

type Controller struct {
	ctrlclient.Client
	Manager   ctrlmanager.Manager
	Log       logr.Logger
	Clientset kubernetes.Interface

	ctx       context.Context
	ioStreams io.Streams
	opts      *options.Options
	ioMu      sync.Mutex
}

var _ ctrlreconcile.Reconciler = (*Controller)(nil)

// New returns a new Controller registered with the Manager.
func New(ctx context.Context, ioStreams io.Streams, mgr ctrlmanager.Manager, opts *options.Options) (*Controller, error) {
	lvl := zap.NewAtomicLevelAt(zap.DebugLevel)
	logger := ctrlzap.New(func(o *ctrlzap.Options) {
		o.Level = &lvl
		o.Development = true
		o.DestWritter = ioStreams.ErrOut
	})
	ctrllog.SetLogger(logger)

	c := &Controller{
		Client:    mgr.GetClient(),
		Manager:   mgr,
		Log:       logger.WithName("controller"),
		ctx:       ctx,
		ioStreams: ioStreams,
		opts:      opts,
	}

	if err := c.SetupWithManager(mgr); err != nil {
		c.Log.Error(err, "failed to create controller")
		return nil, err
	}

	return c, nil
}

var colorList = [][2]*color.Color{
	{color.New(color.FgHiCyan), color.New(color.FgCyan)},
	{color.New(color.FgHiGreen), color.New(color.FgGreen)},
	{color.New(color.FgHiMagenta), color.New(color.FgMagenta)},
	{color.New(color.FgHiYellow), color.New(color.FgYellow)},
	{color.New(color.FgHiBlue), color.New(color.FgBlue)},
	{color.New(color.FgHiRed), color.New(color.FgRed)},
}

func findColors(podName string) (podColor, containerColor *color.Color) {
	digest := xxhash.New()
	digest.Write(unsafes.Slice(podName))
	idx := digest.Sum64() % uint64(len(colorList))

	colors := colorList[idx]
	return colors[0], colors[1]
}

type LogEvent struct {
	// Message is the log message itself
	Message string `json:"message"`

	// PodName of the pod
	PodName string `json:"podName"`

	// ContainerName of the container
	ContainerName string `json:"containerName"`

	// Namespace of the pod
	Namespace string `json:"namespace"`

	// Timestamp of the pod
	Timestamp *time.Time `json:"timestamp"`

	PodColor       *color.Color `json:"-"`
	ContainerColor *color.Color `json:"-"`
}

func (c *Controller) Reconcile(req ctrlreconcile.Request) (result ctrlreconcile.Result, err error) {
	log := c.Log.WithName("Reconcile")

	var pod corev1.Pod
	if err := c.Get(c.ctx, req.NamespacedName, &pod); err != nil {
		if ctrlclient.IgnoreNotFound(err) != nil {
			log.Error(err, "failed to get pod")
			return result, err
		}
		return result, nil
	}

	logOpts := corev1.PodLogOptions{
		Follow:     true,
		Timestamps: c.opts.Timestamps,
	}
	if c.opts.Lines > 0 {
		logOpts.TailLines = &c.opts.Lines
	}
	if c.opts.Timestamps {
		sec := int64(c.opts.Since.Seconds())
		logOpts.SinceSeconds = &sec
	}

	podColor, containerColor := findColors(pod.Name)

	boff := backoff.NewExponentialBackOff()
	for i := range pod.Spec.Containers {
		container := pod.Spec.Containers[i]
		logOpts.Container = container.Name

		stream, err := c.Clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &logOpts).Stream()
		if err != nil {
			if errStatus, ok := err.(apierrors.APIStatus); ok {
				switch errStatus.Status().Code {
				case http.StatusBadRequest:
					time.Sleep(boff.GetElapsedTime())
					continue
				case http.StatusNotFound:
					return result, err
				}
			}
			return result, err
		}

		go func(container corev1.Container, logOpts corev1.PodLogOptions) {
			r := bufio.NewReader(stream)

			for {
				l, err := r.ReadBytes('\n')
				if err != nil {
					if errors.Is(err, io.EOF) {
						stream.Close()
						break
					}
					c.Log.Error(err, "failed to ReadBytes")
					return
				}
				line := trimSpace(l)

				parts := strings.SplitN(line, " ", 2)
				if len(parts) < 2 {
					c.Log.Info("failed to split line", "line", line)
					continue
				}

				timeString, message := parts[0], parts[1]
				event := &LogEvent{
					Message:        message,
					PodName:        pod.Name,
					ContainerName:  container.Name,
					PodColor:       podColor,
					ContainerColor: containerColor,
				}

				if c.opts.Timestamps {
					timestamp, err := time.Parse(time.RFC3339Nano, timeString)
					if err != nil {
						c.Log.Error(err, "failed to parse timestamp", "timeString", timeString)
						return
					}
					event.Timestamp = &timestamp
				}

				// TODO(zchee): use goroutine
				c.ioMu.Lock()
				err = c.opts.Template.Execute(c.ioStreams.Out, event)
				c.ioMu.Unlock()
				if err != nil {
					c.Log.Error(err, "failed to tmpl.Execute", "event", event)
					return
				}
			}
		}(container, logOpts)

		go func() {
			<-c.ctx.Done()
			stream.Close()
		}()
	}

	log.Info("end of Reconcile")
	return result, nil
}

func (c *Controller) SetupWithManager(mgr ctrlmanager.Manager) error {
	ctrlOpts := ctrlcontroller.Options{
		Reconciler:              c,
		MaxConcurrentReconciles: c.opts.Concurrency,
	}

	return ctrlbuilder.ControllerManagedBy(mgr).For(&corev1.Pod{}).WithOptions(ctrlOpts).Complete(c)
}

func trimSpace(buf []byte) string {
	line := unsafes.String(buf)
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[0 : len(line)-1]
	}
	for len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[0 : len(line)-1]
	}

	return line
}
