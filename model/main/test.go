package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sql_bank/model"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/sqlBank?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	_ = db.AutoMigrate(&model.User{},
		&model.AccountType{},
		&model.Account{},
		&model.Transaction{},
		&model.Loan{},
		&model.Overdraft{},
	)
}
