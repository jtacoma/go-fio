// See the NOTICE file distributed with this work for
// additional information regarding copyright ownership.
// Joshua Tacoma licenses this file to you under the Apache
// License, Version 2.0 (the "License"); you may not use this
// file except in compliance with the License.  You may obtain
// a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package fionet

import (
	"io"
	"net"
	"time"

	"github.com/jtacoma/go-fio"
)

type FConn struct {
	inner  net.Conn
	writer fio.Writer
}

func NewFConn(c net.Conn) *FConn {
	if itscool, ok := c.(*FConn); ok {
		return itscool
	}
	fconn := &FConn{
		inner:  c,
		writer: fio.NewWriter(c),
	}
	return fconn
}

func Pipe() (*FConn, *FConn) {
	a, b := net.Pipe()
	return NewFConn(a), NewFConn(b)
}

func (c *FConn) WriteFrame(f fio.Frame) error { return c.writer.WriteFrame(f) }

func (c *FConn) Close() error {
	done := make(chan fio.Wrote)
	c.writer.WriteFrame(fio.Callback(fio.BytesFrame{}, done))
	c.writer.Close()
	<-done
	return c.inner.Close()
}

func (c *FConn) Read(b []byte) (n int, err error)  { return io.ReadFull(c.inner, b) }
func (c *FConn) Write(b []byte) (n int, err error) { return c.writer.Write(b) }

func (c *FConn) LocalAddr() net.Addr                { return c.inner.LocalAddr() }
func (c *FConn) RemoteAddr() net.Addr               { return c.inner.RemoteAddr() }
func (c *FConn) SetDeadline(t time.Time) error      { return c.inner.SetDeadline(t) }
func (c *FConn) SetReadDeadline(t time.Time) error  { return c.inner.SetReadDeadline(t) }
func (c *FConn) SetWriteDeadline(t time.Time) error { return c.inner.SetWriteDeadline(t) }
