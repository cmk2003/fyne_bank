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
