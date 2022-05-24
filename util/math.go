package util

import (
	"math"
)

func RoundNumber(num float64, prec int) float64 {
	return math.Round(num*math.Pow10(prec)) / math.Pow10(prec)
}
