// Copyright 2013 Joshua Tacoma
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zio1

import (
	"io"

	"github.com/jtacoma/go-fio"
)

// A Message consists of 1 or more frames.
//
type Message struct {
	F []fio.Frame
}

// ReadMessage reads one Message from r or returns an error.
//
// ReadMessage returns io.ErrUnexpectedEOF if EOF is reached before reading a
// final frame.
//
func ReadMessage(r io.Reader) (*Message, error) {
	var message Message
	for {
		f, err := ReadZFrame(r)
		if f == nil {
			return nil, err
		}
		message.F = append(message.F, f.Body)
		if f.Flags&More == 0 {
			return &message, nil
		} else if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
	}
}

// Len returns the total number of bytes required to represent m.
//
func (m Message) Len() (length int) {
	for _, f := range m.F {
		length += ZFrame{Body: f}.Len()
	}
	return
}

// Read puts a representation of m, as Len() bytes, at the beginning of buf.
//
// Read will panic if m contains zero frames or buf is too small.
//
func (m Message) Read(buf []byte) (n int, err error) {
	var fn int
	if len(m.F) == 0 {
		panic(EmptyMessage)
	}
	for i, f := range m.F {
		zf := ZFrame{Body: f}
		if i+1 < len(m.F) {
			zf.Flags |= More
		}
		fn, err = zf.Read(buf[n:])
		n += fn
		if err != nil {
			break
		}
	}
	return
}
