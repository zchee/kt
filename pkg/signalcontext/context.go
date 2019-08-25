// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package signalcontext provides the context interface which handle of termination OS signals.
package signalcontext

import (
	"context"
	"os"
	"os/signal"
	"time"

	errors "golang.org/x/xerrors"
)

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)

	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. exit directly.
	}()

	return stop
}

// NewContext creates a new context with SetupSignalHandler() as our Done() channel.
func NewContext() context.Context {
	return &signalContext{
		stopCh: SetupSignalHandler(),
	}
}

type signalContext struct {
	stopCh <-chan struct{}
}

// Deadline implements context.Context.
func (*signalContext) Deadline() (deadline time.Time, ok bool) { return }

// Done implements context.Context.
func (sctx *signalContext) Done() <-chan struct{} {
	return sctx.stopCh
}

// Err implements context.Context.
func (sctx *signalContext) Err() error {
	select {
	case _, ok := <-sctx.Done():
		if !ok {
			return errors.New("received a termination signal")
		}
	default:
	}

	return nil
}

// Value implements context.Context.
func (*signalContext) Value(key interface{}) interface{} { return nil }
