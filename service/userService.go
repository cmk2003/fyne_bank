package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sql_bank/global"
	"sql_bank/model"
)

type UserService struct {
}

func (u *UserService) LoginSys(username string, password string) (bool, model.User) {
	var user model.User
	// 查找用户
	result := global.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, model.User{}
		}
		fmt.Println("Error querying the database: ", result.Error)
		return false, model.User{}
	}
	// 检查密码
	if user.Password == password {
		user.Password = ""
		return true, user
	} else {
		user.Password = ""
		return false, model.User{}
	}
}
