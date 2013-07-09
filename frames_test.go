package fio

import (
	"testing"
)

var bytesFrames = []BytesFrame{
	BytesFrame{},
	BytesFrame("test"),
}

func TestBytesFrame_Read(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()
	unit := BytesFrame("test")
	short := make([]byte, 2)
	unit.Read(short)
}

func TestBytesFrame_Len(t *testing.T) {
	for itest, test := range bytesFrames {
		actual := test.Len()
		if actual != len(test) {
			t.Errorf("%d: expected %d, got %d.", itest, len(test), actual)
		}
	}
}

func TestStringFrame_Read(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("did not panic")
		}
	}()
	unit := StringFrame("test")
	short := make([]byte, 2)
	unit.Read(short)
}
