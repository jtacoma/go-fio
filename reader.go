// See the NOTICE file distributed with this work for
// additional information regarding copyright ownership.
// Joshua Tacoma licenses this file to you under the Apache
// License, Version 2.0 (the "License"); you may not use this
// file except in compliance with the License.  You may obtain
// a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
	if _, already := r.(*reader); already {
		return r
	} else {
		return &reader{r}
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	return io.ReadFull(r.inner, buf)
}
