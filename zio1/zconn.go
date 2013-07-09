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
	"net"

	"github.com/jtacoma/go-fio/fionet"
)

type ZConn struct {
	*fionet.FConn
}

func NewZConn(c net.Conn) *ZConn {
	return &ZConn{fionet.NewFConn(c)}
}

func Pipe() (*ZConn, *ZConn) {
	a, b := net.Pipe()
	return NewZConn(a), NewZConn(b)
}

func (c *ZConn) ReadZFrame() (*ZFrame, error) {
	return ReadZFrame(c)
}

func (c *ZConn) ReadMessage() (*Message, error) {
	return ReadMessage(c)
}

func (c *ZConn) WriteZFrame(f *ZFrame) error {
	return c.WriteFrame(f)
}

func (c *ZConn) WriteMessage(m *Message) error {
	return c.WriteFrame(m)
}
