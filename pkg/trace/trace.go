// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"context"
	"net/http"
)

// WithClientTrace returns a new context with
// an embedded otelhttptrace.NewClientTrace based on the parent ctx.
//
//nolint:gocritic
func WithClientTrace(ctx context.Context, req *http.Request) context.Context {
	return ctx
	// props := propagation.New(propagation.WithExtractors(trace.TraceContext{}))
	// return propagation.ExtractHTTP(ctx, props, req.Header)
}
