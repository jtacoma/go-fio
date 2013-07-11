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

package fio

import (
	"encoding/binary"
	"io"
)

type Flags byte

const (
	More Flags = 1 << iota
)

type Zmtp1ReadFramer struct {
	r io.Reader
}

type Zmtp1WriteFramer struct {
	w io.Writer
}

type Zmtp1ReadWriteFramer struct {
	*Zmtp1ReadFramer
	*Zmtp1WriteFramer
}

func NewZmtp1ReadWriteFramer(r io.Reader, w io.Writer) *Zmtp1ReadWriteFramer {
	return &Zmtp1ReadWriteFramer{
		NewZmtp1ReadFramer(r),
		NewZmtp1WriteFramer(w),
	}
}

func NewZmtp1ReadFramer(r io.Reader) *Zmtp1ReadFramer {
	return &Zmtp1ReadFramer{r}
}

func (z *Zmtp1ReadFramer) ReadFrame() (flags Flags, length uint64, r io.Reader, err error) {
	buf := make([]byte, 10)
	if _, err = io.ReadFull(z.r, buf[:2]); err != nil {
		return
	}
	if buf[0] < 0xFF {
		length = uint64(buf[0])
		flags = Flags(buf[1])
	} else {
		if _, err = io.ReadFull(z.r, buf[2:]); err != nil {
			return
		}
		length = binary.BigEndian.Uint64(buf[1:9])
		flags = Flags(buf[9])
	}
	length -= 1
	r = io.LimitReader(z.r, int64(length)) // TODO: handle uint64
	return
}

func NewZmtp1WriteFramer(w io.Writer) *Zmtp1WriteFramer {
	return &Zmtp1WriteFramer{w}
}

func (z *Zmtp1WriteFramer) WriteFrame(flags Flags, length uint64) (w io.Writer, err error) {
	length += 1
	if length < 255 {
		if _, err = z.w.Write([]byte{byte(length), byte(flags)}); err != nil {
			return
		}
	} else {
		if _, err = z.w.Write([]byte{0xFF}); err != nil {
			return
		}
		longbuf := make([]byte, 8)
		binary.BigEndian.PutUint64(longbuf, uint64(length))
		longbuf[9] = byte(flags)
		if _, err = z.w.Write(longbuf); err != nil {
			return
		}
	}
	w = z.w
	return
}
