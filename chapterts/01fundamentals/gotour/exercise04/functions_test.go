package exercise04

import "testing"

func TestFibonacci(t *testing.T) {
	f := fibonacci()

	// Test the first few Fibonacci numbers
	expected := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	for i, exp := range expected {
		actual := f()
		if actual != exp {
			t.Errorf("Test %d: Expected %d, got %d", i+1, exp, actual)
		}
	}

	// Test that the sequence continues correctly after some calls
	for i := 0; i < 5; i++ {
		f() // Advance the sequence
	}
	expectedContinued := []int{987, 1597, 2584}
	for i, exp := range expectedContinued {
		actual := f()
		if actual != exp {
			t.Errorf("Continued Test %d: Expected %d, got %d", i+1, exp, actual)
		}
	}

	// Test creating a new Fibonacci generator and its initial values
	f2 := fibonacci()
	expectedNew := []int{0, 1, 1, 2}
	for i, exp := range expectedNew {
		actual := f2()
		if actual != exp {
			t.Errorf("New Generator Test %d: Expected %d, got %d", i+1, exp, actual)
		}
	}
}

func BenchmarkFibonacci(b *testing.B) {
	f := fibonacci()
	for n := 0; n < b.N; n++ {
		f()
	}
}
