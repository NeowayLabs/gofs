package gofs

import "io"

// ReaderFS is the interface that groups the basic Read file system methods.
type ReaderFS interface {
	// Open opens the provided path for reading
	Open(path string) (io.ReadCloser, error)

	// ReadAll opens the path for reading and reads all its contents
	// returning them as an byte array.
	ReadAll(path string) ([]byte, error)
}

// WriterFS is the interface that groups the basic Read file system methods.
type WriterFS interface {
	// Create creates the provided path as a file ready for writing.
	// If the file already exists the file will be truncated.
	Create(path string) (io.WriteCloser, error)

	// Create creates the provided path as a file ready for writing
	// and writes all the contents to it, closing it afterwards.
	WriteAll(path string, contents []byte) error
}

// FS is the interface that groups all file system operations
type FS interface {
	ReaderFS
	WriterFS

	// Remove removes the provided path from the underlying storage.
	Remove(path string) error
}
