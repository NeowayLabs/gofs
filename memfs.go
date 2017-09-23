package gofs

import (
	"bytes"
	"io"
	"io/ioutil"
)

type MemFS struct {
	fs map[string]*bytes.Buffer
}

func (m *MemFS) Open(path string) (io.ReadCloser, error) {
	contents := bytes.NewBuffer(m.fs[path].Bytes())
	return ioutil.NopCloser(contents), nil
}

func (m *MemFS) ReadAll(path string) ([]byte, error) {
	return nil, nil
}

func (m *MemFS) Create(path string) (io.WriteCloser, error) {
	buf := &bytes.Buffer{}
	m.fs[path] = buf
	return newWriterNopCloser(buf), nil
}

func (m *MemFS) WriteAll(path string, contents []byte) error {
	return nil
}

func (m *MemFS) Remove(path string) error {
	return nil
}

func NewMemFS() *MemFS {
	return &MemFS{
		fs: map[string]*bytes.Buffer{},
	}
}

type writerNopCloser struct {
	io.Writer
}

func (writerNopCloser) Close() error { return nil }

func newWriterNopCloser(w io.Writer) io.WriteCloser {
	return writerNopCloser{w}
}
