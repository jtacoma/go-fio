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
	"bufio"
	"io"
)

type Encoding interface {
	Len([][]byte) int
	Decode(io.Reader) ([][]byte, error)
	Encode(io.Writer, [][]byte) error
}

type Sender struct {
	writer   *bufio.Writer
	encoding Encoding
	sending  chan cue
	sleeping chan bool
	fault    error
}

type cue struct {
	m [][]byte
	c chan<- error
}

func NewSender(w *bufio.Writer, e Encoding) *Sender {
	s := Sender{
		writer:   w,
		encoding: e,
		sending:  make(chan cue),
		sleeping: make(chan bool, 1),
	}
	s.sleeping <- true
	return &s
}

func (s *Sender) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	close(s.sending)
	return
}

func (s *Sender) Send(f ...[]byte) (err error) {
	return s.SendNotify(nil, f...)
}

func (s *Sender) SendNotify(c chan<- error, f ...[]byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	select {
	case s.sending <- cue{f, c}:
	case <-s.sleeping:
		go s.consumeAvailable()
		s.sending <- cue{f, c}
	}
	return
}

func (s *Sender) consumeAvailable() {
	var (
		cue cue
		ok  bool
	)
	defer func() {
		select {
		case s.sleeping <- true:
		default:
		}
	}()
	if s.fault != nil {
		return
	} else {
		defer func() {
			if s.fault != nil {
				close(s.sending)
			} else {
				s.sleeping <- true
			}
		}()
	}
	for {
		select {
		case cue, ok = <-s.sending:
			if !ok {
				return
			}
		default:
			return
		}
		func() {
			var err error
			defer func() {
				if cue.c != nil {
					cue.c <- err
				}
			}()
			mlen := s.encoding.Len(cue.m)
			if mlen < 0 {
				err = ErrNegativeCount
			} else if mlen > s.writer.Buffered()+s.writer.Available() {
				err = ErrTooLong
			} else if mlen > s.writer.Available() {
				s.writer.Flush()
			}
			if err = s.encoding.Encode(s.writer, cue.m); err != nil {
				s.fault = err
			}
		}()
	}
}
