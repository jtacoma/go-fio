// See the NOTICE file distributed with this work for
// additional information regarding copyright ownership.
// Joshua Tacoma licenses this file to you under the Apache
// License, Version 2.0 (the "License"); you may not use this
// file except in compliance with the License.  You may obtain
// a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package fio

import (
	"io"
)

// A Frame is a finite sequence of bytes.
//
// Rather than storing an entire array of bytes in memory as such, this
// interface allows an array of bytes to be constructed on demand in a buffer
// where it may be adjacent to any number of other frames, ready to be sent all
// at once.
//
// Zero-length frames are valid.
//
// Frames that return a negative value from Len() have undefined consequences.
//
type Frame interface {
	Len() int
	Read(buf []byte) (int, error)
}

// Wrote represents the results of a call to io.Writer.Write().
//
type Wrote struct {
	N   int
	Err error
}

// BytesFrame represents an array of bytes as a Frame.
//
type BytesFrame []byte

func (f BytesFrame) Len() int { return len(f) }

func (f BytesFrame) Read(buf []byte) (n int, err error) {
	n = copy(buf, f)
	if n == len(f) {
		err = io.EOF
	} else {
		panic("buffer is too small")
	}
	return
}

// StringFrame represents a string as a Frame.
//
type StringFrame string

func (f StringFrame) Len() int { return len([]byte(f)) }
func (f StringFrame) Read(buf []byte) (n int, err error) {
	n = copy(buf, f)
	if n == len(f) {
		err = io.EOF
	} else {
		panic("buffer is too small")
	}
	return
}

type callback struct {
	inner Frame
	C     chan<- Wrote
}

// Callback wraps f so that when it is written to the underlying io.Writer of a
// fio.Writer the result of that io.Writer.Write() call will be sent back on ch.
//
func Callback(f Frame, c chan<- Wrote) Frame {
	return &callback{
		inner: f,
		C:     c,
	}
}
func (f *callback) Len() int                   { return f.inner.Len() }
func (f *callback) Read(p []byte) (int, error) { return f.inner.Read(p) }
