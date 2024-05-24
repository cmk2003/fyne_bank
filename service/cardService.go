package service

import (
	"errors"
	"sql_bank/global"
	"sql_bank/model"
	"strconv"
)

type CardService struct {
}

type ResCardType struct {
	ID          uint   `gorm:"primary_key" json:"id"`                        // 账户类型ID
	Name        string `gorm:"type:varchar(20);not null;unique" json:"name"` // 账户类型名称，如"招商银行一卡通"
	Description string `gorm:"type:varchar(255);" json:"description"`        // 账户类型描述
}

var (
	transferService = TransferService{}
)

func (c *CardService) GetCardType(id uint) []ResCardType {
	//查找当前用户已经开户的分类
	var account []model.Account
	tx := global.DB.Preload("AccountType").Where("user_id = ?", id).Find(&account)
	//属性拷贝
	var resCardType []ResCardType
	for _, v := range account {
		resCardType = append(resCardType, ResCardType{
			ID:          v.AccountTypeID,
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

func (c *CardService) GetAccountByCardNumber(value string) model.Account {
	//根据卡号获取账户信息
	var account model.Account
	tx := global.DB.Where("account_number = ? ", value).First(&account)
	if tx.Error != nil {
		return model.Account{}
	}
	return account

}

func (c *CardService) WithDrawOverDraft(s string, s2 string, amount float64, id uint) error {
	//根据cardnumber获取账户信息
	var account model.Account
	_ = global.DB.Where("account_number = ? and account_type_id = ?", s, id).First(&account)
	if account.PasswordHash != s2 {
		return errors.New("密码错误")
	}
	//判断是否透支
	if account.IsOverdraftLimitReached {
		return errors.New("透支已经达到限额")
	}
	// 余额置为0
	account.Balance = 0
	account.OverdraftLimit -= amount
	// 未还款才需要置为true
	//account.IsOverdraftLimitReached = true

	global.DB.Save(&account)
	return nil
}

func (c *CardService) WithDraw(number string, password string, amount float64, id uint) error {
	//根据cardnumber获取账户信息
	var account model.Account
	_ = global.DB.Where("account_number = ? and account_type_id = ?", number, id).First(&account)
	if account.PasswordHash != password {
		return errors.New("密码错误")
	}
	if account.Balance < amount {
		return errors.New("余额不足")
	}
	account.Balance -= amount
	global.DB.Save(&account)
	return nil

}

func (c *CardService) VerifyPassword(number string, password string) bool {
	//验证密码
	var account model.Account
	_ = global.DB.Where("account_number = ?", number).First(&account)
	if account.PasswordHash == password {
		return true
	}
	return false
}

func (c *CardService) AddBalanceFromLoan(number string, float float64) error {
	//根据卡号增加余额
	var account model.Account
	_ = global.DB.Where("account_number = ?", number).First(&account)
	account.Balance += float
	tx := global.DB.Save(&account)
	if tx.Error != nil {
		return errors.New("存款失败")
	}
	return nil
}

// 转账
func (c *CardService) Transfer(selected string, text2 string, text3 string, transaction model.Transaction) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return errors.New("开启事务失败")
	}

	// 查询发起方账户
	var account model.Account
	if err := tx.Where("account_number = ?", selected).First(&account).Error; err != nil {
		tx.Rollback()
		return errors.New("查询账户失败")
	}

	// 转换金额
	val, err := strconv.ParseFloat(text3, 64)
	if err != nil {
		tx.Rollback()
		return errors.New("金额转换失败")
	}

	// 检查余额是否充足
	if account.Balance < val {
		tx.Rollback()
		return errors.New("余额不足")
	}

	// 扣减发起方账户余额
	account.Balance -= val
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return errors.New("更新账户余额失败")
	}

	// 查询接收方账户
	var account2 model.Account
	if err := tx.Where("account_number = ?", text2).First(&account2).Error; err != nil {
		tx.Rollback()
		return errors.New("查询对方账户失败")
	}

	// 增加接收方账户余额
	account2.Balance += val
	if err := tx.Save(&account2).Error; err != nil {
		tx.Rollback()
		return errors.New("更新对方账户余额失败")
	}
	// 插入一条初始状态的事务记录
	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("初始化事务失败")
	}

	// 记录交易
	transaction.Status = "success"

	// 使用事务保存交易记录
	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		return errors.New("更新事务状态失败")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return errors.New("提交事务失败")
	}

	return nil
}

func (c *CardService) UpdateAccount(info model.Account) {
	//更新账户信息
	global.DB.Save(&info)
}
