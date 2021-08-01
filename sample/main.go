package main

import "math"

// pi launches n goroutines to compute an
// approximation of pi.
func pi(n float64) float64 {
	// ch := make(chan float64)
	f := 0.0
	for k := 0.0; k < n; k++ {
		f += 4 * math.Pow(-1, k) / (2*k + 1)
	}
	return f
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b uint) uint {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// Euler’s Totient Function
func Phi(n uint) uint {
	var result uint = 1
	var i uint
	for i = 2; i < n; i++ {
		if GCD(i, n) == 1 {
			result++
		}
	}
	return result
}

func main() {
	println("Hello, World!")
	println(Phi(1337))
	println(pi(10000))
}
