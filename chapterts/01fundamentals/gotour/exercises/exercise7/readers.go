package exercise7

type MyReader struct{}

func (m MyReader) Read(b []byte) (n int, err error) {
	for i := range b {
		n += 1
		b[i] = 'A'
	}

	return
}
