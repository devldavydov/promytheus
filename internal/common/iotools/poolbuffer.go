package iotools

import "io"

type PoolBuffer interface {
	io.ReadWriter
	Reset()
}
