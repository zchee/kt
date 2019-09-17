// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"bufio"
	"context"
	"fmt"
	iopkg "io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/go-logr/logr"
	color "github.com/zchee/color/v2"
	"go.uber.org/zap"
	errors "golang.org/x/xerrors"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/workqueue"
	ctrlbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	ctrlevent "sigs.k8s.io/controller-runtime/pkg/event"
	ctrlhandler "sigs.k8s.io/controller-runtime/pkg/handler"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
	ctrlpredicate "sigs.k8s.io/controller-runtime/pkg/predicate"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
	ctrlsource "sigs.k8s.io/controller-runtime/pkg/source"

	"github.com/zchee/kt/pkg/internal/unsafes"
	"github.com/zchee/kt/pkg/io"
	"github.com/zchee/kt/pkg/options"
)

// Controller implements a reconcile.Reconciler.
type Controller struct {
	ctrlclient.Client
	controller   ctrlcontroller.Controller
	Manager      ctrlmanager.Manager
	Predicate    ctrlpredicate.Predicate
	EventHandler ctrlhandler.EventHandler
	Log          logr.Logger
	Clientset    kubernetes.Interface

	ctx       context.Context
	ioStreams io.Streams
	ioMu      sync.Mutex
	opts      *options.Options
}

var _ ctrlreconcile.Reconciler = (*Controller)(nil)

// New returns the new Controller registered with the manager.Manager.
func New(ctx context.Context, ioStreams io.Streams, mgr ctrlmanager.Manager, opts *options.Options) (*Controller, error) {
	lvl := zap.NewAtomicLevelAt(zap.DebugLevel)
	logger := ctrlzap.New(func(o *ctrlzap.Options) {
		o.Level = &lvl
		o.Development = true
		o.DestWritter = ioStreams.ErrOut
	})
	ctrllog.SetLogger(logger.WithName("controller"))

	state := new(sync.Map)
	predicateFilter := &PredicatePodEventFilter{
		ioStreams: ioStreams,
		states:    state,
		log:       logger.WithName("predicate"),
	}
	eventHandler := &PodEventHandler{
		ioStreams:    ioStreams,
		states:       state,
		log:          logger.WithName("eventHandler"),
		isNamespaced: (opts.AllNamespaces || len(opts.Namespaces) > 0),
	}

	c := &Controller{
		Client:       mgr.GetClient(),
		Manager:      mgr,
		Predicate:    predicateFilter,
		EventHandler: eventHandler,
		Log:          logger,
		ctx:          ctx,
		ioStreams:    ioStreams,
		opts:         opts,
	}

	if err := c.SetupWithManager(mgr); err != nil {
		c.Log.Error(err, "failed to setup controller with manager", "Controller", c)
		return nil, err
	}

	return c, nil
}

const (
	namespaceFmt       = "%s %s%s » %s\n" // (+|-) Namespace/PodName » ContainerName
	nonNamespaceFmt    = "%s %s » %s\n"   // (+|-) PodName » ContainerName
	namespaceSeparator = "/"

	createMark = "+"
	deleteMark = "-"
)

// PredicatePodEventFilter filters events before they are provided to handler.EventHandlers.
type PredicatePodEventFilter struct {
	ioStreams io.Streams
	log       logr.Logger
	states    *sync.Map
}

var _ ctrlpredicate.Predicate = (*PredicatePodEventFilter)(nil)

// Create implements predicate.Predicate.
func (e *PredicatePodEventFilter) Create(event ctrlevent.CreateEvent) bool {
	return true
}

// Delete implements predicate.Predicate.
func (e *PredicatePodEventFilter) Delete(event ctrlevent.DeleteEvent) bool {
	return true
}

// Update implements predicate.Predicate.
func (e *PredicatePodEventFilter) Update(event ctrlevent.UpdateEvent) bool {
	return true
}

// Generic implements predicate.Predicate.
func (e *PredicatePodEventFilter) Generic(event ctrlevent.GenericEvent) bool {
	return true
}

// PodEventHandler enqueues reconcile.Requests in response to only of pods events.
type PodEventHandler struct {
	ioStreams    io.Streams
	log          logr.Logger
	states       *sync.Map
	isNamespaced bool
}

var _ ctrlhandler.EventHandler = (*PodEventHandler)(nil)

// Create implements handler.EventHandler.
func (e *PodEventHandler) Create(event ctrlevent.CreateEvent, q workqueue.RateLimitingInterface) {
	if event.Meta == nil {
		e.log.Error(nil, "CreateEvent received with no metadata", "event", event)
		return
	}

	pod, ok := event.Object.(*corev1.Pod)
	if !ok {
		return
	}

	mark := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	p, c := findColors(pod.Name)

	printFunc := func(pod *corev1.Pod, container corev1.Container) {
		format := nonNamespaceFmt
		args := []interface{}{mark("+"), p.SprintfFunc()(pod.Name), c.SprintfFunc()(container.Name)}

		if e.isNamespaced {
			format = namespaceFmt
			args = append([]interface{}{args[0], p.SprintfFunc()(pod.Namespace + namespaceSeparator)}, args[1:]...)
		}

		fmt.Fprintf(e.ioStreams.Out, format, args...)
	}

	for i, s := range pod.Status.InitContainerStatuses {
		if s.State.Running != nil {
			printFunc(pod, pod.Spec.InitContainers[i])
		}
	}
	for i, s := range pod.Status.ContainerStatuses {
		if s.State.Running != nil {
			printFunc(pod, pod.Spec.Containers[i])
		}
	}

	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.Meta.GetName(),
		Namespace: event.Meta.GetNamespace(),
	}})
}

