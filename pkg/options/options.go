// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"text/template"
	"time"

	errors "golang.org/x/xerrors"

	corev1 "k8s.io/api/core/v1"

	regexp "github.com/zchee/kt/pkg/internal/lazyregexp"
)

// Options represents a filtered log options.
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

// Query represents a filtered log regexp queries.
type Query struct {
	PodQuery              *regexp.Regexp
	ContainerState        ContainerState
	ContainerQuery        *regexp.Regexp
	ExcludeContainerQuery *regexp.Regexp
	ExcludeQuery          []*regexp.Regexp
	IncludeQuery          []*regexp.Regexp
}

// ContainerState represents a stete of container.
type ContainerState string

// State of container.
const (
	Running    ContainerState = "running"    // container is running
	Waiting    ContainerState = "waiting"    // container is waiting
	Terminated ContainerState = "terminated" // container is terminated
)

// NewContainerState returns the ContainerState from state.
func NewContainerState(state string) (ContainerState, error) {
	switch ContainerState(state) {
	case Running:
		return Running, nil
	case Waiting:
		return Waiting, nil
	case Terminated:
		return Terminated, nil
	}

	return "", errors.New("containerState should be one of 'running', 'waiting', or 'terminated'")
}

// Match returns whether the match state to cs.
func (cs ContainerState) Match(state corev1.ContainerState) bool {
	switch cs {
	case Running:
		return state.Running != nil
	case Waiting:
		return state.Waiting != nil
	case Terminated:
		return state.Terminated != nil
	default:
		return false
	}
}
