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

type Message struct {
	F []fio.Frame
}

func ReadMessage(r io.Reader) (*Message, error) {
	var message Message
	for {
		f, err := ReadZFrame(r)
		if err != nil {
			return nil, err
		} else {
			message.F = append(message.F, f.Body)
		}
		if !f.More {
			break
		}
	}
	return &message, nil
}

func (m *Message) zframes() (frames []ZFrame) {
	for _, f := range m.F {
		frames = append(frames, ZFrame{More: true, Body: f})
	}
	frames[len(frames)-1].More = false
	return frames
}

func (m *Message) Len() (length int) {
	for _, f := range m.zframes() {
		length += f.Len()
	}
	return
}

func (m *Message) Read(buf []byte) (n int, err error) {
	var fn int
	for _, f := range m.zframes() {
		fn, err = f.Read(buf[n:])
		n += fn
		if err != nil {
			break
		}
	}
	return
}
