package fio

import (
	"io"
)

// A Frame is a finite sequence of bytes.
//
// Rather than storing an entire array of bytes in memory as such, this
// interface allows an array of bytes to be constructed on demand in a buffer
// where it may be adjacent to any number of other frames, ready to be sent all
// at once.
//
// Zero-length frames are valid.
//
// Frames that return a negative value from Len() have undefined consequences.
//
type Frame interface {
	Len() int // number of bytes that will be written
	WriteTo(w io.Writer) error
}

// Wrote represents the results of a call to io.Writer.Write().
//
type Wrote struct {
	N   int
	Err error
}

// BytesFrame represents an array of bytes as a Frame.
//
type BytesFrame []byte

func (f BytesFrame) Len() int { return len(f) }
func (f BytesFrame) WriteTo(w io.Writer) error {
	_, err := w.Write(f)
	return err
}

// StringFrame represents a string as a Frame.
//
type StringFrame string

func (f StringFrame) Len() int { return len([]byte(f)) }
func (f StringFrame) WriteTo(w io.Writer) error {
	_, err := w.Write([]byte(f))
	return err
}

type callback struct {
	inner Frame
	C     chan<- Wrote
}

// Callback wraps f so that when it is written to the underlying io.Writer of a
// fio.Writer the result of that io.Writer.Write() call will be sent back on ch.
//
func Callback(f Frame, c chan<- Wrote) Frame {
	return &callback{
		inner: f,
		C:     c,
	}
}
func (f *callback) Len() int                  { return f.inner.Len() }
func (f *callback) WriteTo(w io.Writer) error { return f.inner.WriteTo(w) }