// Update implements handler.EventHandler.
func (e *PodEventHandler) Update(event ctrlevent.UpdateEvent, q workqueue.RateLimitingInterface) {
	if event.MetaOld == nil {
		e.log.Error(nil, "UpdateEvent received with no old metadata", "event", event)
	}
	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.MetaOld.GetName(),
		Namespace: event.MetaOld.GetNamespace(),
	}})

	if event.MetaNew == nil {
		e.log.Error(nil, "UpdateEvent received with no new metadata", "event", event)
	}
	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.MetaNew.GetName(),
		Namespace: event.MetaNew.GetNamespace(),
	}})
}

// Delete implements handler.EventHandler.
func (e *PodEventHandler) Delete(event ctrlevent.DeleteEvent, q workqueue.RateLimitingInterface) {
	if event.Meta == nil {
		e.log.Error(nil, "DeleteEvent received with no metadata", "event", event)
		return
	}
	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.Meta.GetName(),
		Namespace: event.Meta.GetNamespace(),
	}})
}

// Generic implements handler.EventHandler.
func (e *PodEventHandler) Generic(event ctrlevent.GenericEvent, q workqueue.RateLimitingInterface) {
	if event.Meta == nil {
		e.log.Error(nil, "GenericEvent received with no metadata", "event", event)
		return
	}
	q.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{
		Name:      event.Meta.GetName(),
		Namespace: event.Meta.GetNamespace(),
	}})
}

func (c *Controller) Watch(ioStreams io.Streams) error {
	return c.controller.Watch(
		&ctrlsource.Kind{
			Type: &corev1.Pod{},
		},
		&PodEventHandler{
			ioStreams:    ioStreams,
			log:          c.Log.WithName("enqueueEventHandler"),
			isNamespaced: (c.opts.AllNamespaces || len(c.opts.Namespaces) > 0),
		},
	)
}

const (
	lineDelim = '\n'
)

// Reconcile implements a ctrlreconcile.Reconciler.
func (c *Controller) Reconcile(req ctrlreconcile.Request) (result ctrlreconcile.Result, err error) {
	log := c.Log.WithName("Reconcile").WithValues("req.Namespace", req.Namespace, "req.Name", req.Name)

	var pod corev1.Pod
	if err := c.Get(c.ctx, req.NamespacedName, &pod); err != nil {
		if ctrlclient.IgnoreNotFound(err) != nil {
			log.Error(err, "failed to get pod")
			return result, err
		}
		return result, nil
	}

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

	podColor, containerColor := findColors(pod.Name)

	boff := backoff.NewExponentialBackOff()
	for i := range pod.Spec.Containers {
		container := pod.Spec.Containers[i]
		podLogOpts := new(corev1.PodLogOptions)
		*podLogOpts = *logOpts // shallow copy
		podLogOpts.Container = container.Name

		stream, err := c.Clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, podLogOpts).Context(c.ctx).Stream()
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

		go func(container corev1.Container, stream iopkg.ReadCloser) {
			r := bufio.NewReader(stream)

			for {
				l, err := r.ReadBytes(lineDelim)
				if err != nil {
					if errors.Is(err, iopkg.EOF) {
						stream.Close()
						return
					}
					c.Log.Error(err, "failed to ReadBytes")
					return
				}
				line := trimSpace(unsafes.String(l))

				var (
					timeString string
					message    string
				)
				parts := strings.SplitN(line, " ", 3)
				switch len(parts) {
				case 2:
					timeString = parts[0]
					message = parts[1]
				case 3:
					timeString = parts[0]
					message = parts[2]
					if c.opts.Timestamps {
						message = parts[1] + " " + message
					}
				default:
					message = line // fellback
				}

				event := &LogEvent{
					Message:        message,
					PodName:        pod.Name,
					ContainerName:  container.Name,
					PodColor:       podColor,
					ContainerColor: containerColor,
				}
				if c.opts.Timestamps {
					timestamp, err := time.Parse(time.RFC3339Nano, timeString)
					if err == nil {
						event.Timestamp = &timestamp // omit error handling
					}
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
		}(container, stream)

		go func() {
			<-c.ctx.Done()
			stream.Close()
		}()
	}

	log.Info("end of Reconcile")
	return result, nil
}

// SetupWithManager setups the Controller with manager.Manager.
func (c *Controller) SetupWithManager(mgr ctrlmanager.Manager) (err error) {
	ctrlOpts := ctrlcontroller.Options{
		Reconciler:              c,
		MaxConcurrentReconciles: c.opts.Concurrency,
	}

	c.controller, err = ctrlbuilder.ControllerManagedBy(mgr).For(&corev1.Pod{}).WithEventFilter(c.Predicate).WithOptions(ctrlOpts).Build(c)
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
