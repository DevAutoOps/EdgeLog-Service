package tools

import (
	"math"
	"strconv"
)

func IntToString(e int) string {
	return strconv.Itoa(e)
}

func Int64ToString(e int64) string {
	return strconv.FormatInt(e, 10)
}

func Round(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n // TODO +0.5  It's for rounding ， If you don't want this, remove this
}
