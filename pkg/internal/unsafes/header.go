// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafes

import "unsafe"

// SliceHeader is the same as reflect.SliceHeader but with unsafe.Pointers to
// guarantee they don't get collected by the GC.
type SliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

// StringHeader is the same as reflect.StringHeader but with unsafe.Pointers to
// guarantee they don't get collected by the GC.
type StringHeader struct {
	Data unsafe.Pointer
	Len  int
}
