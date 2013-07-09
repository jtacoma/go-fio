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
	"sync"
	"testing"
	"time"
)

var batchWriteTests = []struct {
	Frames      []string
	Sets        []int
	Limit       int
	Performance float64
}{
	{[]string{"test"}, []int{1}, 5, 1.0},
	{[]string{"test"}, []int{2}, 15, 1.0},
	{[]string{"1", "2", "3", "4"}, []int{10, 10, 10, 10}, 80, 40.0},
}

func TestWriter_Write(t *testing.T) {
	for itest, test := range batchWriteTests {
		writer := testWriter{MaxLength: -1}
		unit := NewWriter(&writer)
		var expected string
		var totalcount int
		var wg sync.WaitGroup
		for iset, reps := range test.Sets {
			wait := time.Duration(iset*50) * time.Millisecond
			for i := 0; i < reps; i += 1 {
				for _, s := range test.Frames {
					expected += s
					totalcount += 1
					wg.Add(1)
					go func() {
						defer wg.Done()
						time.Sleep(wait)
						if _, err := unit.Write([]byte(s)); err != nil {
							t.Errorf("%d: write returned error: %s", itest, err.Error())
						}
					}()
				}
			}
		}
		wg.Wait()
		actual := string(writer.Buffer.Bytes())
		if len(actual) != len(expected) {
			t.Errorf("%d: expected length %d, got %v", itest, len(expected), actual)
		}
		performance := float64(totalcount) / float64(writer.WriteCount)
		if performance < test.Performance {
			t.Errorf("%d: averaged %f frames per write but expected %f", itest, performance, test.Performance)
		}
	}
}

func TestWriter_Close(t *testing.T) {
	unit := NewWriter(&bytes.Buffer{})
	if err := unit.Close(); err != nil {
		t.Errorf("close an open writer: %s", err)
	}
	if err := unit.Close(); err != nil {
		t.Errorf("close a closed writer: %s", err)
	}
	if _, err := unit.Write([]byte("what?")); err == nil {
		t.Errorf("write on a closed writer: no error??")
	}
}
