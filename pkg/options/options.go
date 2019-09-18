// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"regexp"
	"text/template"
	"time"
)

type Options struct {
	// global filters
	Exclude []string
	Include []string

	// kubeconfig and context
	KubeConfig  string
	KubeContext string

	// pod filters
	Container        string
	ContainerState   string
	ExcludeContainer string
	Namespaces       []string
	Selector         string
	UseColor         string
	Format           string
	Output           string
	Since            time.Duration
	Concurrency      int

	// misc options
	Lines         int64
	Template      *template.Template
	AllNamespaces bool
	Timestamps    bool

	// query
	PodQuery              *regexp.Regexp
	ContainerQuery        *regexp.Regexp
	ExcludeContainerQuery *regexp.Regexp
	ExcludeQuery          []*regexp.Regexp
	IncludeQuery          []*regexp.Regexp
}
