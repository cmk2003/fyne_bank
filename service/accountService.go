package service

import (
	"sql_bank/global"
	"sql_bank/model"
)

type AccountService struct {
}

// 获取账户中的卡类型id
func (a *AccountService) GetCardTypeId() []uint {
	var account []model.Account
	tx := global.DB.Find(&account)
	if tx.Error != nil {
		return []uint{}
	}
	ids := make([]uint, 0)
	for i := 0; i < len(account); i++ {
		ids = append(ids, account[i].AccountTypeID)
	}
	return ids
}

func (a *AccountService) GetAccountList() []model.Account {
	var accountList []model.Account
	tx := global.DB.Preload("AccountType").Find(&accountList)
	if tx.Error != nil {
		return []model.Account{}
	}
	return accountList
}

func (a *AccountService) AddAccount(account model.Account) (model.Account, error) {
	tx := global.DB.Create(&account)
	if tx.Error != nil {
		return model.Account{}, tx.Error
	}
	return account, nil
}

func (a *AccountService) UpdateAccount(account model.Account) interface{} {
	tx := global.DB.Model(&model.Account{}).Where("id = ?", account.ID).Updates(&account)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *AccountService) DeleteAccount(account model.Account) interface{} {
	tx := global.DB.Delete(&account)
	if tx.Error != nil {
		return tx.Error
	}
	return nil

}

// 账号 TODO 多字段搜索
func (a *AccountService) SearchAccount(text string) []model.Account {
	var accountList []model.Account
	tx := global.DB.Where("account_number LIKE ?", "%"+text+"%").Find(&accountList)
	if tx.Error != nil {
		return []model.Account{}
	}
	return accountList
}

// 给账户充值
func (a *AccountService) Recharge(accountNumber string, amount float64) interface{} {
	var account model.Account
	tx := global.DB.Where("account_number = ?", accountNumber).First(&account)
	if tx.Error != nil {
		return tx.Error
	}
	account.Balance += amount
	tx = global.DB.Save(&account)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 账户取款
func (a *AccountService) Withdraw(accountNumber string, amount float64) interface{} {
	var account model.Account
	tx := global.DB.Where("account_number = ?", accountNumber).First(&account)
	if tx.Error != nil {
		return tx.Error
	}
	if account.Balance < amount {
		return "余额不足"
	}
	account.Balance -= amount
	tx = global.DB.Save(&account)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *AccountService) BanAccount(id uint) {
	var account model.Account
	tx := global.DB.Where("id = ?", id).First(&account)
	if tx.Error != nil {
		return
	}
	account.IsOverdraftLimitReached = true
	tx = global.DB.Save(&account)
	if tx.Error != nil {
		return
	}
}
