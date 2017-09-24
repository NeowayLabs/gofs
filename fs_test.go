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

	t.Run("WriteAllTruncatesExistentPath", func(t *testing.T) {
		testWriteAllTruncatesExistentPath(t, newfs)
	})

	t.Run("RemoveFile", func(t *testing.T) {
		testRemoveFile(t, newfs)
	})

	t.Run("RemoveDir", func(t *testing.T) {
		testRemoveDir(t, newfs)
	})

	t.Run("RemoveNonExistentFile", func(t *testing.T) {
		testRemoveNonExistentFile(t, newfs)
	})

	t.Run("ReadNonExistentFile", func(t *testing.T) {
		testReadNonExistentFile(t, newfs)
	})

	t.Run("OpenNonExistentFile", func(t *testing.T) {
		testOpenNonExistentFile(t, newfs)
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
	fs := newfs(t)
	path := newtestpath()
	expectedContents := []byte(path)

	assertNoError(t, fs.WriteAll(path, expectedContents), "writing contents to path[%s]", path)

	contents, err := fs.ReadAll(path)
	assertNoError(t, err, "reading file[%s]", path)

	assertEqualBytes(t, expectedContents, contents)
}

func testWriteAllTruncatesExistentPath(t *testing.T, newfs fsBuilder) {
	fs := newfs(t)
	path := newtestpath()
	expectedContents := []byte(path)

	incompleteData := expectedContents[0 : len(expectedContents)/2]
	assertNoError(t, fs.WriteAll(path, incompleteData), "writing incomplete contents to path[%s]", path)
	assertNoError(t, fs.WriteAll(path, expectedContents), "writing contents to path[%s]", path)

	contents, err := fs.ReadAll(path)
	assertNoError(t, err, "reading file[%s]", path)

	assertEqualBytes(t, expectedContents, contents)
}

func testRemoveDir(t *testing.T, newfs fsBuilder) {
}

func testRemoveFile(t *testing.T, newfs fsBuilder) {
	fs := newfs(t)
	path := newtestpath()
	contents := []byte(path)

	assertNoError(t, fs.WriteAll(path, contents), "writing contents to path[%s]", path)
	assertNoError(t, fs.Remove(path))

	_, err := fs.ReadAll(path)
	assertError(t, err, "reading file[%s]", path)
}

func testRemoveNonExistentFile(t *testing.T, newfs fsBuilder) {
}

func testOpenNonExistentFile(t *testing.T, newfs fsBuilder) {
	fs := newfs(t)
	path := newtestpath()
	_, err := fs.Open(path)
	assertError(t, err, "opening file[%s]", path)
}

func testReadNonExistentFile(t *testing.T, newfs fsBuilder) {
	fs := newfs(t)
	path := newtestpath()
	_, err := fs.ReadAll(path)
	assertError(t, err, "reading file[%s]", path)
}

func newtestpath() string {
	return fmt.Sprintf(
		"/gofs/file-%d-%d",
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
