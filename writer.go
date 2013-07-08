package fio

import (
	"io"
)

// Writer writes frames to an underlying io.Writer in batches.
//
type Writer chan<- Frame

// NewWriter returns a fio.Writer with w as its underlying io.Writer.
//
func NewWriter(w io.Writer) Writer {
	frames := make(chan Frame)
	buffer := buffer{Writer: w}
	go func() {
		for buffer.consume(frames, -1) == nil {
		}
	}()
	return Writer(frames)
}

// WriteFrame simply enqueues f in the channel.
//
// WriteFrame will block if the channel is full and return an error if it is
// closed.
//
func (w Writer) WriteFrame(f Frame) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	w <- f
	return
}

// Write blocks until p is written to the underlying io.Writer.
//
// Write will block if the channel is full and return an error if it is closed.
//
func (w Writer) Write(p []byte) (int, error) {
	ch := make(chan Wrote)
	if err := w.WriteFrame(Callback(BytesFrame(p), ch)); err != nil {
		return 0, err
	} else {
		wrote := <-ch
		return wrote.N, wrote.Err
	}
}

// Close closes the channel; not the underlying io.Writer.
//
func (w Writer) Close() error {
	defer recover()
	close(w)
	return nil
}
