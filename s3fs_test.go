package gofs_test

import (
	"testing"

	"github.com/NeowayLabs/gofs"
)

func TestS3FS(t *testing.T) {
	testFS(t, func(t *testing.T) gofs.FS {
		return gofs.NewS3()
	})
}
