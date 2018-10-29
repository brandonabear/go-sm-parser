package parser

import (
	"errors"
	"testing"
)

func TestCheckError(t *testing.T) {
	var tests = []struct {
		err error
	}{
		{errors.New("error")},
		{nil},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("The code did not panic.")
		}
	}()
	for _, test := range tests {
		CheckError(test.err)
	}
}
