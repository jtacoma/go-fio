// Copyright 2013 Joshua Tacoma
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fio

import (
	"io"
	"math"
)

// A LimitedWriter writes to W but limits the amount of data written to just N
// bytes. Each call to Write updates N to reflect the new amount remaining.
//
type LimitedWriter struct {
	N uint64
	W io.Writer
}

// LimitWriter returns a Writer that writes to w but stops with ErrOverwrite
// when greater than n bytes are written. The underlying implementation is a
// *LimitedWriter.
//
func LimitWriter(w io.Writer, n uint64) *LimitedWriter {
	if n < 0 {
		panic(ErrNegativeLength)
	}
	return &LimitedWriter{n, w}
}

func (w *LimitedWriter) Write(p []byte) (n int, err error) {
	if w.N == 0 || uint64(len(p)) > w.N {
		n, err = w.W.Write(p[:w.N])
		w.N -= uint64(n)
		if err == nil {
			err = ErrOverwrite
		}
	} else {
		n, err = w.W.Write(p)
	}
	return
}

// LimitReader returns a Reader that reads from r but stops with EOF after n
// bytes. The underlying implementation is a *io.LimitedReader unless n is
// greater than the maximum value representable as an int64, in which case it
// uses io.MultiReader with a pair of *io.LimitReader.
//
func LimitReader(r io.Reader, n uint64) io.Reader {
	if n > uint64(math.MaxInt64) {
		return io.MultiReader(io.LimitReader(r, math.MaxInt64), io.LimitReader(r, int64(n-uint64(math.MaxInt64))))
	} else {
		return io.LimitReader(r, int64(n))
	}
}
