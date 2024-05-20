package service

import (
	"errors"
	"sql_bank/global"
	"sql_bank/model"
)

type CardService struct {
}

type ResCardType struct {
	ID          uint   `gorm:"primary_key" json:"id"`                        // 账户类型ID
	Name        string `gorm:"type:varchar(20);not null;unique" json:"name"` // 账户类型名称，如"招商银行一卡通"
	Description string `gorm:"type:varchar(255);" json:"description"`        // 账户类型描述
}

func (c *CardService) GetCardType(id uint) []ResCardType {
	//查找当前用户已经开户的分类
	var account []model.Account
	tx := global.DB.Preload("AccountType").Where("user_id = ?", id).Find(&account)
	//属性拷贝
	var resCardType []ResCardType
	for _, v := range account {
		resCardType = append(resCardType, ResCardType{
			ID:          v.ID,
			Name:        v.AccountType.Name,
			Description: v.AccountType.Description,
		})
	}
	if tx.Error != nil {
		return []ResCardType{}
	}
	return resCardType
}

func (c *CardService) Saving(number string, password string, amount float64, id uint) error {
	//判断卡号和类别是否对应
	var account model.Account
	tx := global.DB.Where("account_number = ? and account_type_id = ?", number, id).First(&account)
	if tx.Error != nil {
		//自定义error
		return errors.New("卡号和类别不匹配")
	}
	//验证number和密码 TODO 密码进行hash加密
	if account.PasswordHash != password {
		return errors.New("密码错误")
	}
	//存款
	account.Balance += amount
	tx = global.DB.Save(&account)
	if tx.Error != nil {
		return errors.New("存款失败")
	}
	return nil

}

func (c *CardService) GetAccountInfo(id uint) model.Account {
	//获取账户信息
	var account model.Account
	tx := global.DB.Where("user_id = ?", id).First(&account)
	if tx.Error != nil {
		return model.Account{}
	}
	account.PasswordHash = ""
	return account

}

func (c *CardService) GetCardNumber(id uint, id2 uint) []model.Account {
	//根据用户id和typeid获取卡号
	var account []model.Account
	tx := global.DB.Where("user_id = ? and account_type_id = ?", id, id2).Find(&account)
	if tx.Error != nil {
		return []model.Account{}
	}
	return account
}