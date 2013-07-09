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
	"testing"

	"github.com/jtacoma/go-fio"
)

var readFrameTests = []*ZFrame{
	{false, fio.StringFrame("test")},
	{true, fio.BytesFrame{}},
}

func TestZFrameReader_Read(t *testing.T) {
	for itest, test := range readFrameTests {
		var buffer bytes.Buffer
		buffer.ReadFrom(test)
		if frame, err := ReadZFrame(&buffer); err != nil {
			t.Errorf("%d: ReadFrame: %s", itest, err.Error())
		} else {
			if frame.More != test.More {
				t.Errorf("%d: ReadFrame returned More==%v", itest, frame.More)
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
