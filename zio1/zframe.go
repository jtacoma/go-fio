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
	"encoding/binary"
	"io"

	"github.com/jtacoma/go-fio"
)

type Flags byte

const (
	More Flags = 1 << iota
)

// A ZFrame is a Frame that consists of Flags and a Body.
//
type ZFrame struct {
	Flags Flags
	Body  fio.Frame
}

// ReadZFrame reads one ZFrame from r or returns an error.
//
func ReadZFrame(r io.Reader) (f *ZFrame, err error) {
	buf := make([]byte, 1)
	if _, err = io.ReadFull(r, buf); err != nil {
		return
	} else if buf[0] < 0xFF {
		buf = make([]byte, buf[0])
	} else {
		buf = make([]byte, 8)
		if _, err = io.ReadFull(r, buf); err != nil {
			return
		}
		length := binary.BigEndian.Uint64(buf)
		buf = make([]byte, length)
	}
	if _, err = io.ReadFull(r, buf); err == nil {
		f = &ZFrame{
			Flags: Flags(buf[0]),
			Body:  fio.BytesFrame(buf[1:]),
		}
	}
	return
}

// Len returns the total number of bytes required to represent f.
//
// While ZMTP/1.0 defines a frame length as the length of Body + 1, for purposes
// of integration with "fio" the Len() of a ZFrame adds to this the number of
// bytes required to encode its length.
//
func (f ZFrame) Len() (length int) {
	length = 1 + f.Body.Len()
	if length < 255 {
		length += 2
	} else {
		length += 10
	}
	return
}

// Read puts a representation of f, as Len() bytes, at the beginning of buf.
//
// Read will panic if buf is too small.
//
func (f ZFrame) Read(buf []byte) (n int, err error) {
	length := 1 + uint64(f.Body.Len())
	if length < 255 {
		buf[0] = byte(length)
		buf[1] = byte(f.Flags)
		n, err = f.Body.Read(buf[2:])
		n += 2
	} else {
		buf[0] = 0xFF
		binary.BigEndian.PutUint64(buf[1:9], length)
		buf[10] = byte(f.Flags)
		n, err = f.Body.Read(buf[11:])
		n += 11
	}
	return
}
