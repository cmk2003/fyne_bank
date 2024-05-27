package task

import (
	"fmt"
	"sql_bank/model"
	"sql_bank/service"
	"strconv"
	"time"
)

var (
	overDraftService = service.OverdraftService{}
	accountService   = service.AccountService{}
	//tranService      = service.TransferService{}
	cardService = service.CardService{}
)

// 定时清除透支金额
func RepayOverDraftsTask() {
	fmt.Println("RepayOverDraftsTask定时任务执行")
	//读取overdraft表
	drafts := overDraftService.GetAllNoPaidOverDraft()
	//遍历透支记录
	for _, draft := range drafts {
		if draft.Repaid == true {
			continue
		}
		//如果逾期
		fmt.Println(draft.RepaymentDueDate)
		fmt.Println(time.Now())
		if !draft.RepaymentDueDate.After(time.Now()) {
			fmt.Println("逾期")
			//标记为逾期
			draft.Repaid = false
			//更新透支记录
			overDraftService.UpdateOverDraft(draft)
			//把账户封禁
			accountService.BanAccount(draft.AccountID)
		} else {
			//如果余额大于透支金额
			fmt.Println(draft.Account)
			if draft.Account.Balance > draft.Amount {
				fmt.Println("如果余额大于透支金额")
				//还透支 记录交易
				transaction := model.Transaction{
					AccountID:       draft.AccountID,
					Amount:          draft.Amount,
					TransactionType: 7,
					TransactionDate: time.Now(),
					CardNumber:      draft.Account.AccountNumber,
					ToCardNumber:    "00000000",
					Status:          "fail",
				}
				mountStr := strconv.FormatFloat(draft.Amount, 'f', 2, 64)
				err := cardService.Transfer(draft.Account.AccountNumber, "00000000", mountStr, transaction)
				if err != nil {
					panic(err)
				}
				//标记为已还款
				draft.Repaid = true
				//更新透支记录
				overDraftService.UpdateOverDraft(draft)
				//透支额度加上去
				cardInfo := cardService.GetAccountByCardNumber(draft.Account.AccountNumber)
				//还款时间在规定时间之前
				if draft.RepaymentDueDate.After(time.Now()) {
					cardInfo.OverdraftLimit += draft.Amount * 1.1
				} else {
					cardInfo.OverdraftLimit += draft.Amount * 0.9
				}
				cardService.UpdateAccount(cardInfo)
				//解封账户
				accountService.UnBanAccount(draft.AccountID)
			}

		}

	}
}
