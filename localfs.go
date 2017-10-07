package gofs

import "io"

type LocalFS struct {
}

func (l *LocalFS) Open(path string) (io.ReadCloser, error) {
	return nil, nil
}

func (l *LocalFS) ReadAll(path string) ([]byte, error) {
	return nil, nil
}

func (l *LocalFS) Create(path string) (io.WriteCloser, error) {
	return nil, nil
}

func (l *LocalFS) WriteAll(path string, contents []byte) error {
	return nil
}

func (l *LocalFS) Remove(path string) error {
	return nil
}

func NewLocalFS() *LocalFS {
	return &LocalFS{}
}
