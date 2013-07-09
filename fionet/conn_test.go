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

package fionet

import (
	"testing"

	"github.com/jtacoma/go-fio"
)

func TestFConn(t *testing.T) {
	a, b := Pipe()
	go a.WriteFrame(fio.StringFrame("One"))
	go a.WriteFrame(fio.StringFrame("Two"))
	buf := make([]byte, len([]byte("One")))
	b.Read(buf)
	s := string(buf)
	if s != "One" {
		t.Fatal("single-frame message was not received.")
	}
}
