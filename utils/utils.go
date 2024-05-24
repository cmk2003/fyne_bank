package utils

import (
	"math/rand"
	"time"
)

func ContainsInt(slice []uint, value int) bool {
	for _, v := range slice {
		if v == uint(value) { // Convert the int value to uint before comparing
			return true
		}
	}
	return false
}
func GenerateRandomAccountID(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
