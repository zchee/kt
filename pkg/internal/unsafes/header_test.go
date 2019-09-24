// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unsafes

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestHeaders(t *testing.T) {
	tests := []struct {
		Unsafes, Reflect interface{}
		name             string
	}{
		{StringHeader{}, reflect.StringHeader{}, "String"},
		{SliceHeader{}, reflect.SliceHeader{}, "Slice"},
	}

	for _, test := range tests {
		d := reflect.TypeOf(test.Unsafes)
		r := reflect.TypeOf(test.Reflect)
		if d.NumField() != r.NumField() {
			t.Errorf("%sHeader and reflect.%sHeader has different number of fields.", test.name, test.name)
			continue
		}

		for i := 0; i < r.NumField(); i++ {
			df := d.Field(i)
			rf := r.Field(i)
			if df.Name != rf.Name {
				t.Errorf("%sHeader field %d is called %s, expected %s", test.name, i, df.Name, rf.Name)
				break
			}

			if rf.Type == reflect.TypeOf(uintptr(0)) {
				if df.Type != reflect.TypeOf(unsafe.Pointer(nil)) {
					t.Errorf("reflect.%sHeader field %d is a uintptr but %sHeader field %d is a %s, expected unsafe.Pointer", test.name, i, test.name, i, df.Type.String())
				}
				continue
			}

			if df.Type != rf.Type {
				t.Errorf("%sHeader field %d is of type %s, expected %s", test.name, i, df.Type.String(), rf.Type.String())
			}
		}
	}
}
