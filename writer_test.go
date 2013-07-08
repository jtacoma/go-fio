package fio

import (
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
