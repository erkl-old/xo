package xo

type Reader interface {
	// Read implements io.Reader.
	Read(buf []byte) (int, error)

	// Peek returns at least n bytes of unread bytes from the Reader's internal
	// buffer, without consuming them, reading more data into the internal
	// buffer first if necessary. The byte slice is only valid until the next
	// read operation.
	//
	// If Peek returns less than n bytes, it must also return an error
	// explaining why.
	Peek(n int) ([]byte, error)

	// Consume discards n bytes from the Reader's internal buffer.
	Consume(n int) error
}

type Writer interface {
	// Write implements io.Writer.
	Write(buf []byte) (int, error)

	// Reserve returns at least n bytes of scratch space from the Writer's
	// internal buffer, flushing data to make room if necessary. The scratch
	// space is only valid until the next write operation.
	//
	// If Reserve returns less than n bytes, it must also return an error
	// explaining why.
	Reserve(n int) ([]byte, error)

	// Commit commits the first n bytes of scratch space to be written.
	Commit(n int) error

	// Flush writes all buffered data.
	Flush() error
}
