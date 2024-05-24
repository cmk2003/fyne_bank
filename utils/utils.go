package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"sql_bank/global"
	"strings"
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
func CheckRedisKeyExist(key string) bool {
	exists, err := global.RDB.Exists(global.Ctx, key).Result()
	if err != nil {
		log.Printf("Error checking if key exists: %v\n", err)
		return false
	}
	if exists == 0 {
		fmt.Println("Key does not exist:", key)
		//不提醒直接放行
		return false
	}
	global.RDB.Del(global.Ctx, key)
	return true
}
func GenerateRandomSalt(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// MD5Hash generates the MD5 hash of the input string.
func MD5Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// DjangoHash generates a Django-compatible hashed password.
func DjangoHash(password string) string {
	salt := GenerateRandomSalt(12)
	hash := MD5Hash(salt + password)
	return fmt.Sprintf("md5$%s$%s", salt, hash)
}
func EncryptPassword(password string, salt string) string {
	//salt := GenerateRandomSalt(12)
	hash := MD5Hash(salt + password)
	return fmt.Sprintf("md5$%s$%s", salt, hash)
}
func ParseDjangoHash(djangoHash string) (algorithm, salt, hash string, err error) {
	parts := strings.Split(djangoHash, "$")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid Django hash format")
	}
	algorithm = parts[0]
	salt = parts[1]
	hash = parts[2]
	return algorithm, salt, hash, nil
}
