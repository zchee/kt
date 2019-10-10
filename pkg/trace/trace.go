// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"net/http"
	"sync"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	errors "golang.org/x/xerrors"
)

var (
	Views = []*view.View{
		ochttp.ClientSentBytesDistribution,
		ochttp.ClientReceivedBytesDistribution,
		ochttp.ClientRoundtripLatencyDistribution,
	}
)

var registerOnce sync.Once

// RegisterViews appends few custom views to default views and registers to the defaultWorker.
// This function will called only once.
func RegisterViews(views ...*view.View) error {
	if err := view.Register(append(Views, views...)...); err != nil {
		return errors.Errorf("failed register views: %w", err)
	}

	return nil

}

// Transport is an http.RoundTripper that instruments all outgoing requests with
// OpenCensus stats and tracing.
func Transport() http.RoundTripper {
	return &ochttp.Transport{}
}
