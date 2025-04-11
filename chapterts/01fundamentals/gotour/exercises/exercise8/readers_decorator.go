package exercise8

import "io"

type Rot13Reader struct {
	r io.Reader
}

func (r Rot13Reader) Read(b []byte) (n int, err error) {
	if n, err = r.r.Read(b); err != nil {
		return 0, err
	}

	for i := 0; i < n; i++ {
		b[i] = rot13(b[i])
	}

	return n, nil
}

func rot13(b byte) byte {
	if b >= 'a' && b <= 'm' || b >= 'A' && b <= 'M' {
		return b + 13
	}
	if b >= 'n' && b <= 'z' || b >= 'N' && b <= 'Z' {
		return b - 13
	}
	return b
}
