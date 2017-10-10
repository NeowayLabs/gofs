package gofs_test

import (
	"io/ioutil"
	"testing"

	"github.com/NeowayLabs/gofs"
)

func TestLocalFS(t *testing.T) {
	testFS(t, func(t *testing.T) gofs.FS {
		dir, err := ioutil.TempDir("", "gofs")
		assertNoError(t, err)
		return gofs.NewLocalFS(dir)
	})
}
