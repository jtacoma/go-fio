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

// Writer writes frames to an underlying io.Writer in batches.
//
type Writer chan<- Frame

// NewWriter returns a fio.Writer with w as its underlying io.Writer.
//
func NewWriter(w io.Writer) Writer {
	frames := make(chan Frame)
	buffer := buffer{Writer: w}
	go func() {
		for buffer.consume(frames, -1) == nil {
		}
	}()
	return Writer(frames)
}

// WriteFrame simply enqueues f in the channel.
//
// WriteFrame will block if the channel is full and return an error if it is
// closed.
//
func (w Writer) WriteFrame(f Frame) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	w <- f
	return
}

// Write blocks until p is written to the underlying io.Writer.
//
// Write will block if the channel is full and return an error if it is closed.
//
func (w Writer) Write(p []byte) (int, error) {
	ch := make(chan Wrote)
	if err := w.WriteFrame(Callback(BytesFrame(p), ch)); err != nil {
		return 0, err
	} else {
		wrote := <-ch
		return wrote.N, wrote.Err
	}
}

// Close closes the channel; not the underlying io.Writer.
//
func (w Writer) Close() error {
	defer func() {
		recover()
	}()
	close(w)
	return nil
}
