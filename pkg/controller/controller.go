// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"flag"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

type Controller struct {
	log                  logr.Logger
	metricsAddr          string
	enableLeaderElection bool
}

var (
	scheme = runtime.NewScheme()
)

func NewController() (ctrlmanager.Manager, error) {
	kubescheme.AddToScheme(scheme)

	c := &Controller{
		log: ctrllog.NewDelegatingLogger(ctrllog.NullLogger{}),
	}
	flag.StringVar(&c.metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&c.enableLeaderElection, "enable-leader-election", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrllog.SetLogger(ctrlzap.Logger(true))

	mgr, err := ctrlmanager.New(ctrlconfig.GetConfigOrDie(), ctrlmanager.Options{
		Scheme:             scheme,
		MetricsBindAddress: c.metricsAddr,
		LeaderElection:     c.enableLeaderElection,
	})
	if err != nil {
		c.log.Error(err, "unable to start manager")
		return nil, err
	}

	return mgr, nil
}
