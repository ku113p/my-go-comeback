package exercise5

import (
	"net"
	"testing"
)

func TestIPAddr_String(t *testing.T) {
	tests := []struct {
		name string
		ip   IPAddr
		want string
	}{
		{
			name: "Basic IPv4 Address",
			ip:   IPAddr{192, 168, 1, 100},
			want: "192.168.1.100",
		},
		{
			name: "All Zeros",
			ip:   IPAddr{0, 0, 0, 0},
			want: "0.0.0.0",
		},
		{
			name: "All 255s",
			ip:   IPAddr{255, 255, 255, 255},
			want: "255.255.255.255",
		},
		{
			name: "Leading Zeros",
			ip:   IPAddr{10, 0, 1, 5},
			want: "10.0.1.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ip.String(); got != tt.want {
				t.Errorf("IPAddr.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Example of a test using the net package for validation
func TestIPAddr_String_NetIPComparison(t *testing.T) {
	tests := []struct {
		name string
		ip   IPAddr
	}{
		{
			name: "Valid IPv4",
			ip:   IPAddr{172, 16, 0, 1},
		},
		{
			name: "Another Valid IPv4",
			ip:   IPAddr{10, 20, 30, 40},
		},
		{
			name: "Loopback",
			ip:   IPAddr{127, 0, 0, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ipStr := tt.ip.String()
			parsedIP := net.ParseIP(ipStr)
			if parsedIP == nil || parsedIP.To4() == nil {
				t.Errorf("IPAddr.String() = %v, which is not a valid IPv4 address according to net.ParseIP", ipStr)
			}
			// Compare the byte values
			parsedBytes := parsedIP.To4()
			for i := 0; i < 4; i++ {
				if tt.ip[i] != parsedBytes[i] {
					t.Errorf("IPAddr.String() resulted in %v, but net.ParseIP parsed it as %v, byte at index %d differs (%d != %d)", ipStr, parsedBytes, i, tt.ip[i], parsedBytes[i])
					return // Exit the inner loop if a difference is found
				}
			}
		})
	}
}
