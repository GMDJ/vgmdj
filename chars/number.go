package chars

import (
	"github.com/vgmdj/utils/logger"
	"math"
	"strconv"
)

func ToInt(num interface{}) int {
	switch num.(type) {
	default:
		logger.Error("invalid type ", num)
		return 0

	case string:
		result, _ := strconv.Atoi(num.(string))
		return result

	case int:
		return num.(int)

	case int32:
		return int(num.(int32))

	case int64:
		return int(num.(int64))

	case float64:
		return int(num.(float64))

	case bool:
		if num.(bool) {
			return 1
		}
		return 0

	case []byte:
		result := 0
		for k, v := range num.([]byte) {
			result += (int(v) - 48) * int(math.Pow10(len(num.([]byte))-k-1))
		}
		return result

	}
}

func ToFloat64(num interface{}) float64 {
	switch num.(type) {
	default:
		logger.Error("invalid type ", num)
		return 0

	case float64:
		return num.(float64)

	case string:
		result, _ := strconv.ParseFloat(num.(string), 64)
		return result

	case int:
		return float64(num.(int))

	case int32:
		return float64(num.(int32))

	case int64:
		return float64(num.(int64))
	}
}

func ToString(num interface{}, prec ...int) string {
	var p int = 4
	if len(prec) != 0 {
		p = prec[0]
	}

	switch num.(type) {
	default:
		logger.Error("invalid type ", num)
		return ""

	case string:
		return num.(string)

	case int:
		return strconv.Itoa(num.(int))

	case int32:
		return strconv.Itoa(int(num.(int32)))

	case int64:
		return strconv.FormatInt(num.(int64), 10)

	case float64:
		return strconv.FormatFloat(num.(float64), 'f', p, 64)

	}
}

func TakeLeftInt(num int, n int) int {
	for num >= int(math.Pow10(n)) {
		num /= 10
	}

	return num
}

//TODO finish this function
func TekeRightInt(num int, n int) int {

	return 0
}
