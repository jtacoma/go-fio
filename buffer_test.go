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

package fio

import (
	"bytes"
	"errors"
	"io"
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
	Frames []StringFrame
	Limit  int
}{
	{[]StringFrame{"test"}, 0},
	{[]StringFrame{"test"}, 5},
	{[]StringFrame{"test", "test", "test"}, 15},
	{[]StringFrame{"test", "test", "test"}, 10},
}

func TestBatcher(t *testing.T) {
	for itest, test := range bufferTests {
		written := make(chan Wrote, len(test.Frames))
		var frames []Frame
		for _, s := range test.Frames {
			frames = append(frames, Callback(s, written))
		}
		totalbytes := 0
		writer := &testWriter{MaxLength: test.Limit}
		unit := buffer{Writer: writer}
		for _, f := range frames {
			totalbytes += f.Len()
			unit.delay(f)
		}
		if err := unit.flush(); err != nil {
			if totalbytes <= test.Limit {
				t.Errorf("%d: unexpected error from flush(): %s", itest, err.Error())
			}
		} else {
			if totalbytes > test.Limit {
				t.Errorf("%d: no error from flush().", itest)
			}
		}
		bytecount := 0
		for i, f := range frames {
			bytecount += f.Len()
			var wrote Wrote
			select {
			case wrote = <-written:
			case <-time.After(10 * time.Millisecond):
				t.Errorf("%d: did not receive results: %d", itest, i)
			}
			if bytecount <= test.Limit {
				if wrote.Err != nil {
					t.Errorf("%d: error (frame %d?): %s", itest, i, wrote.Err.Error())
				}
			} else {
				if wrote.Err == nil {
					t.Errorf("%d: buffer called frame[%d].Done() instead of .Fail()", itest, i)
				}
			}
		}
		if writer.WriteCount != 1 {
			t.Errorf("%d: writer.WriteCount == %d", itest, writer.WriteCount)
		}
	}
}

var bufferConsumeTests = []struct {
	Frames []StringFrame
	Limit  int
	Close  bool
}{
	{[]StringFrame{"test"}, 0, false},
	{[]StringFrame{"test"}, 5, false},
	{[]StringFrame{"test", "test", "test"}, 15, false},
	{[]StringFrame{"test", "test", "test"}, 10, false},
	{[]StringFrame{"test", "test", "test"}, 15, true},
	{[]StringFrame{"test", "test", "test"}, 10, true},
}

func TestBatcher_Consume(t *testing.T) {
	for itest, test := range bufferConsumeTests {
		var (
			ch         = make(chan Frame, len(test.Frames))
			written    = make(chan Wrote, len(test.Frames))
			writer     = &testWriter{MaxLength: test.Limit}
			unit       = buffer{Writer: writer}
			sent       []Frame
			totalbytes int
		)
		for _, s := range test.Frames {
			f := Callback(s, written)
			totalbytes += f.Len()
			ch <- f
			sent = append(sent, f)
		}
		if test.Close {
			close(ch)
		}
		var err error
		go func() {
			err = unit.consume(ch, 0)
		}()
		for _ = range sent {
			<-written
		}
		if totalbytes > test.Limit {
			if err == nil || err == io.EOF {
				t.Errorf("%d: consume should've returned a non-EOF error.", itest)
			}
		} else if err != nil {
			if test.Close && err != io.EOF {
				t.Errorf("%d: consume should've returned EOF: %s", itest, err.Error())
			} else if !test.Close {
				t.Errorf("%d: consume returned an error: %s", itest, err.Error())
			}
		}
	}
}
