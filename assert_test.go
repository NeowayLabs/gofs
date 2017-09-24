package gofs_test

import (
	"fmt"
	"testing"
)

func assertEqualBytes(t *testing.T, expectedContents []byte, contents []byte) {
	t.Helper()

	expectedstr := string(expectedContents)
	gotstr := string(contents)

	if expectedstr != gotstr {
		t.Fatalf("expected contents[%s] but got[%s]", expectedstr, gotstr)
	}
}

func assertNoError(t *testing.T, err error, args ...interface{}) {
	t.Helper()

	if err != nil {
		errmsg := formatErrMsg(args)
		t.Fatalf("unexpected error[%s] %s", err, errmsg)
	}
}

func assertError(t *testing.T, err error, args ...interface{}) {
	t.Helper()

	if err == nil {
		errmsg := formatErrMsg(args)
		t.Fatalf("expected error, got success: %s", errmsg)
	}
}

func formatErrMsg(args []interface{}) string {
	if len(args) == 0 {
		return ""
	}

	if len(args) == 1 {
		return args[0].(string)
	}

	return fmt.Sprintf(args[0].(string), args[1:]...)
}
