/*
The math package provides common math functionality.
*/
package math

// Abs returns the absolute value `x`.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // Return correctly abs(-0).
	}
	return x
}
