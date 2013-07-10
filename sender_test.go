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
	"bytes"
	"errors"
	"math"
	"testing"
	"time"
)

type testWriter struct {
	Buffer     bytes.Buffer
	MaxLength  int
	WriteCount int
}

var maxlengthExceeded = errors.New("maximum buffer length exceeded")

func (b *testWriter) Write(p []byte) (n int, err error) {
	b.WriteCount += 1
	if b.MaxLength >= 0 {
		n = int(math.Min(float64(b.MaxLength), float64(b.Buffer.Len()+len(p)))) - b.Buffer.Len()
	} else {
		n = len(p)
	}
	if n, err = b.Buffer.Write(p[:n]); err == nil && n < len(p) {
		err = maxlengthExceeded
	}
	return
}

var bufferTests = []struct {
	Frames []string
	Limit  int
}{
	{[]string{"test"}, 0},
	{[]string{"test"}, (4 + 2)},
	{[]string{"test"}, (4 + 2) - 1},
	{[]string{"test", "test", "test"}, 3 * (4 + 2)},
	{[]string{"test", "test", "test"}, 3*(4+2) - 1},
	{[]string{"test", "test", "test"}, 2 * (4 + 2)},
}

func TestSender(t *testing.T) {
	for itest, test := range bufferTests {
		func() {
			var frames [][]byte
			for _, s := range test.Frames {
				frames = append(frames, []byte(s))
			}
			totalbytes := Zio1.Len(frames)
			writer := &testWriter{MaxLength: test.Limit}
			buffer := bufio.NewWriter(writer)
			unit := NewSender(buffer, Zio1)
			defer unit.Close()
			errors := make(chan error)
			if err := unit.SendNotify(errors, frames...); err != nil {
				t.Errorf("%d: SendNotify: %s", err)
			}
			var ferrs []error
			select {
			case err := <-errors:
				ferrs = append(ferrs, err)
				if err != nil {
					t.Errorf("%d: received %v", itest, err)
					return
				}
			case <-time.After(10 * time.Millisecond):
				t.Fatalf("%d: timed out waiting on error channel.", itest)
			}
			err := buffer.Flush()
			if err != nil {
				if totalbytes <= test.Limit {
					t.Errorf("%d: unexpected error: %v", itest, err)
					return
				}
			} else {
				if totalbytes > test.Limit {
					t.Errorf("%d: expected error after %d bytes, got none after %d.", itest, test.Limit, totalbytes)
					return
				}
			}
			if writer.WriteCount != 1 {
				t.Errorf("%d: writer.WriteCount == %d", itest, writer.WriteCount)
			}
		}()
	}
}
