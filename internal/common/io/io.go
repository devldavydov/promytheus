package io

import "io"

// ReadWriteReseter presents interface with io.ReadWriter and Reset.
// Can be used where io.ReadWriter should be placed in sync.Pool.
type ReadWriteReseter interface {
	io.ReadWriter
	Reset()
}
