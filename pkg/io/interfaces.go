// Copyright 2019 The kt Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"io"
)

// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
var ErrShortWrite = io.ErrShortWrite

// ErrShortBuffer means that a read required a longer buffer than was provided.
var ErrShortBuffer = io.ErrShortBuffer

// EOF is the error returned by Read when no more input is available.
var EOF = io.EOF

// ErrUnexpectedEOF means that EOF was encountered in the middle of reading a fixed-size block or data structure.
var ErrUnexpectedEOF = io.ErrUnexpectedEOF

// ErrNoProgress is returned by some clients of an io.Reader when
// many calls to Read have failed to return any data or error,
// usually the sign of a broken io.Reader implementation.
var ErrNoProgress = io.ErrNoProgress

// Reader is the interface that wraps the basic Read method.
type Reader = io.Reader

// Writer is the interface that wraps the basic Write method.
type Writer = io.Writer

// Closer is the interface that wraps the basic Close method.
type Closer = io.Closer

// Seeker is the interface that wraps the basic Seek method.
type Seeker = io.Seeker

// ReadWriter is the interface that groups the basic Read and Write methods.
type ReadWriter = io.ReadWriter

// ReadCloser is the interface that groups the basic Read and Close methods.
type ReadCloser = io.ReadCloser

// WriteCloser is the interface that groups the basic Write and Close methods.
type WriteCloser = io.WriteCloser

// ReadWriteCloser is the interface that groups the basic Read, Write and Close methods.
type ReadWriteCloser = io.ReadWriteCloser

// ReadSeeker is the interface that groups the basic Read and Seek methods.
type ReadSeeker = io.ReadSeeker

// WriteSeeker is the interface that groups the basic Write and Seek methods.
type WriteSeeker = io.WriteSeeker

// ReadWriteSeeker is the interface that groups the basic Read, Write and Seek methods.
type ReadWriteSeeker = io.ReadWriteSeeker

// ReaderFrom is the interface that wraps the ReadFrom method.
type ReaderFrom = io.ReaderFrom

// WriterTo is the interface that wraps the WriteTo method.
type WriterTo = io.WriterTo

// ReaderAt is the interface that wraps the basic ReadAt method.
type ReaderAt = io.ReaderAt

// WriterAt is the interface that wraps the basic WriteAt method.
type WriterAt = io.WriterAt

// ByteReader is the interface that wraps the ReadByte method.
type ByteReader = io.ByteReader

// ByteScanner is the interface that adds the UnreadByte method to the
// basic ReadByte method.
type ByteScanner = io.ByteScanner

// ByteWriter is the interface that wraps the WriteByte method.
type ByteWriter = io.ByteWriter

// RuneReader is the interface that wraps the ReadRune method.
type RuneReader = io.RuneReader

// RuneScanner is the interface that adds the UnreadRune method to the
// basic ReadRune method.
type RuneScanner = io.RuneScanner

// StringWriter is the interface that wraps the WriteString method.
type StringWriter = io.StringWriter
