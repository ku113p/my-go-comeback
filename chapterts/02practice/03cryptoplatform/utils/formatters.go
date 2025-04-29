package utils

import (
	"strconv"
	"strings"
)

func FloatComma(n float64) string {
	if n < 0 {
		return "-" + FloatComma(-n)
	}

	// Use strconv.FormatFloat to preserve precision up to 17 digits
	s := strconv.FormatFloat(n, 'f', -1, 64)

	// Split into integer and fractional parts
	parts := strings.SplitN(s, ".", 2)
	intPart := parts[0]
	intStrWithCommas := addCommas(intPart)

	if len(parts) == 1 {
		return intStrWithCommas // no fractional part
	}

	fracPart := parts[1]

	// Insert commas every 3 digits in the fractional part
	var fracWithCommas strings.Builder
	for i := 0; i < len(fracPart); i++ {
		if i > 0 && i%3 == 0 {
			fracWithCommas.WriteByte(',')
		}
		fracWithCommas.WriteByte(fracPart[i])
	}

	return intStrWithCommas + "." + fracWithCommas.String()
}

func addCommas(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	var result []string
	for i := 0; i < n; i++ {
		if i != 0 && (n-i)%3 == 0 {
			result = append(result, ",")
		}
		result = append(result, string(s[i]))
	}
	return strings.Join(result, "")
}
