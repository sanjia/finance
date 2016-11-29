package stock
import (
	"strconv"
)

func F(s string) float64 {
	f,_:=strconv.ParseFloat(s, 32)
	return f
}
