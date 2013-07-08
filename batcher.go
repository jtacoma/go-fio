package zero

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type Message interface {
	WriteTo(w io.Writer) error
}

type CallbackMessage interface {
	Message
	Wrote(n int, err error)
}

type batcher struct {
	mutex    sync.Mutex
	buffered []CallbackMessage
	indices  []int
	buffer   bytes.Buffer
	Writer   io.Writer
}

func (b *batcher) delay(m CallbackMessage) error {
	if err := m.WriteTo(&b.buffer); err != nil {
		return err
	}
	b.buffered = append(b.buffered, m)
	b.indices = append(b.indices, b.buffer.Len())
	return nil
}

func (b *batcher) flush() error {
	if n, err := b.Writer.Write(b.buffer.Bytes()); err != nil {
		defer func() {
			for i, s := range b.indices {
				if s > n {
					b.buffered = b.buffered[i:]
					break
				}
			}
		}()
		for i, m := range b.buffered {
			if b.indices[i] <= n {
				if i == 0 {
					m.Wrote(b.indices[0], nil)
				} else {
					m.Wrote(b.indices[i]-b.indices[i-1], nil)
				}
			} else if i == 0 {
				m.Wrote(n, err)
			} else if b.indices[i-1] <= n {
				m.Wrote(n-b.indices[i-1], err)
			} else {
				m.Wrote(0, err)
			}
		}
		return err
	} else {
		for i, m := range b.buffered {
			m.Wrote(b.indices[i], nil)
		}
		b.buffered = b.buffered[0:0]
		b.indices = b.indices[0:0]
		return nil
	}
}

func (b *batcher) Consume(ch <-chan CallbackMessage, initialWait time.Duration) (err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	var (
		m  CallbackMessage
		ok bool
	)
	if initialWait < 0 {
		m, ok = <-ch
	} else {
		select {
		case m, ok = <-ch:
			if !ok {
				err = io.EOF
			}
		case <-time.After(initialWait):
		}
	}
	if ok {
		b.delay(m)
		for ok {
			select {
			case m, ok = <-ch:
				if ok {
					b.delay(m)
				} else {
					err = io.EOF
				}
			default:
				ok = false
				break
			}
		}
		writerErr := b.flush()
		if writerErr != nil {
			err = writerErr
		}
	}
	return
}
