package gofs_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/NeowayLabs/gofs"
)

// TODO:
// test N writes/reads
// test concurrently accessing different files

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
		testRemoveFile(t, newfs)
	})

	t.Run("RemoveNonExistent", func(t *testing.T) {
		testRemoveFileNonExistent(t, newfs)
	})

	t.Run("ReadNonExistent", func(t *testing.T) {
		testReadNonExistent(t, newfs)
	})
}

func testReadWrite(t *testing.T, newfs fsBuilder) {
	fs := newfs(t)
	path := newtestpath()
	writer, err := fs.Create(path)

	assertNoError(t, err, "creating writer for file[%s]", path)
	defer closeIO(t, writer)

	expectedContents := []byte(path)
	n, err := writer.Write(expectedContents)
	assertNoError(t, err, "writing file[%s]", path)

	if n != len(expectedContents) {
		t.Fatal("expected to write[%d], wrote just[%d]", n, len(expectedContents))
	}

	reader, err := fs.Open(path)

	assertNoError(t, err, "creating reader for file[%s]", path)
	defer closeIO(t, reader)

	contents, err := ioutil.ReadAll(reader)
	assertNoError(t, err, "reading file[%s]", path)

	assertEqualBytes(t, expectedContents, contents)
}

func testReadWriteAll(t *testing.T, newfs fsBuilder) {
}

func testTruncatingExistentPath(t *testing.T, newfs fsBuilder) {
}

func testCloseTwice(t *testing.T, newfs fsBuilder) {
}

func testRemoveFile(t *testing.T, newfs fsBuilder) {
}

func testRemoveFileNonExistent(t *testing.T, newfs fsBuilder) {
}

func testReadNonExistent(t *testing.T, newfs fsBuilder) {
}

func newtestpath() string {
	return fmt.Sprintf(
		"gofs-%d-%d",
		time.Now().Unix(),
		rand.Intn(99999999),
	)
}

func closeIO(t *testing.T, closer io.Closer) {
	if closer == nil {
		t.Fatal("unexpected nil closer")
	}
	err := closer.Close()
	if err != nil {
		t.Fatal(err)
	}
}
