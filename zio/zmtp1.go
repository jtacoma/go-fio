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

package zio

import (
	"encoding/binary"
	"io"

	"github.com/jtacoma/go-fio"
)

type Flags byte

const (
	More Flags = 1 << iota
)

type Zmtp1FrameRecver struct {
	R io.Reader
}

func NewZmtp1FrameRecver(r io.Reader) FrameRecver {
	return &Zmtp1FrameRecver{r}
}

func (z *Zmtp1FrameRecver) RecvFrame() (flags Flags, length uint64, r io.Reader, err error) {
	buf := make([]byte, 10)
	if _, err = io.ReadFull(z.R, buf[:2]); err != nil {
		return
	}
	if buf[0] < 0xFF {
		length = uint64(buf[0])
		flags = Flags(buf[1])
	} else {
		if _, err = io.ReadFull(z.R, buf[2:]); err != nil {
			return
		}
		length = binary.BigEndian.Uint64(buf[1:9])
		flags = Flags(buf[9])
	}
	length -= 1
	r = fio.LimitReader(z.R, length)
	return
}

type Zmtp1FrameSender struct {
	W io.Writer
}

func NewZmtp1FrameSender(w io.Writer) FrameSender {
	return &Zmtp1FrameSender{w}
}

func (z *Zmtp1FrameSender) SendFrame(flags Flags, length uint64) (w io.Writer, err error) {
	if length < 254 {
		if _, err = z.W.Write([]byte{byte(length + 1), byte(flags)}); err != nil {
			return
		}
	} else {
		if _, err = z.W.Write([]byte{0xFF}); err != nil {
			return
		}
		longbuf := make([]byte, 8)
		binary.BigEndian.PutUint64(longbuf, length+1)
		longbuf[9] = byte(flags)
		if _, err = z.W.Write(longbuf); err != nil {
			return
		}
	}
	w = fio.LimitWriter(z.W, length)
	return
}

type zmtp1FrameRecvSender struct {
	FrameRecver
	FrameSender
}

func NewZmtp1FrameRecvSender(r io.Reader, w io.Writer) FrameRecvSender {
	return &zmtp1FrameRecvSender{
		NewZmtp1FrameRecver(r),
		NewZmtp1FrameSender(w),
	}
}
