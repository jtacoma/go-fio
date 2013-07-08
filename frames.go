package fio

import (
	"io"
)

// A Frame is a finite sequence of bytes.
//
// Rather than storing the entire array of bytes in memory as such, this
// interface allows a byte array to be effectively constructed, on demand, into
// a buffer where it may be concatenated with any number of other frames before
// being sent all at once.
//
// Zero-length frames are valid.
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

// Callback wraps f so that when it is written to an underlying io.Writer the
// results of the io.Writer.Write() call to be sent back on ch.
//
func Callback(f Frame, c chan<- Wrote) Frame {
	return &callback{
		inner: f,
		C:     c,
	}
}
func (f *callback) Len() int                  { return f.inner.Len() }
func (f *callback) WriteTo(w io.Writer) error { return f.inner.WriteTo(w) }
