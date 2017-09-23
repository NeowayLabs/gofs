package gofs_test

import (
	"testing"

	"github.com/NeowayLabs/gofs"
)

type fsBuilder func(t *testing.T) gofs.FS

func testFS(t *testing.T, newfs fsBuilder) {
}

func testReadWrite(t *testing.T, newfs fsBuilder) {
}

func testReadWriteAll(t *testing.T, newfs fsBuilder) {
}

func testTruncatingExistentPath(t *testing.T, newfs fsBuilder) {
}

func testCloseTwice(t *testing.T, newfs fsBuilder) {
}

func testRemove(t *testing.T, newfs fsBuilder) {
}

func testRemoveNonExistent(t *testing.T, newfs fsBuilder) {
}
