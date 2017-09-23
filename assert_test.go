package gofs_test

import (
	"fmt"
	"testing"
)

func assertEqualBytes(t *testing.T, expectedContents []byte, contents []byte) {
	expectedstr := string(expectedContents)
	gotstr := string(contents)

	if expectedstr != gotstr {
		t.Fatalf("expected contents[%s] but got[%s]", expectedstr, gotstr)
	}
}

func assertNoError(t *testing.T, err error, args ...interface{}) {
	t.Helper()

	if err != nil {
		errmsg := ""
		if len(args) > 0 {
			errmsg = fmt.Sprintf(args[0].(string), args[1:]...)
		}
		t.Fatalf("unexpected error[%s] %s", err, errmsg)
	}
}
