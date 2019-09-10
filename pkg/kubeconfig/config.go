// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kubeconfig

import (
	"path/filepath"

	errors "golang.org/x/xerrors"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func RestConfig(kc string) (config *rest.Config, err error) {
	kc, err = filepath.Abs(kc)
	if err != nil {
		return nil, errors.Errorf("failed to get %s absolute path: %w", kc, err)
	}

	config, err = clientcmd.BuildConfigFromFlags("", kc)
	if err != nil {
		return nil, errors.Errorf("failed to build rest config from %s: %w", kc, err)
	}

	return config, nil
}
