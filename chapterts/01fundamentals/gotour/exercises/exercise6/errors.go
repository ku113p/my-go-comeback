package exercise6

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}

	if x == 0 || x == 1 {
		return x, nil
	}

	z := 1.0
	pZ := z * 2

	for i := 0; math.Abs(z/pZ-1) > 0.0000001; i++ {
		pZ = z
		z -= (z*z - x) / (2 * z)
		fmt.Printf("Iteration %d: z = %f, from: %f\n", i+1, z, pZ)
	}

	return z, nil
}
