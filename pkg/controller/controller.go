// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"bufio"
	"context"
	iopkg "io"
	"net/http"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/dgraph-io/ristretto"
	"github.com/go-logr/logr"
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
	ctrlpredicate "sigs.k8s.io/controller-runtime/pkg/predicate"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/zchee/kt/pkg/internal/unsafes"
	"github.com/zchee/kt/pkg/io"
	"github.com/zchee/kt/pkg/options"
)

// Controller implements a reconcile.Reconciler.
type Controller struct {
	client     ctrlclient.Client
	controller ctrlcontroller.Controller
	clientset  kubernetes.Interface
	mgr        ctrlmanager.Manager
	predicate  ctrlpredicate.Predicate
	log        logr.Logger

	ctx       context.Context
	ioStreams io.Streams
	ioMu      sync.Mutex
	opts      *options.Options
}

var _ ctrlreconcile.Reconciler = (*Controller)(nil)

// New returns the new Controller registered with the manager.Manager.
func New(ctx context.Context, ioStreams io.Streams, mgr ctrlmanager.Manager, opts *options.Options) (c *Controller, err error) {
	lvl := zap.NewAtomicLevelAt(zap.DebugLevel)
	log := ctrlzap.New(func(o *ctrlzap.Options) {
		o.Level = &lvl
		o.Development = true
		o.DestWritter = ioStreams.ErrOut
	}).WithName("controller")

	ctrllog.SetLogger(log)

	state, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	if err != nil {
		return nil, errors.Errorf("failed to create state cache: %w", err)
	}
	predicateFilter := &PredicateEventFilter{
		ioStreams:    ioStreams,
		state:        state,
		log:          log.WithName("predicate"),
		isNamespaced: (opts.AllNamespaces || len(opts.Namespaces) > 0),
		query:        opts.Query,
	}

	c = &Controller{
		client:    mgr.GetClient(),
		mgr:       mgr,
		predicate: predicateFilter,
		log:       log,
		ctx:       ctx,
		ioStreams: ioStreams,
		opts:      opts,
	}

	c.clientset, err = kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, errors.Errorf("failed to new clientset: %w", err)
	}

	if err := c.SetupWithManager(mgr); err != nil {
		c.log.Error(err, "failed to setup controller with manager", "Controller", c)
		return nil, err
	}

	return c, nil
}

const (
	lineDelim = '\n'
)

// Reconcile implements a ctrlreconcile.Reconciler.
func (c *Controller) Reconcile(req ctrlreconcile.Request) (result ctrlreconcile.Result, err error) {
	log := c.log.WithName("Reconcile").WithValues("req.Namespace", req.Namespace, "req.Name", req.Name)

	var pod corev1.Pod
	if err := c.client.Get(c.ctx, req.NamespacedName, &pod); err != nil {
		if ctrlclient.IgnoreNotFound(err) != nil {
			log.Error(err, "failed to get pod")
			return result, err
		}
		return result, nil
	}

	podColor, containerColor := findColors(pod.Name)

	logOpts := &corev1.PodLogOptions{
		Follow:     true,
		Timestamps: c.opts.Timestamps,
	}
	if c.opts.Lines > 0 {
		logOpts.TailLines = &c.opts.Lines
	}
	if c.opts.Since > 0 {
		sec := int64(c.opts.Since.Seconds())
		logOpts.SinceSeconds = &sec
	}

	boff := backoff.NewExponentialBackOff()
	for i := range pod.Spec.Containers {
		container := pod.Spec.Containers[i]
		podLogOpts := new(corev1.PodLogOptions)
		*podLogOpts = *logOpts // shallow copy
		podLogOpts.Container = container.Name

		stream, err := c.clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, podLogOpts).Context(c.ctx).Stream()
		if err != nil {
			if errStatus, ok := err.(apierrors.APIStatus); ok {
				switch errStatus.Status().Code {
				case http.StatusBadRequest:
					time.Sleep(boff.GetElapsedTime())
					continue
				case http.StatusNotFound:
					return result, nil
				}
			}
			return result, err
		}

		go func(containerName string, stream iopkg.ReadCloser) {
			defer stream.Close()

			r := bufio.NewReader(stream)
			for {
				l, err := r.ReadBytes(lineDelim)
				if err != nil {
					if errors.Is(err, iopkg.EOF) {
						stream.Close()
						return
					}
					c.log.Error(err, "failed to ReadBytes")
					return
				}
				line := trimSpace(unsafes.String(l))

				event := &LogEvent{
					Message:        line,
					PodName:        pod.Name,
					ContainerName:  containerName,
					PodColor:       podColor,
					ContainerColor: containerColor,
				}
				if c.opts.AllNamespaces || len(c.opts.Namespaces) > 0 {
					event.Namespace = pod.Namespace
				}

				// TODO(zchee): use goroutine
				c.ioMu.Lock()
				err = c.opts.Template.Execute(c.ioStreams.Out, event)
				c.ioMu.Unlock()
				if err != nil {
					c.log.Error(err, "failed to tmpl.Execute", "event", event)
					return
				}
			}
		}(container.Name, stream)
	}

	return result, nil
}

// SetupWithManager setups the Controller with manager.Manager.
func (c *Controller) SetupWithManager(mgr ctrlmanager.Manager) (err error) {
	ctrlOpts := ctrlcontroller.Options{
		Reconciler:              c,
		MaxConcurrentReconciles: c.opts.Concurrency,
	}

	c.controller, err = ctrlbuilder.ControllerManagedBy(mgr).For(&corev1.Pod{}).WithEventFilter(c.predicate).WithOptions(ctrlOpts).Build(c)
	if err != nil {
		return err
	}

	return nil
}

func trimSpace(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[0 : len(s)-1]
	}
	for len(s) > 0 && s[len(s)-1] == '\r' {
		s = s[0 : len(s)-1]
	}

	return s
}
