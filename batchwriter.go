package fio

import (
	"io"
)

type BufferedWriter chan Frame

func NewBufferedWriter(w io.Writer) BufferedWriter {
	var (
		back      = batcher{Writer: w}
		front=make(BufferedWriter)
	)
	go func() {
		// Passing -1 here causes Consume() to wait forever (for either a
		// message or the closing of the channel).
		for back.Consume(front, -1) == nil {
		}
	}()
	return front
}

func (w BufferedWriter) Write(p []byte) (int, error) {
	ch := make(chan Wrote)
	w <- Callback(FrameBytes(p), ch)
	wrote := <-ch
	return wrote.N, wrote.Err
}

func (w BufferedWriter) Close() error {
	close(w)
	return nil
}
