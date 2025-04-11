package exercise1

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	pZ := z * 2

	for i := 0; math.Abs(z/pZ-1) > 0.0000001; i++ {
		pZ = z
		z -= (z*z - x) / (2 * z)
		fmt.Printf("Iteration %d: z = %f, from: %f\n", i+1, z, pZ)
	}

	return z
}
