package gofs_test

import (
	"testing"

	"github.com/NeowayLabs/gofs"
)

type fsBuilder func(t *testing.T) gofs.FS

func testFS(t *testing.T, newfs fsBuilder) {

	t.Run("ReadWrite", func(t *testing.T) {
		testReadWrite(t, newfs)
	})

	t.Run("ReadWriteAll", func(t *testing.T) {
		testReadWriteAll(t, newfs)
	})

	t.Run("TruncateExistentPath", func(t *testing.T) {
		testTruncatingExistentPath(t, newfs)
	})

	t.Run("CloseTwice", func(t *testing.T) {
		testCloseTwice(t, newfs)
	})

	t.Run("Remove", func(t *testing.T) {
		testRemove(t, newfs)
	})

	t.Run("RemoveNonExistent", func(t *testing.T) {
		testRemoveNonExistent(t, newfs)
	})

	t.Run("ReadNonExistent", func(t *testing.T) {
		testReadNonExistent(t, newfs)
	})
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

func testReadNonExistent(t *testing.T, newfs fsBuilder) {
}
