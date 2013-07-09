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
	"bytes"
	"io"
	"time"
)

type buffer struct {
	Writer  io.Writer
	frames  []Frame
	indices []int
	raw     bytes.Buffer
}

func (b *buffer) delay(f Frame) (err error) {
	b.raw.ReadFrom(f) // TODO: ensure this is not doing extra copying
	b.frames = append(b.frames, f)
	b.indices = append(b.indices, b.raw.Len())
	return nil
}

func (b *buffer) flush() error {
	defer func() {
		b.raw.Reset()
		b.frames = b.frames[0:0]
		b.indices = b.indices[0:0]
	}()
	if n, err := b.Writer.Write(b.raw.Bytes()); err != nil {
		for i, f := range b.frames {
			if cb, ok := f.(*callback); ok {
				if b.indices[i] <= n {
					if i == 0 {
						cb.C <- Wrote{b.indices[0], nil}
					} else {
						cb.C <- Wrote{b.indices[i] - b.indices[i-1], nil}
					}
				} else if i == 0 {
					cb.C <- Wrote{n, err}
				} else if b.indices[i-1] <= n {
					cb.C <- Wrote{n - b.indices[i-1], err}
				} else {
					cb.C <- Wrote{0, err}
				}
			}
		}
		return err
	} else {
		for i, f := range b.frames {
			if cb, ok := f.(*callback); ok {
				cb.C <- Wrote{b.indices[i], nil}
			}
		}
		return nil
	}
}

func (b *buffer) consume(ch <-chan Frame, initialWait time.Duration) (err error) {
	var (
		f  Frame
		ok bool
	)
	if initialWait < 0 {
		f, ok = <-ch
		if !ok {
			err = io.EOF
		}
	} else {
		select {
		case f, ok = <-ch:
			if !ok {
				err = io.EOF
			}
		case <-time.After(initialWait):
		}
	}
	if ok {
		err = b.delay(f)
		for ok && err == nil {
			select {
			case f, ok = <-ch:
				if ok {
					err = b.delay(f)
				} else {
					err = io.EOF
				}
			default:
				ok = false
			}
		}
		writerErr := b.flush()
		if writerErr != nil {
			err = writerErr
		}
	}
	return
}
