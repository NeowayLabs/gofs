package gofs_test

import (
	"testing"

	"github.com/NeowayLabs/gofs"
)

func TestMemFS(t *testing.T) {
	testFS(t, func(t *testing.T) gofs.FS {
		return gofs.NewMemFS()
	})
}
