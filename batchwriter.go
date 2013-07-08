package zero

import (
	"io"
	"sync"
)

type bytesMessage struct {
	sync.WaitGroup
	N     int
	Err   error
	Bytes []byte
}

func (m *bytesMessage) WriteTo(w io.Writer) error {
	_, err := w.Write(m.Bytes)
	return err
}

func (m *bytesMessage) Wrote(n int, err error) {
	m.N = n
	m.Err = err
	m.Done()
}

type batchWriter struct {
	ch chan CallbackMessage
}

func BatchWriter(w io.Writer) io.Writer {
	var (
		b      = batcher{Writer: w}
		ch     = make(chan CallbackMessage)
		result = batchWriter{ch: ch}
	)
	go func() {
		for b.Consume(ch, -1) == nil {
		}
	}()
	return &result
}

func (w *batchWriter) Write(p []byte) (int, error) {
	m := &bytesMessage{Bytes: p}
	m.Add(1)
	w.ch <- m
	m.Wait()
	return m.N, m.Err
}
