package exercise8

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRot13Reader_Read(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "empty input",
			input: "",
			want:  "",
		},
		{
			name:  "lowercase letters",
			input: "abcdefghijklmnopqrstuvwxyz",
			want:  "nopqrstuvwxyzabcdefghijklm",
		},
		{
			name:  "uppercase letters",
			input: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			want:  "NOPQRSTUVWXYZABCDEFGHIJKLM",
		},
		{
			name:  "mixed case letters",
			input: "Hello World",
			want:  "Uryyb Jbeyq",
		},
		{
			name:  "numbers and symbols",
			input: "12345 !@#$%",
			want:  "12345 !@#$%",
		},
		{
			name:  "mixed characters",
			input: "Test 123 string.",
			want:  "Grfg 123 fgevat.",
		},
		{
			name:  "longer string",
			input: "The quick brown fox jumps over the lazy dog.",
			want:  "Gur dhvpx oebja sbk whzcf bire gur ynml qbt.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := Rot13Reader{r: strings.NewReader(tt.input)}
			output := &bytes.Buffer{}
			_, err := io.Copy(output, reader)

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got := output.String(); got != tt.want {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRot13Reader_Read_PartialBuffer(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz"
	want := "nopqrstuvwxyzabcdefghijklm"
	reader := Rot13Reader{r: strings.NewReader(input)}
	output := make([]byte, 5)
	var got bytes.Buffer

	for {
		n, err := reader.Read(output)
		if n > 0 {
			got.Write(output[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Read() error: %v", err)
		}
	}

	if got.String() != want {
		t.Errorf("Read() with partial buffer got = %v, want %v", got.String(), want)
	}
}

func TestRot13Reader_Read_EmptyBuffer(t *testing.T) {
	reader := Rot13Reader{r: strings.NewReader("test")}
	emptyBuf := make([]byte, 0)
	n, err := reader.Read(emptyBuf)
	if n != 0 {
		t.Errorf("Read() with empty buffer should return n=0, got %d", n)
	}
	if err != nil {
		t.Errorf("Read() with empty buffer should return nil error, got %v", err)
	}
}
