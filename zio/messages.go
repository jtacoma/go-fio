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
	"io"

	"github.com/jtacoma/go-fio"
)

type MsgRecver struct {
	R FrameRecver
}

func RecvMsg(r FrameRecver) *MsgRecver {
	return &MsgRecver{r}
}

func (m *MsgRecver) RecvMsgFrame() (more bool, length uint64, r io.Reader, err error) {
	if m.R == nil {
		err = io.EOF
	} else {
		var flags Flags
		flags, length, r, err = m.R.RecvFrame()
		more = (flags & More) == More
		if err != nil {
			m.R = nil
		} else if !more {
			err = io.EOF
		}
	}
	return
}

type MsgSender struct {
	S FrameSender
}

func SendMsg(s FrameSender) *MsgSender {
	return &MsgSender{s}
}

func (m *MsgSender) SendMsgFrame(more bool, length uint64) (w io.Writer, err error) {
	if m.S == nil {
		err = fio.ErrLongWrite
	} else {
		var flags Flags
		if more {
			flags |= More
		}
		w, err = m.S.SendFrame(flags, length)
	}
	return
}
