// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/go-logr/logr"
	ants "github.com/panjf2000/ants/v2"
	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	"github.com/zchee/kt/pkg/options"
	"github.com/zchee/kt/pkg/stdio"
)

// Controller represents a tail Kubernetes resource logs.
//
// Implements a reconcile.Reconciler.
type Controller struct {
	mgr       ctrlmanager.Manager
	client    ctrlclient.Client
	clientset kubernetes.Interface
	predicate ctrlpredicate.Predicate
	log       logr.Logger

	ioStreams stdio.Streams
	ioMu      sync.Mutex // mutex lock of ioStreams
	gp        *ants.PoolWithFunc
	opts      *options.Options
}

// compile time check whether the Controller implements ctrlreconciler.Reconciler interface.
var _ ctrlreconcile.Reconciler = (*Controller)(nil)

const numWorkers = 64

type gpLogger struct {
	logr.Logger
}

// Printf implements github.com/panjf2000/ants/v2.Logger
func (gl gpLogger) Printf(format string, args ...interface{}) {
	gl.Info(format, args...)
}

// New returns the new Controller registered with the manager.Manager.
func New(ioStreams stdio.Streams, mgr ctrlmanager.Manager, opts *options.Options) (c *Controller, err error) {
	lv := zap.NewAtomicLevelAt(zap.ErrorLevel)
	if opts.Debug {
		lv.SetLevel(zap.DebugLevel)
	}

	zapOpts := []ctrlzap.Opts{
		ctrlzap.WriteTo(ioStreams.ErrOut),
		ctrlzap.Level(&lv),
		ctrlzap.UseDevMode(lv.Enabled(zap.DebugLevel)),
	}
	logger := ctrlzap.New(zapOpts...).WithName("controller")
	ctrllog.SetLogger(logger)

	predicate := &PredicateEventFilter{
		ioStreams:    ioStreams,
		log:          logger.WithName("predicate"),
		isNamespaced: (opts.AllNamespaces || len(opts.Namespaces) > 0),
		query:        opts.Query,
	}

	c = &Controller{
		client:    mgr.GetClient(),
		mgr:       mgr,
		predicate: predicate,
		log:       logger,
		ioStreams: ioStreams,
		opts:      opts,
	}

	workerPanicHandler := func(i interface{}) {
		switch i := i.(type) {
		case nil:
			// nothing to do
		case error:
			c.log.Error(i, "paniced on worker")
		default:
			panic(fmt.Errorf("controller.panic: %v", i))
		}
	}
	gpLogger := gpLogger{Logger: logger}
	gpOpts := []ants.Option{
		ants.WithNonblocking(true),
		ants.WithPreAlloc(true),
		ants.WithPanicHandler(workerPanicHandler),
		ants.WithLogger(gpLogger),
	}
	gp, err := ants.NewPoolWithFunc(numWorkers, c.ReadStream, gpOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create goroutine pool: %w", err)
	}
	c.gp = gp

	c.clientset, err = kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to new clientset: %w", err)
	}

	if err := c.SetupWithManager(mgr); err != nil {
		c.log.Error(err, "failed to setup controller with manager", "Controller", c)
		return nil, err
	}

	return c, nil
}

// SetupWithManager setups the Controller with manager.Manager.
func (c *Controller) SetupWithManager(mgr ctrlmanager.Manager) (err error) {
	ctrlOpts := ctrlcontroller.Options{
		MaxConcurrentReconciles: c.opts.Concurrency,
	}

	return ctrlbuilder.ControllerManagedBy(mgr).For(&corev1.Pod{}).WithEventFilter(c.predicate).WithOptions(ctrlOpts).Complete(c)
}

// Reconcile implements a ctrlreconcile.Reconciler.
func (c *Controller) Reconcile(ctx context.Context, req ctrlreconcile.Request) (result ctrlreconcile.Result, err error) {
	log := c.log.WithName("Reconcile").WithValues("req.Namespace", req.Namespace, "req.Name", req.Name)

	var pod corev1.Pod
	if err := c.client.Get(ctx, req.NamespacedName, &pod); err != nil {
		if ctrlclient.IgnoreNotFound(err) != nil {
			log.Error(err, "failed to get pod")
			return result, err
		}
		return result, nil
	}
	if !c.opts.Query.PodQuery.MatchString(pod.GetName()) {
		return result, nil // skip if not matched PodQuery
	}

	podColor, containerColor := findColors(pod.GetName())

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

		stream, err := c.clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.GetName(), podLogOpts).Stream(ctx)
		if err != nil {
			switch apierrors.ReasonForError(err) {
			case metav1.StatusReasonNotFound:
				return result, nil // ignore NotFound error
			case metav1.StatusReasonBadRequest:
				time.Sleep(boff.GetElapsedTime()) // retry after exponential back off delay
				continue
			}

			// fallthrough
			return result, err
		}

		if err := c.gp.Invoke(&eventStream{
			stream: stream,
			LogEvent: LogEvent{
				PodName:        pod.GetName(),
				ContainerName:  container.Name,
				Namespace:      pod.GetNamespace(),
				PodColor:       podColor,
				ContainerColor: containerColor,
			},
		}); err != nil {
			return result, err
		}
	}

	return result, nil
}

type eventStream struct {
	LogEvent

	stream io.ReadCloser
}

func (c *Controller) ReadStream(v interface{}) {
	es := v.(*eventStream)
	defer es.stream.Close()

	r := bufio.NewReader(es.stream)
	for {
		l, err := r.ReadBytes('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			c.log.Error(err, "failed to ReadBytes")
			return
		}
		line := trimSpace(unsafes.String(l))

		event := es.LogEvent
		event.Message = line
		if !c.opts.AllNamespaces && len(c.opts.Namespaces) == 0 {
			event.Namespace = "" // remove Namespace
		}

		c.ioMu.Lock()
		if err = c.opts.Template.Execute(c.ioStreams.Out, event); err != nil {
			c.log.Error(err, "failed to tmpl.Execute", "event", event)
			return
		}
		c.ioMu.Unlock()
	}
}

// Close closes the goroutine pool.
func (c *Controller) Close() {
	c.gp.Release()
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
