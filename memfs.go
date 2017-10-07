package gofs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
)

type MemFS struct {
	fs     map[string]*bytes.Buffer
	fslock sync.Mutex
}

func NewMemFS() *MemFS {
	return &MemFS{
		fs: map[string]*bytes.Buffer{},
	}
}

func (m *MemFS) Open(path string) (io.ReadCloser, error) {
	contents, err := m.getcontents(path)
	if err != nil {
		return nil, err
	}
	return ioutil.NopCloser(bytes.NewReader(contents.Bytes())), nil
}

func (m *MemFS) ReadAll(path string) ([]byte, error) {
	contents, err := m.getcontents(path)
	if err != nil {
		return nil, err
	}
	return contents.Bytes(), nil
}

func (m *MemFS) Create(path string) (io.WriteCloser, error) {
	buf := &bytes.Buffer{}
	m.setcontents(path, buf)
	return newWriterNopCloser(buf), nil
}

func (m *MemFS) WriteAll(path string, contents []byte) error {
	buf := &bytes.Buffer{}
	io.Copy(buf, bytes.NewReader(contents))
	m.setcontents(path, buf)
	return nil
}

func (m *MemFS) Remove(path string) error {
	if m.isFile(path) {
		m.lock()
		delete(m.fs, path)
		m.unlock()
		return nil
	}
	return m.removeDir(path)
}

func (m *MemFS) isFile(path string) bool {
	m.lock()
	defer m.unlock()
	_, ok := m.fs[path]
	return ok
}

func (m *MemFS) removeDir(dir string) error {
	m.lock()
	defer m.unlock()

	err := fmt.Errorf("removing non existent path[%s]", dir)
	for storedFile, _ := range m.fs {
		if strings.HasPrefix(storedFile, dir) {
			delete(m.fs, storedFile)
			err = nil
		}
	}

	return err
}

func (m *MemFS) lock() {
	m.fslock.Lock()
}

func (m *MemFS) unlock() {
	m.fslock.Unlock()
}

func (m *MemFS) setcontents(path string, contents *bytes.Buffer) {
	m.lock()
	defer m.unlock()
	m.fs[path] = contents
}

func (m *MemFS) getcontents(path string) (*bytes.Buffer, error) {
	m.lock()
	defer m.unlock()

	contents, ok := m.fs[path]
	if !ok {
		return nil, fmt.Errorf("unable to find file[%s]", path)
	}
	return contents, nil
}

type writerNopCloser struct {
	io.Writer
}

func (writerNopCloser) Close() error { return nil }

func newWriterNopCloser(w io.Writer) io.WriteCloser {
	return writerNopCloser{w}
}
