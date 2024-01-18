package utils

import "math"

// Go does not mod properly when you do d % m with negative numbers
func ModProperly(d, m int) int {
   var res int = d % m
   if ((res < 0 && m > 0) || (res > 0 && m < 0)) {
      return res + m
   }
   return res
}

// gcd calculates the greatest common divisor of a and b.
func Gcd(a int, b int) int {
	if b == 0 {
		return a
	}
	return Gcd(b, a%b)
}

// lcm calculates the least common multiple of a and b.
func Lcm(a int, b int) int {
	return (a / Gcd(a, b)) * b
}

// abd does the absolute value for ints
func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// pow does powers for ints
func Pow(x, y int) int {
    return int(math.Pow(float64(x), float64(y)))
}