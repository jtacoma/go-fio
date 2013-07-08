package fio

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type batcher struct {
	mutex    sync.Mutex
	buffered []Frame
	indices  []int
	buffer   bytes.Buffer
	Writer   io.Writer
}

func (b *batcher) delay(f Frame) error {
	if err := f.WriteTo(&b.buffer); err != nil {
		return err
	}
	b.buffered = append(b.buffered, f)
	b.indices = append(b.indices, b.buffer.Len())
	return nil
}

func (b *batcher) flush() error {
	defer func() {
		b.buffer.Reset()
		b.buffered = b.buffered[0:0]
		b.indices = b.indices[0:0]
	}()
	if n, err := b.Writer.Write(b.buffer.Bytes()); err != nil {
		for i, f := range b.buffered {
			if cb, ok := f.(*callback); ok {
				if b.indices[i] <= n {
					if i == 0 {
						cb.C <- Wrote{b.indices[0], nil}
					} else {
						cb.C <- Wrote{b.indices[i] - b.indices[i-1], nil}
					}
				} else if i == 0 {
					cb.C <- Wrote{n, err}
				} else if b.indices[i-1] <= n {
					cb.C <- Wrote{n - b.indices[i-1], err}
				} else {
					cb.C <- Wrote{0, err}
				}
			}
		}
		return err
	} else {
		for i, f := range b.buffered {
			if cb, ok := f.(*callback); ok {
				cb.C <- Wrote{b.indices[i], nil}
			}
		}
		return nil
	}
}

func (b *batcher) Consume(ch <-chan Frame, initialWait time.Duration) (err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	var (
		f  Frame
		ok bool
	)
	if initialWait < 0 {
		f, ok = <-ch
	} else {
		select {
		case f, ok = <-ch:
			if !ok {
				err = io.EOF
			}
		case <-time.After(initialWait):
		}
	}
	if ok {
		err = b.delay(f)
		for ok && err == nil {
			select {
			case f, ok = <-ch:
				if ok {
					err = b.delay(f)
				} else {
					err = io.EOF
				}
			default:
				ok = false
			}
		}
		writerErr := b.flush()
		if writerErr != nil {
			err = writerErr
		}
	}
	return
}
