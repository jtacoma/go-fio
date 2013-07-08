package fio

import (
	"io"
)

type batchWriter struct {
	ch chan Frame
}

func Writer(w io.Writer) io.WriteCloser {
	var (
		b      = batcher{Writer: w}
		ch     = make(chan Frame)
		result = batchWriter{ch: ch}
	)
	go func() {
		// Passing -1 here causes Consume() to wait forever (for either a
		// message or the closing of the channel).
		for b.Consume(ch, -1) == nil {
		}
	}()
	return &result
}

func (w *batchWriter) Write(p []byte) (int, error) {
	ch := make(chan Wrote)
	w.ch <- Callback(FrameBytes(p), ch)
	wrote := <-ch
	return wrote.N, wrote.Err
}

func (w *batchWriter) Close() error {
	close(w.ch)
	return nil
}
