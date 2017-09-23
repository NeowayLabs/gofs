# GoFS

GoFS provides multiple implementations of a generic file
system interface.  Don't get confused with a operational
system file system, the generic interface is aimed at the
Go runtime not an operational system file system.

The name may be misleading, but it means a file system for
the Go runtime, not a file system implemented in Go.

## Why

This project is the consolidation of 3 different implementations
of the same idea (AKA duplication), a way to abstract common
storage operations. Basically:

* Read files
* Write files
* Delete files

For testing purposes we required in-memory and host file system
implementations, in production we used Amazon S3. Them the
day came where we needed to migrate all our storage implementations
to Azure. This seemed like a good sign that it was time to
remove the duplication and consolidate the library.

We provide a common interface to access the following storages:

* In Memory
* Local
* Amazon S3
* Azure Blob Storage

The file system interface has a smaller surface than a usual
file system since we don't need all operations and some cloud services
simply won't provide all those operations for blob storage.
