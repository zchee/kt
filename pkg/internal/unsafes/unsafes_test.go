// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafes_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/zchee/kt/pkg/internal/unsafes"
)

func TestByteSlice(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []byte
	}{
		{
			name: "ToByteSlice",
			s:    "ToByteSlice",
			want: []byte("ToByteSlice"),
		},
		{
			name: "empty",
			s:    "",
			want: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(unsafes.ByteSlice(tt.s), tt.want); diff != "" {
				t.Errorf("%s: (-got, +want)\n%s", tt.name, diff)
				return
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want string
	}{
		{
			name: "ToString",
			b:    []byte("ToString"),
			want: "ToString",
		},
		{
			name: "empty",
			b:    nil,
			want: "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if diff := cmp.Diff(unsafes.String(tt.b), tt.want); diff != "" {
				t.Errorf("%s: (-got, +want)\n%s", tt.name, diff)
				return
			}
		})
	}
}
