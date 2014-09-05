package xo

import (
	"errors"
)

var (
	// ErrBufferTooSmall may be used by Reader or Writer implementations to
	// indicate that their internal buffers are too small to fulfill a Peek or
	// Reserve request.
	ErrBufferTooSmall = errors.New("xo: buffer too small")

	// ErrShortPeek and ErrShortReserve describe the case that Reader.Peek
	// or Writer.Reserve returned a smaller byte slice than requested, without
	// also providing an error explaining why.
	ErrShortPeek    = errors.New("xo: short peek")
	ErrShortReserve = errors.New("xo: short reserve")

	// ErrInvalidConsumeSize and ErrInvalidCommitSize may be used by Reader or
	// Writer implementations to indicate that the size argument in a call to
	// Consume or Commit is invalid.
	ErrInvalidConsumeSize = errors.New("xo: invalid consume size")
	ErrInvalidCommitSize  = errors.New("xo: invalid commit size")
)
