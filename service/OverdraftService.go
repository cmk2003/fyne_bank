package service

import (
	"sql_bank/global"
	"sql_bank/model"
	"time"
)

type OverdraftService struct {
}

func (s OverdraftService) AddOverdraft(id uint, amount float64) error {
	var overDraft model.Overdraft
	overDraft.AccountID = id
	overDraft.Amount = amount
	//当前日期加一个月
	overDraft.RepaymentDueDate = time.Now().AddDate(0, 1, 0)
	tx := global.DB.Create(&overDraft)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s OverdraftService) GetOverDraftByAccountId(id uint) []model.Overdraft {
	var overDraft []model.Overdraft
	tx := global.DB.Where("account_id = ?", id).Find(&overDraft)
	if tx.Error != nil {
		return []model.Overdraft{}
	}
	return overDraft
}

// 获取所有没还的透支记录
func (s OverdraftService) GetAllNoPaidOverDraft() []model.Overdraft {
	var overDraft []model.Overdraft
	tx := global.DB.Where("repaid = ?", false).Find(&overDraft)

	for _, i2 := range overDraft {
		first := global.DB.Preload("Account").First(&i2)
		if first.Error != nil {
			return []model.Overdraft{}
		}
	}

	if tx.Error != nil {
		return []model.Overdraft{}
	}
	return overDraft
}

func (s OverdraftService) UpdateOverDraft(draft model.Overdraft) {
	global.DB.Where("id = ?", draft.ID).Save(&draft)
}
