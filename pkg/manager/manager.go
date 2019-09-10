// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manager

import (
	"k8s.io/apimachinery/pkg/runtime"
	kubescheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	ctrlzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	ctrlmanager "sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	scheme     = runtime.NewScheme()
	managerLog = ctrllog.Log.WithName("manager")
)

// Manager represents a ctrlmanager.Manager.
type Manager struct {
	ctrlmanager.Manager
}

// NewManager returns a new Manager for creating Controllers.
func NewManager(config *rest.Config) (*Manager, error) {
	kubescheme.AddToScheme(scheme)

	ctrllog.SetLogger(ctrlzap.Logger(true))

	mgrOpts := ctrlmanager.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
	}

	mgr, err := ctrlmanager.New(config, mgrOpts)
	if err != nil {
		managerLog.Error(err, "failed to create manager")
		return nil, err
	}

	return &Manager{
		Manager: mgr,
	}, nil
}
