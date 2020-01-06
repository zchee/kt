// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trace

import (
	"context"
	"net/http/httptrace"

	otelhttptrace "go.opentelemetry.io/otel/plugin/httptrace"
)

// WithClientTrace returns a new context with
// an embedded otelhttptrace.NewClientTrace based on the parent ctx.
func WithClientTrace(ctx context.Context) context.Context {
	return httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
}
