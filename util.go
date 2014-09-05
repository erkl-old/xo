package xo

// NewReadWriter combines a Reader and a Writer into a ReadWriter.
func NewReadWriter(r Reader, w Writer) ReadWriter {
	return &readwriter{r, w}
}

type readwriter struct {
	Reader
	Writer
}

// WriteString writes a string to w.
func WriteString(w Writer, s string) (int, error) {
	var n int

	// If s is larger than w's internal buffer, the write has to be
	// split into multiple chunks.
	for len(s) > 0 {
		buf, err := w.Reserve(len(s))
		if len(buf) == 0 {
			if err == nil {
				err = ErrShortReserve
			}
			return n, err
		}

		nc := copy(buf, s)

		if err := w.Commit(nc); err != nil {
			return n, err
		}

		s = s[nc:]
		n += nc
	}

	return n, nil
}
