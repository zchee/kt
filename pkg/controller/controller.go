// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"context"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	toolscache "k8s.io/client-go/tools/cache"
	ctrlbuilder "sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	ctrlcontroller "sigs.k8s.io/controller-runtime/pkg/controller"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	controllerLog = ctrllog.Log.WithName("controller")
)

type Controller struct {
	ctrlclient.Client
	Manager ctrlmanager.Manager
	Log     logr.Logger

	ctx         context.Context
	concurrency int
}

var _ ctrlreconcile.Reconciler = (*Controller)(nil)

// New returns a new Controller registered with the Manager.
func NewController(ctx context.Context, mgr ctrlmanager.Manager, concurrency int) (*Controller, error) {
	ctrllog.SetLogger(ctrlzap.Logger(true))

	c := &Controller{
		Client:      mgr.GetClient(),
		Manager:     mgr,
		Log:         controllerLog,
		ctx:         ctx,
		concurrency: concurrency,
	}
	if err := c.SetupWithManager(mgr); err != nil {
		controllerLog.Error(err, "failed to create controller")
	}

	return c, nil
}

func (c *Controller) Reconcile(req ctrlreconcile.Request) (result ctrlreconcile.Result, err error) {
	log := c.Log.WithValues("controller", "Reconcile")

	var pod corev1.Pod
	if err := c.Get(c.ctx, req.NamespacedName, &pod); err != nil {
		if ctrlclient.IgnoreNotFound(err) != nil {
			log.Error(err, "failed to get pod")
			return ctrlreconcile.Result{}, err
		}
		return result, nil
	}

	cache := c.Manager.GetCache()
	informer, err := cache.GetInformer(&pod)
	// informer, err := cache.GetInformerForKind(pod.GroupVersionKind())
	if err != nil {
		log.Error(err, "failed to get informer")
		return result, err
	}
	informer.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
		AddFunc:    func(obj interface{}) { log.Info("AddFunc", "obj", obj) },
		UpdateFunc: func(oldObj, newObj interface{}) { log.Info("UpdateFunc", "oldObj", oldObj, "newObj", newObj) },
		DeleteFunc: func(obj interface{}) { log.Info("DeleteFunc", "obj", obj) },
	})

	return result, nil
}

func (c *Controller) SetupWithManager(mgr ctrlmanager.Manager) error {
	ctrlOpts := ctrlcontroller.Options{
		Reconciler:              c,
		MaxConcurrentReconciles: c.concurrency,
	}
	return ctrlbuilder.ControllerManagedBy(mgr).For(&corev1.Pod{}).WithOptions(ctrlOpts).Complete(c)
}
