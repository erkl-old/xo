package xo

import (
	"io"
)

// NewWriter creates a new Writer, using buf as its internal buffer.
func NewWriter(w io.Writer, buf []byte) Writer {
	return &writer{wr: w, buf: buf}
}

// writer implements the Writer interface.
type writer struct {
	wr  io.Writer
	err error
	buf []byte
	n   int
}

func (w *writer) Write(buf []byte) (int, error) {
	var n int

	// If the internal buffer already has data in it, keep filling it up.
	if w.n > 0 {
		n = copy(w.buf[w.n:], buf)
		w.n += n

		// If the internal buffer is now full, flush it.
		if w.n == len(w.buf) {
			err := w.Flush()
			if err != nil {
				return n, err
			}
		}
	}

	// Are we done already?
	if n == len(buf) {
		return n, nil
	}

	// If the remaining input fit in the internal buffer, simply copy it.
	// Otherwise write it straight to the destination writer.
	if len(buf)-n < len(w.buf) {
		nc := copy(w.buf[w.n:], buf[n:])
		w.n += nc
		return n + nc, nil
	} else {
		nw, err := w.write(buf[n:])
		return n + nw, err
	}
}

func (w *writer) Reserve(n int) ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}

	// Flush the buffer to make room, if necessary.
	if n > len(w.buf)-w.n && w.n > 0 {
		if err := w.Flush(); err != nil {
			return nil, err
		}
	}

	// If we're returning a slice containing fewer than n bytes because the
	// internal buffer isn't large enough, explain this with an error.
	if n > len(w.buf) {
		return w.buf[w.n:], ErrBufferTooSmall
	} else {
		return w.buf[w.n:], nil
	}
}

func (w *writer) Commit(n int) error {
	switch {
	case w.err != nil:
		return w.err
	case n < 0:
		return nil
	case n > len(w.buf)-w.n:
		return ErrInvalidCommitSize
	default:
		w.n += n
		return nil
	}
}

func (w *writer) Flush() error {
	if w.n == 0 || w.err != nil {
		return w.err
	}

	_, err := w.write(w.buf[:w.n])
	if err != nil {
		return err
	}

	w.n = 0
	return nil
}

func (w *writer) write(buf []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	var n int

	for n < len(buf) {
		nw, err := w.wr.Write(buf[n:])
		if nw > 0 {
			n += nw
		} else {
			if err == nil {
				err = io.ErrShortWrite
			}
			w.err = err
			return n, err
		}
	}

	return n, nil
}
