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

// 获取用户列表
func (u *UserService) GetUserList() []model.User {
	var users []model.User
	global.DB.Find(&users)
	return users
}

// 新增用户
func (u *UserService) AddUser(user model.User) error {
	tx := global.DB.Create(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 删除用户
func (u *UserService) DeleteUser(id uint) interface{} {
	var user model.User
	tx := global.DB.Where("id = ?", id).Delete(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 修改用户
func (u *UserService) UpdateUser(user model.User) error {
	tx := global.DB.Model(&model.User{}).Where("username = ?", user.Username).
		Updates(map[string]interface{}{"role": user.Role, "gender": user.Gender})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 根据用户名搜索
func (u *UserService) SearchUser(username string) []model.User {
	var users []model.User
	global.DB.Where("username like ?", "%"+username+"%").Find(&users)
	return users
}

func (u *UserService) GetUserNameById(id uint) string {
	var user model.User
	global.DB.Where("id = ?", id).First(&user)
	return user.Username
}
