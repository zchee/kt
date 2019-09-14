// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controller

import (
	"time"

	color "github.com/zchee/color/v2"
)

// LogEvent represents a Pod log event.
type LogEvent struct {
	// Message is the log message itself
	Message string `json:"message"`

	// PodName of the pod
	PodName string `json:"podName"`

	// ContainerName of the container
	ContainerName string `json:"containerName"`

	// Namespace of the pod
	Namespace string `json:"namespace"`

	// Timestamp of the pod
	Timestamp *time.Time `json:"timestamp,omitempty"`

	PodColor       *color.Color `json:"-"`
	ContainerColor *color.Color `json:"-"`
}
