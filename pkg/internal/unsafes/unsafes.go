// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafes

import (
	"unsafe"
)

// String transforms a slice of byte into a string without doing the actual copy
// of the data.
func String(b []byte) string {
	if b == nil {
		return ""
	}

	return *(*string)(unsafe.Pointer(&b))
}

// ByteSlice converts a strings into the equivalent byte slice without doing the
// actual copy of the data. The slice returned by this function may be read-only.
// See examples for more details.
func ByteSlice(s string) []byte {
	if s == "" {
		return nil
	}

	sh := *(*StringHeader)(unsafe.Pointer(&s))
	bh := SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

//go:nosplit

// NoEscape hides a pointer from escape analysis.
func NoEscape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

//go:nosplit

// NoEscapeUintPtr hides a uintptr from escape analysis.
func NoEscapeUintPtr(x uintptr) unsafe.Pointer {
	return unsafe.Pointer(x ^ 0)
}
