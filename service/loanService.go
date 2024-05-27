package service

import (
	"sql_bank/global"
	"sql_bank/model"
)

type LoanService struct {
}

func (s LoanService) AddLoan(record model.Loan) error {
	tx := global.DB.Create(&record)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s LoanService) GetLoanList(ids []uint) []model.Loan {
	var loans []model.Loan
	tx := global.DB.Where("account_id in (?)", ids).Find(&loans)
	if tx.Error != nil {
		return []model.Loan{}
	}
	//获取卡号
	//tx = global.DB.Preload("Account").Find(&loans)
	for a, i := range loans {
		tx = global.DB.Preload("Account").First(&i)
		loans[a] = i
	}
	if tx.Error != nil {
		return []model.Loan{}
	}
	return loans
}

func (s LoanService) UpdateLoan(loan model.Loan) {
	global.DB.Save(&loan)
}

func (s LoanService) GetAllLoans() []model.Loan {
	var loans []model.Loan
	tx := global.DB.Preload("Account.User").Preload("Account.AccountType").Find(&loans)
	if tx.Error != nil {
		return []model.Loan{}
	}
	return loans

}

func (s LoanService) SearchLoans(s2 string) []model.Loan {
	// 模糊查询 多字段 卡号 真实姓名 账户类型
	var loans []model.Loan
	//tx := global.DB.Joins("Account").Joins("User").Joins("AccountType").
	//	Where("Account.account_number LIKE ? or User.real_name LIKE ? or AccountType.name LIKE ?", "%"+s2+"%", "%"+s2+"%", "%"+s2+"%").
	//	Find(&loans)
	tx := global.DB.Joins("JOIN account ON account.id = loan.account_id").
		Joins("JOIN user ON user.id = account.user_id").
		Joins("JOIN account_type ON account_type.id = account.account_type_id").
		Where("account.account_number LIKE ? OR user.real_name LIKE ? OR account_type.name LIKE ?", "%"+s2+"%", "%"+s2+"%", "%"+s2+"%").
		Find(&loans)

	//preload
	for a, i := range loans {
		tx = global.DB.Preload("Account").First(&i)
		loans[a] = i
	}
	// preload User
	for a, i := range loans {
		tx = global.DB.Preload("Account.User").First(&i)
		loans[a] = i
	}
	if tx.Error != nil {
		return []model.Loan{}
	}
	return loans
}

func (s2 LoanService) SearchLoansByPage(s string, page int, size int) []model.Loan {
	var loans []model.Loan
	query := global.DB.Joins("JOIN account ON account.id = loan.account_id").
		Joins("JOIN user ON user.id = account.user_id").
		Joins("JOIN account_type ON account_type.id = account.account_type_id").
		Where("account.account_number LIKE ? OR user.real_name LIKE ? OR account_type.name LIKE ?", "%"+s+"%", "%"+s+"%", "%"+s+"%").
		Offset((page - 1) * size).Limit(size)

	// 预加载 Account, Account.User 和 Account.AccountType
	tx := query.Preload("Account.User").Preload("Account.AccountType").Find(&loans)

	if tx.Error != nil {
		return []model.Loan{}
	}
	return loans

}

func (s2 LoanService) GetNoPayLoan(s string, page int, size int) []model.Loan {
	var loans []model.Loan
	query := global.DB.Joins("JOIN account ON account.id = loan.account_id").
		Joins("JOIN user ON user.id = account.user_id").
		Joins("JOIN account_type ON account_type.id = account.account_type_id").
		Where("account.account_number LIKE ? OR user.real_name LIKE ? OR account_type.name LIKE ?", "%"+s+"%", "%"+s+"%", "%"+s+"%").
		Where("loan.status = ?", 0).
		Offset((page - 1) * size).Limit(size)

	// 预加载 Account, Account.User 和 Account.AccountType
	tx := query.Preload("Account.User").Preload("Account.AccountType").Find(&loans)

	if tx.Error != nil {
		return []model.Loan{}
	}
	return loans

}
