package ljconf

import (
	"bufio"
	"bytes"
	"io"
)

// *rcReader is a io.Reader ignoring comment lines
type rcReader struct {
	left []byte
	err  error
	r    *bufio.Reader
}

func (r *rcReader) Read(p []byte) (n int, err error) {
	if len(r.left) > 0 {
		n = copy(p, r.left)
		r.left = r.left[n:]
		p = p[n:]
	}

	if len(r.left) > 0 || len(p) == 0 {
		return n, r.err
	}

	for {
		line, err := r.r.ReadBytes('\n')
		r.err = err
		trimmed := bytes.TrimLeft(line, " \t")
		if bytes.HasPrefix(trimmed, []byte("//")) || bytes.HasPrefix(trimmed, []byte(";")) {
			if r.err != nil {
				continue
			}
			line = nil
		}

		r.left = line

		rd := copy(p, r.left)
		n += rd
		r.left = r.left[rd:]
		p = p[rd:]

		if len(r.left) > 0 || len(p) == 0 || rd == 0 {
			return n, r.err
		}
	}
	// Can't reach here
	return 0, nil
}

func newRcReader(r io.Reader) *rcReader {
	return &rcReader{r: bufio.NewReader(r)}
}
