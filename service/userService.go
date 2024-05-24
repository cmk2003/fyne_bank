package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sql_bank/global"
	"sql_bank/model"
	"sql_bank/utils"
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
	//读取数据库的salt
	_, salt, hash, err := utils.ParseDjangoHash(user.Password)
	if err != nil {
		fmt.Println("Error parsing Django hash:", err)
		return false, model.User{}
	}
	//使用传来的密码做加密
	newHash := utils.EncryptPassword(password, salt)
	fmt.Println(newHash, hash)
	// 检查密码
	if newHash == user.Password {
		//user.Password = ""
		return true, user
	} else {
		//user.Password = ""
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
		Updates(map[string]interface{}{"role": user.Role, "gender": user.Gender, "phone": user.Phone, "email": user.Email, "address": user.Address, "real_name": user.RealName})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 根据用户名搜索
func (u *UserService) SearchUser(username string) []model.User {
	// 模糊查询 多字段搜索 用户真实名字 电话号码 邮箱
	var users []model.User
	searchPattern := "%" + username + "%"
	global.DB.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ? OR email LIKE ?", searchPattern, searchPattern, searchPattern, searchPattern).Find(&users)
	return users
}

func (u *UserService) GetUserNameById(id uint) string {
	var user model.User
	global.DB.Where("id = ?", id).First(&user)
	return user.RealName
}

func (u *UserService) FreezeUser(id uint) {
	var user model.User
	global.DB.Where("id = ?", id).First(&user)
	user.IsFrozen = true
	global.DB.Save(&user)
}

// 解冻
func (u *UserService) UnFreezeUser(id uint) {
	var user model.User
	global.DB.Where("id = ?", id).First(&user)
	user.IsFrozen = false
	global.DB.Save(&user)
}

func (u *UserService) ChangePass(info model.User) interface{} {
	tx := global.DB.Model(&model.User{}).Where("username = ?", info.Username).
		Updates(map[string]interface{}{"password": info.Password})
	if tx.Error != nil {
		return tx.Error
	}
	return nil

}
