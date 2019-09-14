// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"text/template"
	"time"
)

type Options struct {
	// kubeconfig and context
	KubeConfig  string
	KubeContext string

	// global filters
	Exclude []string
	Include []string

	// pod filters
	Container        string
	ContainerState   string
	ExcludeContainer string
	Namespace        string
	AllNamespaces    bool
	Selector         string
	Timestamps       bool
	Since            time.Duration
	Concurrency      int

	// misc options
	Lines          int64
	UseColor       string
	TemplateString string
	Template       *template.Template
	Output         string
	Help           bool
}
