package service

import (
	"fmt"
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
	var overDrafts []model.Overdraft
	tx := global.DB.Preload("Account").Where("repaid = ?", false).Find(&overDrafts)

	if tx.Error != nil {
		fmt.Println("Query Error:", tx.Error)
		return []model.Overdraft{}
	}

	return overDrafts
}

func (s OverdraftService) UpdateOverDraft(draft model.Overdraft) {
	global.DB.Where("id = ?", draft.ID).Save(&draft)
}

func (s OverdraftService) GetOverdraftList(ids []uint) []model.Overdraft {
	var overdraftList []model.Overdraft
	tx := global.DB.Preload("Account").Where("account_id IN (?)", ids).Find(&overdraftList)

	if tx.Error != nil {
		return []model.Overdraft{}
	}

	return overdraftList
}

func (s OverdraftService) UpdateOverdraft(overdraft model.Overdraft) error {
	tx := global.DB.Save(&overdraft)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (s OverdraftService) GetOverdraftListWithNoPay(uints []uint) []model.Overdraft {
	var overdraftList []model.Overdraft
	tx := global.DB.Where("account_id in (?) and repaid = ?", uints, false).Find(&overdraftList)
	for i, i2 := range overdraftList {
		first := global.DB.Preload("Account").First(&i2)
		if first.Error != nil {
			return []model.Overdraft{}
		}
		overdraftList[i] = i2
	}
	if tx.Error != nil {
		return []model.Overdraft{}
	}
	return overdraftList
}
