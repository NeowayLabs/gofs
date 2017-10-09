package gofs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type LocalFS struct {
	basedir string
}

func (l *LocalFS) Open(path string) (io.ReadCloser, error) {
	fullpath := filepath.Join(l.basedir, path)
	return os.Open(fullpath)
}

func (l *LocalFS) ReadAll(path string) ([]byte, error) {
	reader, err := l.Open(path)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func (l *LocalFS) Create(path string) (io.WriteCloser, error) {
	fullpath := filepath.Join(l.basedir, path)
	if err := createDirs(fullpath); err != nil {
		return nil, fmt.Errorf("unable to create file[%s], error[%s]", path, err)
	}
	return os.Create(fullpath)
}

func (l *LocalFS) WriteAll(path string, contents []byte) error {
	writer, err := l.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()
	n, err := io.Copy(writer, bytes.NewReader(contents))
	if n != int64(len(contents)) {
		return fmt.Errorf("unable to write all, expected to write [%i] bytes but wrote[%i]", len(contents), n)
	}
	return err
}

func (l *LocalFS) Remove(path string) error {
	fullpath := filepath.Join(l.basedir, path)
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		return fmt.Errorf("error[%s] removing file[%s]", err, path)
	}
	return os.RemoveAll(fullpath)
}

// NewLocalFS creates a new LocalFS where all
// files are opened and created relative to the
// basedir provided. Very much like a chroot on
// the provided basedir.
func NewLocalFS(basedir string) *LocalFS {
	return &LocalFS{basedir: basedir}
}

func createDirs(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, 0766)
}
