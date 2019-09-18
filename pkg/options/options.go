// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"regexp"
	"text/template"
	"time"

	errors "golang.org/x/xerrors"

	corev1 "k8s.io/api/core/v1"
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

	Query *Query
}

type ContainerState string

const (
	Running    = "running"
	Waiting    = "waiting"
	Terminated = "terminated"
)

func NewContainerState(state string) (ContainerState, error) {
	switch state {
	case Running:
		return Running, nil
	case Waiting:
		return Waiting, nil
	case Terminated:
		return Terminated, nil
	}

	return "", errors.New("containerState should be one of 'running', 'waiting', or 'terminated'")
}

func (cs ContainerState) Match(cState corev1.ContainerState) bool {
	switch cs {
	case Running:
		return cState.Running != nil
	case Waiting:
		return cState.Waiting != nil
	case Terminated:
		return cState.Terminated != nil
	default:
		return false
	}
}

type Query struct {
	PodQuery              *regexp.Regexp
	ContainerState        ContainerState
	ContainerQuery        *regexp.Regexp
	ExcludeContainerQuery *regexp.Regexp
	ExcludeQuery          []*regexp.Regexp
	IncludeQuery          []*regexp.Regexp
}
