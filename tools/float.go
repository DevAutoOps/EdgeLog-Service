package tools

import (
	"math"
	"strconv"
)

//  Keep decimal places
func FormatFloat(num float64, decimal int) string {
	//  Default multiplication 1
	d := float64(1)
	if decimal > 0 {
		// 10 of N Power
		d = math.Pow10(decimal)
	}
	// math.trunc The function is to return the integer part of a floating point number
	//  Divide it back ， Invalid after decimal point 0 It doesn't exist
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}
