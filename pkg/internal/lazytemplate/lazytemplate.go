// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package lazytemplate is a thin wrapper over text/template, allowing the use
// of global template variables without forcing them to be parsed at init.
//
// This package copied from golang/go@e4c3925925d9
package lazytemplate

import (
	"io"
	"sync"
	"text/template"

	"github.com/zchee/kt/pkg/internal/kttesting"
)

// Template is a wrapper around text/template.Template, where the underlying
// template will be parsed the first time it is needed.
type Template struct {
	name string
	text string

	once sync.Once
	tmpl *template.Template
}

func (r *Template) tp() *template.Template {
	r.once.Do(r.build)
	return r.tmpl
}

func (r *Template) build() {
	r.tmpl = template.Must(template.New(r.name).Parse(r.text))
	r.name, r.text = "", ""
}

// Execute applies a parsed template to the specified data object,
// and writes the output to w.
func (r *Template) Execute(w io.Writer, data interface{}) error {
	return r.tp().Execute(w, data)
}

// New creates a new lazy template, delaying the parsing work until it is first
// needed. If the code is being run as part of tests, the template parsing will
// happen immediately.
func New(name, text string) *Template {
	lt := &Template{name: name, text: text}
	if kttesting.InTest() {
		// In tests, always parse the templates early.
		lt.tp()
	}
	return lt
}
