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

const (
	flagsLen int = 1
	shortLen     = 1
	longLen      = 8
)

type ZFrame struct {
	More bool
	Body fio.Frame
}

func ReadZFrame(r io.Reader) (f *ZFrame, err error) {
	buf := make([]byte, 2)
	if _, err = io.ReadFull(r, buf); err != nil {
	} else if buf[0] == 0xFF {
		panic("long messages are not yet supported.")
	} else if buf[1]|0x01 != 0x01 {
		panic("unrecognized flags.")
	} else {
		body := make([]byte, int(buf[0]))
		if _, err = io.ReadFull(r, body); err == nil {
			f = &ZFrame{
				More: buf[1]&0x01 != 0,
				Body: fio.BytesFrame(body),
			}
		}
	}
	return
}

func (f *ZFrame) Len() (length int) {
	length = f.Body.Len()
	if length < 255 {
		length += shortLen + flagsLen
	} else {
		length += longLen + flagsLen
	}
	return
}

func (f *ZFrame) Read(buf []byte) (n int, err error) {
	var flags byte
	if f.More {
		flags |= 0x01
	}
	bodylen := f.Body.Len()
	if bodylen < 255 {
		buf[0] = byte(bodylen)
		buf[1] = flags
		n, err = f.Body.Read(buf[2:])
		n += 2
	} else {
		panic("long frames are not yet supported")
		// TODO: write 0xFF 8-byte-length flags
	}
	return
}
