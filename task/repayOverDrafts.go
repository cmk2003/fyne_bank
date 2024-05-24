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
		if draft.RepaymentDueDate.After(time.Now()) {
			//标记为逾期
			draft.Repaid = false
			//更新透支记录
			overDraftService.UpdateOverDraft(draft)
			//把账户封禁
			accountService.BanAccount(draft.AccountID)
		} else {
			//如果余额大于透支金额
			if draft.Account.Balance > draft.Amount {
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
			}

		}

	}
}
