package service

import (
	"errors"
	"sql_bank/global"
	"sql_bank/model"
	"sql_bank/utils"
)

var (
	accountService = AccountService{}
)

type CardTypeService struct {
}

// GetCardType 获取卡类型
func (c *CardTypeService) GetCardType() []model.AccountType {
	var cardType []model.AccountType
	tx := global.DB.Find(&cardType)
	if tx.Error != nil {
		return []model.AccountType{}
	}
	return cardType
}

// 新增卡类型
func (c *CardTypeService) AddCardType(cardType model.AccountType) (model.AccountType, error) {
	tx := global.DB.Create(&cardType)
	if tx.Error != nil {
		return model.AccountType{}, tx.Error
	}
	return cardType, nil
}

// 修改卡类型
func (c *CardTypeService) UpdateCardType(cardType model.AccountType) error {
	tx := global.DB.Model(&model.AccountType{}).
		Where("id = ?", cardType.ID).Updates(&cardType)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (c *CardTypeService) DeleteCardType(cardType model.AccountType) error {
	//判断当前卡类型是否被使用
	if ids := accountService.GetCardTypeId(); utils.ContainsInt(ids, int(cardType.ID)) {
		//不能删除
		return errors.New("当前卡类型正在被使用，不能删除")
	}
	tx := global.DB.Delete(&cardType)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// 搜索
func (c *CardTypeService) GetCardTypeByName(text string) []model.AccountType {
	var cardType []model.AccountType
	tx := global.DB.Where("name like ?", "%"+text+"%").Find(&cardType)
	if tx.Error != nil {
		return []model.AccountType{}
	}
	return cardType
}

func (c *CardTypeService) GetCardTypeRate(id uint) float64 {
	var cardType model.AccountType
	tx := global.DB.Where("id = ?", id).First(&cardType)
	if tx.Error != nil {
		return 0
	}
	return cardType.InterestRate
}

func (c *CardTypeService) GetBankNameByCardTypeId(id uint) string {
	var cardType model.AccountType
	tx := global.DB.Preload("Bank").Where("id = ?", id).First(&cardType)
	if tx.Error != nil {
		return ""
	}
	return cardType.Bank.Name
}
