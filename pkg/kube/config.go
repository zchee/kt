// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kube

import (
	"path/filepath"

	errors "golang.org/x/xerrors"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func RestConfig(kubeconfig string) (config *rest.Config, err error) {
	kubeconfig, err = filepath.Abs(kubeconfig)
	if err != nil {
		return nil, errors.Errorf("failed to get %s absolute path: %w", kubeconfig, err)
	}

	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Errorf("failed to build rest config from %s: %w", kubeconfig, err)
	}

	return config, nil
}
