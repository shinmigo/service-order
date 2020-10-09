package utils

import (
	"strconv"
	"time"
)

func GetUniqueId() uint64 {
	float, _ := strconv.ParseFloat(time.Now().Format("060102150405.00000"), 64)

	return uint64(float * 1000)
}
