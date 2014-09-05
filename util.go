package xo

// NewReadWriter combines a Reader and a Writer into a ReadWriter.
func NewReadWriter(r Reader, w Writer) ReadWriter {
	return &readwriter{r, w}
}

type readwriter struct {
	Reader
	Writer
}
