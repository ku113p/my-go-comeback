package exercise7

import (
	"io"
	"testing"
)

func TestMyReader_Read(t *testing.T) {
	reader := MyReader{}
	buffer := make([]byte, 10)

	n, err := reader.Read(buffer)

	if err != nil {
		t.Errorf("Read returned an error: %v", err)
	}

	if n != len(buffer) {
		t.Errorf("Read returned incorrect number of bytes read. Expected %d, got %d", len(buffer), n)
	}

	for i := range buffer {
		if buffer[i] != 'A' {
			t.Errorf("Read did not fill the buffer with 'A' at index %d. Got '%c'", i, buffer[i])
		}
	}

	// Test reading into a zero-length buffer
	emptyBuffer := make([]byte, 0)
	nEmpty, errEmpty := reader.Read(emptyBuffer)
	if errEmpty != nil {
		t.Errorf("Read on empty buffer returned an error: %v", errEmpty)
	}
	if nEmpty != 0 {
		t.Errorf("Read on empty buffer should return 0 bytes read, got %d", nEmpty)
	}

	// Test reading into a larger buffer
	largeBuffer := make([]byte, 25)
	nLarge, errLarge := reader.Read(largeBuffer)
	if errLarge != nil {
		t.Errorf("Read on large buffer returned an error: %v", errLarge)
	}
	if nLarge != len(largeBuffer) {
		t.Errorf("Read on large buffer returned incorrect number of bytes read. Expected %d, got %d", len(largeBuffer), nLarge)
	}
	for i := range largeBuffer {
		if largeBuffer[i] != 'A' {
			t.Errorf("Read on large buffer did not fill the buffer with 'A' at index %d. Got '%c'", i, largeBuffer[i])
		}
	}
}

func TestMyReader_Read_EOF(t *testing.T) {
	reader := MyReader{}
	buffer := make([]byte, 5)

	// The MyReader's Read method will always fill the buffer and return nil error.
	// It will never return io.EOF.

	n1, err1 := reader.Read(buffer)
	if err1 != nil {
		t.Errorf("First Read returned an error: %v", err1)
	}
	if n1 != len(buffer) {
		t.Errorf("First Read should have filled the buffer. Read %d bytes, expected %d", n1, len(buffer))
	}

	buffer2 := make([]byte, 3)
	n2, err2 := reader.Read(buffer2)
	if err2 != nil {
		t.Errorf("Second Read returned an error: %v", err2)
	}
	if n2 != len(buffer2) {
		t.Errorf("Second Read should have filled the buffer. Read %d bytes, expected %d", n2, len(buffer2))
	}
}

func TestMyReader_ImplementsReader(t *testing.T) {
	var _ io.Reader = MyReader{}
}
