package fio

import (
	"io"
)

type Frame interface {
	Len() int
	WriteTo(w io.Writer) error
}

type Wrote struct {
	N   int
	Err error
}

type bytesFrame []byte

func FrameBytes(b []byte) Frame { return bytesFrame(b) }
func (f bytesFrame) Len() int   { return len(f) }
func (f bytesFrame) WriteTo(w io.Writer) error {
	_, err := w.Write(f)
	return err
}

type stringFrame string

func FrameString(s string) Frame { return stringFrame(s) }
func (f stringFrame) Len() int   { return len([]byte(f)) }
func (f stringFrame) WriteTo(w io.Writer) error {
	_, err := w.Write([]byte(f))
	return err
}

type callback struct {
	inner Frame
	C     chan<- Wrote
}

func Callback(f Frame, c chan<- Wrote) Frame {
	return &callback{
		inner: f,
		C:     c,
	}
}
func (f *callback) Len() int                  { return f.inner.Len() }
func (f *callback) WriteTo(w io.Writer) error { return f.inner.WriteTo(w) }
