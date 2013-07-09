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
	"testing"

	"github.com/jtacoma/go-fio"
)

var readMessageTests = []Message{
	{[]fio.Frame{fio.StringFrame("test")}},
	{[]fio.Frame{fio.BytesFrame{}}},
}

func TestReadMessage(t *testing.T) {
	var buffer bytes.Buffer
	for itest, test := range readMessageTests {
		buffer.ReadFrom(&test)
		t.Logf("%d: buffer: x%s", itest, hex.EncodeToString(buffer.Bytes()))
		if message, err := ReadMessage(&buffer); err != nil {
			t.Errorf("%d: ReadMessage: %s", itest, err.Error())
		} else {
			if len(message.F) != len(test.F) {
				t.Errorf("%d: ReadMessage returned len(F)==%d", itest, len(message.F))
			}
		}
	}
}
