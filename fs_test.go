package gofs_test

import (
	"testing"

	"github.com/NeowayLabs/gofs"
)

type fsBuilder func(t *testing.T) gofs.FS

func testFS(t *testing.T, newfs fsBuilder) {
}
