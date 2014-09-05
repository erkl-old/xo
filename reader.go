package xo

import (
	"io"
)

// NewReader returns a new Reader, using buf as its internal buffer.
func NewReader(r io.Reader, buf []byte) Reader {
	return &reader{rd: r, buf: buf}
}

// reader implements the Reader interface.
type reader struct {
	rd   io.Reader
	err  error
	buf  []byte
	r, w int
}

func (r *reader) Read(buf []byte) (int, error) {
	// If the buffer is empty, either fill it again, or read straight into the
	// provided buffer – depending on which is larger.
	if r.r == r.w {
		if len(buf) >= len(r.buf) {
			return r.read(buf)
		} else {
			err := r.fill(1)
			if err != nil {
				return 0, err
			}
		}
	}

	n := copy(buf[0:], r.buf[r.r:r.w])
	r.r += n
	return n, nil
}

func (r *reader) Peek(n int) ([]byte, error) {
	err := r.fill(n)
	return r.buf[r.r:r.w], err
}

func (r *reader) Consume(n int) error {
	switch {
	case n < 0:
		return nil
	case n > r.w-r.r:
		return ErrInvalidConsumeSize
	default:
		r.r += n
		return nil
	}
}

func (r *reader) fill(n int) (err error) {
	// We should return as much data as possible, even if we can't
	// fit n bytes in the internal buffer.
	if n > len(r.buf) {
		n = len(r.buf)
		err = ErrBufferTooSmall
	}

	// Return early if we have enough data.
	if r.w-r.r >= n {
		return err
	}

	// If necessary, move buffered data to make more room.
	if r.r == r.w {
		r.r = 0
		r.w = 0
	} else if n > len(r.buf)-r.r {
		r.w = copy(r.buf[0:], r.buf[r.r:r.w])
		r.r = 0
	}

	// Fill the buffer.
	for r.w-r.r < n {
		nr, err := r.read(r.buf[r.w:])
		if err != nil {
			return err
		}
		r.w += nr
	}

	return err
}

func (r *reader) read(buf []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	// Limit read attempts.
	for i := 0; i < 10; i++ {
		n, err := r.rd.Read(buf)
		r.err = err
		if n > 0 {
			return n, nil
		} else if err != nil {
			return 0, err
		}
	}

	r.err = io.ErrNoProgress
	return 0, io.ErrNoProgress
}
