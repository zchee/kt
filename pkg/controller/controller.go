// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"context"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrlbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	scheme        = runtime.NewScheme()
	controllerLog = ctrllog.Log.WithName("controller")
)

type Controller struct {
	ctrlclient.Client
	ctrlmanager.Manager
	Log logr.Logger

	ctx         context.Context
	concurrency int
}

var _ ctrlreconcile.Reconciler = (*Controller)(nil)

func NewController(ctx context.Context, config *rest.Config, concurrency int) (*Controller, error) {
	kubescheme.AddToScheme(scheme)

	ctrllog.SetLogger(ctrlzap.Logger(true))

	mgrOpts := ctrlmanager.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
	}
	mgr, err := ctrlmanager.New(config, mgrOpts)
	if err != nil {
		controllerLog.Error(err, "unable to create manager")
		return nil, err
	}

	c := &Controller{
		Client:      mgr.GetClient(),
		Manager:     mgr,
		Log:         controllerLog,
		ctx:         ctx,
		concurrency: concurrency,
	}
	if err := c.SetupWithManager(mgr); err != nil {
		controllerLog.Error(err, "unable to create controller")
	}

	return c, nil
}

func (c *Controller) Reconcile(req ctrlreconcile.Request) (ctrlreconcile.Result, error) {
	log := c.Log.WithValues("pod", req.NamespacedName)

	var pod corev1.Pod
	if err := c.Get(c.ctx, req.NamespacedName, &pod); err != nil {
		if ctrlclient.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to get pod")
			return ctrlreconcile.Result{}, err
		}
		return ctrlreconcile.Result{}, nil
	}

	log.Info("pod", "pod", pod)

	return ctrlreconcile.Result{}, nil
}

func (c *Controller) SetupWithManager(mgr ctrlmanager.Manager) error {
	ctrlOpts := ctrlcontroller.Options{
		Reconciler:              c,
		MaxConcurrentReconciles: c.concurrency,
	}
	return ctrlbuilder.ControllerManagedBy(mgr).For(&corev1.Pod{}).WithOptions(ctrlOpts).Complete(c)
}
