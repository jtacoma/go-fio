package fio

import (
	"io"
)

type Writer chan<- Frame

// NewWriter returns a writer that writes frames to the given io.Writer in
// batches.
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
// WriteFrame will block if the channel is full and panic if it is closed.
//
func (w Writer) WriteFrame(f Frame) {
	w <- f
}

// Write blocks until p is written to the underlying io.Writer.
//
// Write will panic if this writer has been closed.
//
func (w Writer) Write(p []byte) (int, error) {
	ch := make(chan Wrote)
	w <- Callback(BytesFrame(p), ch)
	wrote := <-ch
	return wrote.N, wrote.Err
}

// Close closes the channel; not the underlying io.Writer.
//
// Close will panic if the channel has already been closed.
//
func (w Writer) Close() error {
	close(w)
	return nil
}
