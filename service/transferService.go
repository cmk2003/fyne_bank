package service

import (
	"sql_bank/global"
	"sql_bank/model"
)

type TransferService struct {
}

// 新增记录
func (t *TransferService) AddTransfer(record model.Transaction) error {
	tx := global.DB.Create(&record)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 读取所有记录
func (t *TransferService) GetAllTransfers() ([]model.Transaction, error) {
	var records []model.Transaction
	tx := global.DB.Find(&records)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return records, nil
}

// 更新记录
func (t *TransferService) UpdateTransfer(transaction model.Transaction) error {
	tx := global.DB.Model(&model.Transaction{}).Where("id = ?", transaction.ID).Updates(&transaction)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
