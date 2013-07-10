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

type zio1Flags byte

const (
	More zio1Flags = 1 << iota
)

type zio1 struct{}

var Zio1 Encoding = zio1{}

func (zio1) Len(m [][]byte) (length int) {
	for _, f := range m {
		length += 2 + len(f)
		if length >= 255 {
			length += 8
		}
	}
	return
}

func (zio1) Decode(r io.Reader) (m [][]byte, err error) {
	flags := More
	for flags&More != 0 {
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
			flags = zio1Flags(buf[0])
			m = append(m, buf[1:])
		}
	}
	return
}

func (zio1) Encode(w io.Writer, m [][]byte) (err error) {
	var longbuf []byte
	flags := More
	for iframe, f := range m {
		flen := 1 + uint64(len(f))
		if iframe == len(m)-1 {
			flags = 0
		}
		if flen < 255 {
			if _, err = w.Write([]byte{byte(flen), byte(flags)}); err != nil {
				return
			}
		} else {
			if longbuf == nil {
				longbuf = make([]byte, 10)
			}
			longbuf[0] = 0xFF
			binary.BigEndian.PutUint64(longbuf[1:9], flen)
			longbuf[9] = byte(flags)
			if _, err = w.Write(longbuf); err != nil {
				return
			}
		}
		if _, err = w.Write(f); err != nil {
			return
		}
	}
	return
}
