// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafes

import (
	"reflect"
	"unsafe"
)

// Slice returns a byte array that points to the given string without a heap allocation.
// The string must be preserved until the  byte arrayis disposed.
func Slice(s string) (p []byte) {
	if s == "" {
		return
	}

	sh := (*reflect.SliceHeader)(unsafe.Pointer(&p))
	sh.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	sh.Len = len(s)
	sh.Cap = len(s)

	return
}

// String returns a string that points to the given byte array without a heap allocation.
// The byte array must be preserved until the string is disposed.
func String(b []byte) (s string) {
	if len(b) == 0 {
		return
	}

	(*reflect.StringHeader)(unsafe.Pointer(&s)).Data = uintptr(unsafe.Pointer(&b[0]))
	(*reflect.StringHeader)(unsafe.Pointer(&s)).Len = len(b)

	return
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
