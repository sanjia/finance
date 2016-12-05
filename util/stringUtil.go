package util

import (
	"strconv"
)

func ToInt(str string) int64 {
	value, _ := strconv.ParseInt(str, 10, 64)
	return value
}

func ToFloat(str string) float64 {
	value, _ := strconv.ParseFloat(str, 64)
	return value
}
