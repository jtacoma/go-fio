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
	"testing"
)

var limitWriterTests = []struct {
	Len    uint64
	Pieces []string
}{
	{0, []string{}},
	{0, []string{""}},
	{5, []string{"Hello"}},
	{0, []string{"Hello", " ", "World", "!"}},
	{11, []string{"Hello", " ", "World", "!"}},
	{12, []string{"Hello", " ", "World", "!"}},
	{13, []string{"Hello", " ", "World", "!"}},
}

func TestLimitedWriter(t *testing.T) {
	for itest, test := range limitWriterTests {
		var buf = bytes.NewBuffer(make([]byte, 0, test.Len))
		var unit = LimitWriter(buf, test.Len)
		var (
			expectN        = unit.N
			triedN  uint64 = 0
		)
		for ipiece, piece := range test.Pieces {
			var content = []byte(piece)
			n, err := unit.Write(content)
			if n < 0 {
				t.Errorf("%d:%d: n=%d < 0", itest, ipiece, n)
			}
			triedN += uint64(len(content))
			expectN -= uint64(n)
			if triedN == test.Len && err == nil {
				// Perfect finish.
				if unit.N != 0 {
					t.Errorf("%d:%d: filled yet unit.N=%d", itest, ipiece, unit.N)
				} else if expectN != 0 {
					t.Errorf("%d:%d: filled yet expectN=%d", itest, ipiece, unit.N)
				}
			} else if triedN > test.Len {
				// ErrLongWrite
				if err != ErrLongWrite {
					t.Errorf("%d:%d: expected ErrLongWrite, got %v", itest, ipiece, err)
				}
				if unit.N != 0 {
					t.Errorf("%d:%d: ErrLongWrite and unit.N=%d", itest, ipiece, unit.N)
				}
				if test.Len >= triedN {
					t.Errorf("%d:%d: ErrLongWrite when writing %d in limit %d", itest, ipiece, triedN, buf.Len())
				}
			} else if err != nil {
				t.Errorf("%d:%d: unexpected error: %s", itest, ipiece, err)
			} else if triedN >= test.Len {
				t.Errorf("%d:%d: expected ErrLongWrite, got nil", itest, ipiece)
			} else if n != len(content) {
				t.Errorf("%d:%d: expected n=%d, got %d", itest, ipiece, len(content), n)
			} else if unit.N != expectN {
				t.Errorf("%d:%d: expected N=%d, got %d", itest, ipiece, expectN, unit.N)
			}
		}
		if triedN < test.Len {
			if uint64(buf.Len()) != triedN {
				t.Errorf("%d: expected buf.Len=%d, got %d", itest, triedN, buf.Len())
			}
		} else if uint64(buf.Len()) != test.Len {
			t.Errorf("%d: expected buf.Len=%d, got %d", itest, test.Len, buf.Len())
		}
	}
}
