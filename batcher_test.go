package fio

import (
	"bytes"
	"errors"
	"io"
	"math"
	"sync"
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
	if(b.MaxLength>=0){
		n = int(math.Min(float64(b.MaxLength), float64(b.Buffer.Len()+len(p)))) - b.Buffer.Len()
	}else{
		n=len(p)
	}
	if n, err = b.Buffer.Write(p[:n]); err == nil && n < len(p) {
		err = maxlengthExceeded
	}
	return
}

var batcherTests = []struct {
	Frames []stringFrame
	Limit  int
}{
	{[]stringFrame{"test"}, 0},
	{[]stringFrame{"test"}, 5},
	{[]stringFrame{"test", "test", "test"}, 15},
	{[]stringFrame{"test", "test", "test"}, 10},
}

func TestBatcher(t *testing.T) {
	for itest, test := range batcherTests {
		written := make(chan Wrote, len(test.Frames))
		var frames []Frame
		for _, s := range test.Frames {
			frames = append(frames, Callback(s, written))
		}
		totalbytes := 0
		writer := &testWriter{MaxLength: test.Limit}
		unit := batcher{Writer: writer}
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
					t.Errorf("%d: batcher called frame[%d].Done() instead of .Fail()", itest, i)
				}
			}
		}
		if writer.WriteCount != 1 {
			t.Errorf("%d: writer.WriteCount == %d", itest, writer.WriteCount)
		}
	}
}

var batcherConsumeTests = []struct {
	Frames []stringFrame
	Limit  int
	Close  bool
}{
	{[]stringFrame{"test"}, 0, false},
	{[]stringFrame{"test"}, 5, false},
	{[]stringFrame{"test", "test", "test"}, 15, false},
	{[]stringFrame{"test", "test", "test"}, 10, false},
	{[]stringFrame{"test", "test", "test"}, 15, true},
	{[]stringFrame{"test", "test", "test"}, 10, true},
}

func TestBatcher_Consume(t *testing.T) {
	for itest, test := range batcherConsumeTests {
		var (
			ch         = make(chan Frame, len(test.Frames))
			written    = make(chan Wrote, len(test.Frames))
			writer     = &testWriter{MaxLength: test.Limit}
			unit       = batcher{Writer: writer}
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
			err = unit.Consume(ch, 0)
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

var batchWriteTests = []struct {
	Frames []string
	Sets   []int
	Limit  int
	Performance float64
}{
	{[]string{"test"}, []int{1}, 5, 1.0},
	{[]string{"test"}, []int{2}, 15, 1.0},
	{[]string{"1", "2", "3", "4"}, []int{10,10,10,10}, 80, 40.0},
}

func TestBufferedWriter_Write(t *testing.T) {
	for itest, test := range batchWriteTests {
		writer := testWriter{MaxLength: -1}
		unit := NewBufferedWriter(&writer)
		var expected string
		var totalcount int
		var wg sync.WaitGroup
		for iset, reps:=range(test.Sets){
			wait:=time.Duration( iset*50)*time.Millisecond
			for i := 0; i < reps; i += 1 {
				for _, s := range test.Frames {
					expected += s
					totalcount +=1
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
		performance:=float64(totalcount)/ float64(writer.WriteCount)
		if performance<test.Performance{
			t.Errorf("%d: averaged %f frames per write but expected %f", itest, performance, test.Performance)
		}
	}
}
