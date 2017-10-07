package gofs_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/NeowayLabs/gofs"
)

type fsBuilder func(t *testing.T) gofs.FS

func testFS(t *testing.T, newfs fsBuilder) {

	t.Run("ReadWrite", func(t *testing.T) {
		testReadWrite(t, newfs)
	})

	t.Run("ReadSameFileMultipleTimes", func(t *testing.T) {
		testReadSameFileMultipleTimes(t, newfs)
	})

	t.Run("ConcurrentReadWrite", func(t *testing.T) {
		testConcurrentReadWrite(t, newfs)
	})

	t.Run("ConcurrentReadWriteSameFile", func(t *testing.T) {
		testConcurrentReadWriteSameFile(t, newfs)
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

func testReadSameFileMultipleTimes(t *testing.T, newfs fsBuilder) {
	// Guarantee orthogonality between open streams
	f := setupWithFile(t, newfs)

	readers := []io.ReadCloser{
		f.Open(t),
		f.Open(t),
		f.Open(t),
		f.Open(t),
		f.Open(t),
		f.Open(t),
	}

	for _, reader := range readers {
		readContents, err := ioutil.ReadAll(reader)
		assertNoError(t, err)
		assertEqualBytes(t, f.contents, readContents)
		reader.Close()
	}
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
	fs := newfs(t)
	dir := newtestpath()
	contents := []byte("echo")
	file1 := dir + "/file1"
	file2 := dir + "/file1"

	assertNoError(t, fs.WriteAll(file1, contents), "writing contents to path[%s]", file1)
	assertNoError(t, fs.WriteAll(file2, contents), "writing contents to path[%s]", file2)

	assertNoError(t, fs.Remove(dir))
	assertFileDontExist(t, fs, file1)
	assertFileDontExist(t, fs, file2)
}

func testRemoveFile(t *testing.T, newfs fsBuilder) {
	f := setupWithFile(t, newfs)

	assertNoError(t, f.fs.WriteAll(f.path, f.contents), "writing contents to path[%s]", f.path)
	assertNoError(t, f.fs.Remove(f.path))

	_, err := f.fs.ReadAll(f.path)
	assertError(t, err, "reading file[%s]", f.path)
}

func testRemoveNonExistentFile(t *testing.T, newfs fsBuilder) {
	fs := newfs(t)
	path := newtestpath()
	assertError(t, fs.Remove(path), "removing file[%s]", path)
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

func testConcurrentReadWrite(t *testing.T, newfs fsBuilder) {
	fs := newfs(t)
	concurrency := 50
	waiter := sync.WaitGroup{}
	waiter.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer waiter.Done()

			path := newtestpath()
			contents := []byte(path)

			testRead := func() {
				reader, err := fs.Open(path)
				assertNoError(t, err)
				defer reader.Close()
				readContents, err := ioutil.ReadAll(reader)
				assertNoError(t, err)
				assertEqualBytes(t, contents, readContents)
			}

			testReadAll := func() {
				readContents, err := fs.ReadAll(path)
				assertNoError(t, err)
				assertEqualBytes(t, contents, readContents)
			}

			writer, err := fs.Create(path)
			assertNoError(t, err)
			n, err := writer.Write(contents)
			assertNoError(t, err)
			if n != len(contents) {
				t.Fatalf("expected to write[%i] wrote[%i]", len(contents), n)
			}

			testRead()
			testReadAll()
		}()
	}

	waiter.Wait()
}

func testConcurrentReadWriteSameFile(t *testing.T, newfs fsBuilder) {
	f := setupWithFile(t, newfs)
	concurrency := 50
	waiter := sync.WaitGroup{}
	waiter.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func() {
			defer waiter.Done()

			testRead := func() {
				reader := f.Open(t)
				defer reader.Close()
				readContents, err := ioutil.ReadAll(reader)
				assertNoError(t, err)
				assertEqualBytes(t, f.contents, readContents)
			}

			testReadAll := func() {
				readContents, err := f.fs.ReadAll(f.path)
				assertNoError(t, err)
				assertEqualBytes(t, f.contents, readContents)
			}

			testRead()
			testReadAll()
		}()
	}

	waiter.Wait()
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

type fixture struct {
	fs       gofs.FS
	path     string
	contents []byte
}

func (f *fixture) Open(t *testing.T) io.ReadCloser {
	r, err := f.fs.Open(f.path)
	assertNoError(t, err)
	return r
}

func setupWithFile(t *testing.T, newfs fsBuilder) fixture {
	fs := newfs(t)
	path := newtestpath()
	contents := []byte(path)

	assertNoError(t, fs.WriteAll(path, contents), "writing contents to path[%s]", path)

	return fixture{
		fs:       fs,
		path:     path,
		contents: contents,
	}
}
