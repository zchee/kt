// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafes

import (
	"reflect"
	"unsafe"
)

// String transforms a slice of byte into a string without doing the actual copy
// of the data.
func String(b []byte) string {
	if b == nil {
		return ""
	}

	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b)).Data)

	var s string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(p)
	hdr.Len = len(b)

	return s
}

// ByteSlice converts a strings into the equivalent byte slice without doing the
// actual copy of the data. The slice returned by this function may be read-only.
// See examples for more details.
func ByteSlice(s string) []byte {
	if s == "" {
		return nil
	}

	p := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)

	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Data = uintptr(p)
	hdr.Cap = len(s)
	hdr.Len = len(s)

	return b
}

// NoEscape hides a pointer from escape analysis.
//go:nosplit
//go:nocheckptr
//go:linkname NoEscape runtime.noescape
func NoEscape(p unsafe.Pointer) unsafe.Pointer
