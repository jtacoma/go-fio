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
	"sync"
	"testing"
)

func TestZmtp1RecvSender(t *testing.T) {
	var wg sync.WaitGroup
	r, w := io.Pipe()
	unit := NewZmtp1FrameRecvSender(r, w)
	wg.Add(2)
	go func() {
		defer wg.Done()
		flags, length, reader, err := unit.RecvFrame()
		if err != nil {
			t.Errorf("ReadFrame: %s", err)
		} else {
			if flags != 0 {
				t.Errorf("ReadFrame: flags=%x", flags)
			}
			if length != 4 {
				t.Errorf("ReadFrame: length=%d, expected %d", length, 4)
			}
			b := []byte("TEST")
			n, err := io.ReadFull(reader, b)
			if err != nil {
				t.Errorf("ReadFull: %s", err)
			}
			if n != len(b) {
				t.Errorf("ReadFull: n=%d instead of %d", n, len(b))
			}
			if string(b) != "test" {
				t.Errorf("ReadFull: expected %v, got %v", "test", string(b))
			}
		}
	}()
	go func() {
		defer func() {
			t.Logf("sent! (?)")
			wg.Done()
		}()
		f := []byte("test")
		t.Logf("writing: %s", string(f))
		if writer, err := unit.SendFrame(0, uint64(len(f))); err != nil {
			t.Fatalf("WriteFrame: %s", err)
		} else if _, err = writer.Write(f); err != nil {
			t.Fatalf("writer.Write: %s", err)
		}
		t.Logf("wrote: %s", string(f))
		w.Close()
	}()
	wg.Wait()
}
