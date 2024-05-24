package main

import (
	"fmt"
	"sql_bank/utils"
)

func main() {
	hash := utils.DjangoHash("123456")
	println(hash)
	algorithm, salt, hash, err := utils.ParseDjangoHash(hash)
	if err != nil {
		fmt.Println("Error parsing Django hash:", err)
		return
	}

	fmt.Println("Algorithm:", algorithm)
	fmt.Println("Salt:", salt)
	fmt.Println("Hash:", hash)
}
