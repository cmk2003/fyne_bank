package initialize

import (
	"fmt"
	"sql_bank/model"
	"sql_bank/service"
	"time"
)

var (
	transferService = service.TransferService{}
	//accountService  = service.AccountService{}
	cardService = service.CardService{}
)

func RecoverData() {
	//读取未转账完成的记录
	transactions, err := transferService.GetAllTransfers()
	if err != nil {
		panic(err)
	}
	//获取状态为fail的记录
	for _, transaction := range transactions {
		if transaction.Status == "fail" {
			//根据记录进行处理 加钱
			//记录交易
			backTransaction := model.Transaction{
				CardNumber:      transaction.ToCardNumber,
				ToCardNumber:    transaction.CardNumber,
				Amount:          transaction.Amount,
				Status:          "success",
				TransactionDate: time.Now(),
				TransactionType: 8,
			}
			err := cardService.Transfer(transaction.ToCardNumber, transaction.CardNumber, fmt.Sprintf("%f", transaction.Amount), backTransaction)
			if err != nil {
				panic(err)
			}
			fmt.Println(backTransaction)
			err = transferService.AddTransfer(backTransaction)
			if err != nil {
				panic(err)
			}
			//处理成功后更新状态
			transaction.Status = "backed"
			err = transferService.UpdateTransfer(transaction)
			if err != nil {
				panic(err)
			}
		}
	}
}
