// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manager

import (
	"go.uber.org/zap"

	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	scheme = runtime.NewScheme()
)

// Manager represents a ctrlmanager.Manager.
type Manager struct {
	ctrlmanager.Manager
}

type Options = ctrlmanager.Options

// New returns a new Manager for creating Controllers.
func New(config *rest.Config, mgrOpts *ctrlmanager.Options) (*Manager, error) {
	kubescheme.AddToScheme(scheme)

	lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
	logger := ctrlzap.New(func(o *ctrlzap.Options) {
		o.Level = &lvl
		o.Development = true
	}).WithName("manager")
	ctrllog.SetLogger(logger)

	mgrOpts.Scheme = scheme
	mgrOpts.MetricsBindAddress = "0" // force disable the metrics serving

	mgr, err := ctrlmanager.New(config, *mgrOpts)
	if err != nil {
		logger.Error(err, "failed to create manager")
		return nil, err
	}

	return &Manager{
		Manager: mgr,
	}, nil
}
