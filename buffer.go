package fio

import (
	"bytes"
	"io"
	"time"
)

type buffer struct {
	Writer  io.Writer
	frames  []Frame
	indices []int
	raw     bytes.Buffer
}

func (b *buffer) delay(f Frame) error {
	if err := f.WriteTo(&b.raw); err != nil {
		return err
	}
	b.frames = append(b.frames, f)
	b.indices = append(b.indices, b.raw.Len())
	return nil
}

func (b *buffer) flush() error {
	defer func() {
		b.raw.Reset()
		b.frames = b.frames[0:0]
		b.indices = b.indices[0:0]
	}()
	if n, err := b.Writer.Write(b.raw.Bytes()); err != nil {
		for i, f := range b.frames {
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
		for i, f := range b.frames {
			if cb, ok := f.(*callback); ok {
				cb.C <- Wrote{b.indices[i], nil}
			}
		}
		return nil
	}
}

func (b *buffer) consume(ch <-chan Frame, initialWait time.Duration) (err error) {
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
