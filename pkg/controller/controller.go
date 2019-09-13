// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"bufio"
	"context"
	"html"
	"io"
	"net/http"
	"strings"
	"text/template"
	"time"
	"unsafe"

	"github.com/go-logr/logr"
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
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/zchee/kt/pkg/cmdoptions"
)

type Controller struct {
	ctrlclient.Client
	Manager ctrlmanager.Manager
	Log     logr.Logger

	clientset   kubernetes.Interface
	ctx         context.Context
	ioStreams   cmdoptions.IOStreams
	concurrency int
	tmpl        *template.Template
}

var _ ctrlreconcile.Reconciler = (*Controller)(nil)

type Options func(*Controller)

func WithConcurrency(concurrency int) Options {
	return func(c *Controller) {
		c.concurrency = concurrency
	}
}

func WithIOStearms(ioStearms cmdoptions.IOStreams) Options {
	return func(c *Controller) {
		c.ioStreams = ioStearms
	}
}

func WithTemplate(tmpl *template.Template) Options {
	return func(c *Controller) {
		c.tmpl = tmpl
	}
}

// New returns a new Controller registered with the Manager.
func NewController(ctx context.Context, mgr ctrlmanager.Manager, opts ...Options) (*Controller, error) {
	lvl := zap.NewAtomicLevelAt(zap.DebugLevel)
	logger := ctrlzap.New(func(o *ctrlzap.Options) {
		o.Level = &lvl
		o.Development = true
	})
	ctrllog.SetLogger(logger)

	c := &Controller{
		Client:      mgr.GetClient(),
		Manager:     mgr,
		Log:         logger.WithName("controller"),
		ctx:         ctx,
		concurrency: 1, // default is 1
	}
	for _, opt := range opts {
		opt(c)
	}

	if err := c.SetupWithManager(mgr); err != nil {
		c.Log.Error(err, "failed to create controller")
		return nil, err
	}

	return c, nil
}

type LogEvent struct {
	Pod       *corev1.Pod
	Container *corev1.Container
	Timestamp *time.Time
	Message   string
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

	c.clientset, err = kubernetes.NewForConfig(c.Manager.GetConfig())
	if err != nil {
		log.Error(err, "failed to new clientset")
		return result, err
	}

	now := metav1.Now()
	logOpts := corev1.PodLogOptions{
		Follow:     true,
		Timestamps: true,
		SinceTime:  &now,
	}

	for i := range pod.Spec.Containers {
		container := pod.Spec.Containers[i]
		logOpts.Container = container.Name

		stream, err := c.clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &logOpts).Stream()
		if err != nil {
			if errStatus, ok := err.(apierrors.APIStatus); ok {
				switch errStatus.Status().Code {
				case http.StatusBadRequest:
					return result, err
				case http.StatusNotFound:
					return result, err
				default:
					return result, err
				}
			}
			return result, nil
		}

		go func(container corev1.Container, logOpts corev1.PodLogOptions) {
			for {
				select {
				default:
					r := bufio.NewReader(stream)
					for {
						l, err := r.ReadBytes('\n')
						if err != nil {
							if err == io.EOF {
								stream.Close()
								break
							}
							return
						}
						line := *(*string)(unsafe.Pointer(&l))

						if len(line) > 0 && line[len(line)-1] == '\n' {
							line = line[0 : len(line)-1]
						}
						for len(line) > 0 && line[len(line)-1] == '\r' {
							line = line[0 : len(line)-1]
						}

						parts := strings.SplitN(line, " ", 2)
						if len(parts) < 2 {
							// TODO: Warn
							return
						}

						timeString, message := parts[0], parts[1]
						timestamp, err := time.Parse(time.RFC3339Nano, timeString)
						if err != nil {
							c.Log.Error(err, "failed to parse timestamp", "timeString", timeString)
							return
						}

						event := LogEvent{
							Pod:       &pod,
							Container: &container,
							Timestamp: &timestamp,
							Message:   html.UnescapeString(message),
						}
						if err := c.tmpl.ExecuteTemplate(c.ioStreams.Out, "line", event); err != nil {
							c.Log.Error(err, "failed to tmpl.Execute", "line", line)
							return
						}
						// fmt.Fprintf(c.ioStreams.Out, "%s %s %s\n", pod.Name, container.Name, line)
					}
				case <-c.ctx.Done():
					stream.Close()
					return
				}
			}
		}(container, logOpts)
	}

	log.Info("end of Reconcile")
	return result, nil
}

func (c *Controller) SetupWithManager(mgr ctrlmanager.Manager) error {
	ctrlOpts := ctrlcontroller.Options{
		Reconciler:              c,
		MaxConcurrentReconciles: c.concurrency,
	}

	return ctrlbuilder.ControllerManagedBy(mgr).For(&corev1.Pod{}).WithOptions(ctrlOpts).Complete(c)
}
