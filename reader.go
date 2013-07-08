package fio

import (
	"io"
)

type reader struct {
	inner io.Reader
}

// NewReader returns an io.Reader implementation that uses io.ReadFull() to
// ensure that every call to Read() either fills the provided buffer with bytes
// from the underlying reader or returns an error.
//
func NewReader(r io.Reader) io.Reader {
	return &reader{r}
}

func (r *reader) Read(buf []byte) (int, error) {
	return io.ReadFull(r.inner, buf)
}
