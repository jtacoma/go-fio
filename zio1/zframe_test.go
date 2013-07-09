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

package zio1

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	"github.com/jtacoma/go-fio"
)

type longBody struct {
	Reps  int
	Piece []byte
}

func (b longBody) Len() int {
	return len(b.Piece) * b.Reps
}

func (b longBody) Read(buf []byte) (n int, err error) {
	n = b.Len()
	if len(buf) < n {
		panic(io.ErrShortBuffer)
	}
	for i := 0; i < b.Reps; i += 1 {
		copy(buf[i*len(b.Piece):(i+1)*len(b.Piece)], b.Piece)
	}
	err = io.EOF
	return
}

var readFrameTests = []*ZFrame{
	{0, fio.StringFrame("test")},
	{More, fio.BytesFrame{}},
	{0, longBody{255, []byte{0x00}}},
}

func TestZFrameReader_Read(t *testing.T) {
	for itest, test := range readFrameTests {
		t.Logf("%d: %x +%d bytes", itest, test.Flags, test.Body.Len())
		raw := make([]byte, test.Len())
		test.Read(raw)
		buffer := bytes.NewBuffer(raw)
		t.Logf("%d: x%s", itest, hex.EncodeToString(buffer.Bytes()))
		if frame, err := ReadZFrame(buffer); err != nil {
			t.Errorf("%d: ReadFrame: %s", itest, err.Error())
		} else {
			if frame.Flags != test.Flags {
				t.Errorf("%d: ReadFrame returned Flags==%x", itest, frame.Flags)
			}
			var testBuf, frameBuf bytes.Buffer
			testBuf.ReadFrom(test.Body)
			frameBuf.ReadFrom(frame.Body)
			if frameBuf.Len() != testBuf.Len() {
				t.Errorf("%d: ReadFrame returned Body.Len()==%d (%d)",
					itest, frame.Body.Len(), frameBuf.Len())
			} else {
				testS := string(testBuf.Bytes())
				frameS := string(frameBuf.Bytes())
				if testS != frameS {
					t.Errorf("%d: expected %v but ReadFrame returned Body==%v ",
						itest, testS, frameS)
				}
			}
		}
	}
}
