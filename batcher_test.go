package zero

import (
	"bytes"
	"errors"
	"io"
	"math"
	"sync"
	"testing"
	"time"
)

type maxlengthWriter struct {
	Buffer     bytes.Buffer
	MaxLength  int
	WriteCount int
}

var maxlengthExceeded = errors.New("maximum buffer length exceeded")

func (b *maxlengthWriter) Write(p []byte) (n int, err error) {
	b.WriteCount += 1
	n = int(math.Min(float64(b.MaxLength), float64(b.Buffer.Len()+len(p)))) - b.Buffer.Len()
	if n, err = b.Buffer.Write(p[:n]); err == nil && n < len(p) {
		err = maxlengthExceeded
	}
	return
}

var batcherTests = []struct {
	Messages []string
	Limit    int
}{
	{[]string{"test"}, 0},
	{[]string{"test"}, 5},
	{[]string{"test", "test", "test"}, 15},
	{[]string{"test", "test", "test"}, 10},
}

func TestBatcher(t *testing.T) {
	for itest, test := range batcherTests {
		var messages []*bytesMessage
		for _, s := range test.Messages {
			messages = append(messages, &bytesMessage{Bytes: []byte(s)})
		}
		totalbytes := 0
		writer := &maxlengthWriter{MaxLength: test.Limit}
		unit := batcher{Writer: writer}
		for _, m := range messages {
			totalbytes += len(m.Bytes)
			m.Add(1)
			unit.delay(m)
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
		for i, m := range messages {
			done := make(chan bool)
			go func() { m.Wait(); done <- true }()
			bytecount += len(m.Bytes)
			if bytecount <= test.Limit {
				select {
				case <-done:
					if m.Err != nil {
						t.Errorf("%d: batcher called message[%d].Fail(%s)", itest, i, m.Err.Error())
					}
				case <-time.After(10 * time.Millisecond):
					t.Errorf("%d: batcher did not call message[%d].Done()", itest, i)
				}
			} else {
				select {
				case <-done:
					if m.Err == nil {
						t.Errorf("%d: batcher called message[%d].Done() instead of .Fail()", itest, i)
					}
				case <-time.After(10 * time.Millisecond):
					t.Errorf("%d: batcher called message[%d].Done() nor .Fail()", itest, i)
				}
			}
		}
		if writer.WriteCount != 1 {
			t.Errorf("%d: writer.WriteCount == %d", itest, writer.WriteCount)
		}
	}
}

var batcherConsumeTests = []struct {
	Messages []string
	Limit    int
	Close    bool
}{
	{[]string{"test"}, 0, false},
	{[]string{"test"}, 5, false},
	{[]string{"test", "test", "test"}, 15, false},
	{[]string{"test", "test", "test"}, 10, false},
	{[]string{"test", "test", "test"}, 15, true},
	{[]string{"test", "test", "test"}, 10, true},
}

func TestBatcher_Consume(t *testing.T) {
	for itest, test := range batcherConsumeTests {
		var (
			ch         = make(chan CallbackMessage, len(test.Messages))
			writer     = &maxlengthWriter{MaxLength: test.Limit}
			unit       = batcher{Writer: writer}
			sent       []*bytesMessage
			totalbytes = 0
		)
		for _, s := range test.Messages {
			m := &bytesMessage{Bytes: []byte(s)}
			totalbytes += len(m.Bytes)
			m.Add(1)
			ch <- m
			sent = append(sent, m)
		}
		if test.Close {
			close(ch)
		}
		done := make(chan error)
		go func() {
			for _, m := range sent {
				m.Wait()
			}
			done <- nil
		}()
		var err error
		go func() {
			err = unit.Consume(ch, 0)
		}()
		select {
		case <-done:
		case <-time.After(10 * time.Millisecond):
			t.Errorf("%d: timed out waiting on messages.", itest)
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
	Messages []string
	Reps     int
	Limit    int
}{
	{[]string{"test"}, 1, 5},
	{[]string{"test"}, 2, 15},
	{[]string{"1", "2", "3", "4"}, 10, 80},
}

func TestBatchWriter_Write(t *testing.T) {
	for itest, test := range batchWriteTests {
		writer := maxlengthWriter{MaxLength: test.Limit}
		unit := BatchWriter(&writer)
		var expected string
		var wg sync.WaitGroup
		for i := 0; i < test.Reps; i += 1 {
			for _, s := range test.Messages {
				expected += s
				wg.Add(1)
				go func() {
					defer wg.Done()
					if _, err := unit.Write([]byte(s)); err != nil {
						t.Errorf("%d: write returned error: %s", itest, err.Error())
					}
				}()
			}
		}
		wg.Wait()
		actual := string(writer.Buffer.Bytes())
		if len(actual) != len(expected) {
			t.Errorf("%d: expected length %d, got %v", itest, len(expected), actual)
		}
		if writer.WriteCount != 1 && writer.WriteCount == test.Reps*len(test.Messages) {
			t.Errorf("%d: writer.WriteCount is exactly %d", itest, writer.WriteCount)
		}
	}
}
