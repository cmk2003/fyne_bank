package page

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/service"
	"strconv"
)

var (
	OverdraftService = service.OverdraftService{}
)

func MakeWithdrawUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {
	//获取用户账户信息
	info := cardService.GetAccountInfo(userInfo.ID)
	fmt.Println(info)
	// 用户信息标签
	balanceLabel := widget.NewLabel("当前余额: N/A")
	OverdraftLimitLabel := widget.NewLabel("当前可透支额度: N/A")
	creditLevelLabel := widget.NewLabel("信誉等级:" + strconv.Itoa(info.CreditRating))
	//overdraftRecordLabel := widget.NewLabel("透支记录: N/A")

	// 卡类别选择 根据用户id选择已经开户的卡
	cardTypeList := cardService.GetCardType(userInfo.ID)

	fmt.Println("cardTypeList", cardTypeList)
	//卡的类型
	cardTypeOptions := make([]string, 0)
	//卡号
	cardNumberOptions := make([]string, 0)
	cardNumberSelect := widget.NewSelect(cardNumberOptions, func(value string) {
		// 在回调拿到卡号
		fmt.Println("Selected Card Number:", value)
		//根据卡号获取到余额
		account := cardService.GetAccountByCardNumber(value)
		balanceLabel.SetText("当前余额: " + strconv.FormatFloat(account.Balance, 'f', 2, 64))
		OverdraftLimitLabel.SetText("当前可透支额度: " + strconv.FormatFloat(account.OverdraftLimit, 'f', 2, 64))
	})

	for _, cardType := range cardTypeList {
		cardTypeOptions = append(cardTypeOptions, cardType.Name)
	}
	fmt.Println(cardTypeOptions)
	// 定义一个变量来保存所选卡类别的ID
	var selectedCardTypeID uint
	cardTypeSelect := widget.NewSelect(cardTypeOptions, func(value string) {
		// 在回调中查找所选项的ID
		for _, cardType := range cardTypeList {
			if cardType.Name == value {
				cardNumberSelect.Options = []string{}
				cardNumberOptions = []string{}
				cardNumberSelect.Selected = ""
				cardNumberSelect.Refresh() // 刷新选择器
				selectedCardTypeID = cardType.ID
				fmt.Println("Selected Card Type ID:", selectedCardTypeID)
				//获取根据卡类别选择卡号userid和typeid 获取countList
				cardInfoList := cardService.GetCardNumber(userInfo.ID, selectedCardTypeID)
				for _, cardInfo := range cardInfoList {
					cardNumberOptions = append(cardNumberOptions, cardInfo.AccountNumber)
				}
				fmt.Println(cardNumberOptions)
				cardNumberSelect.Options = cardNumberOptions // 更新选项
				cardNumberSelect.Refresh()                   // 刷新选择器
				break
			}
		}
	})
	overdraftRecordLink := widget.NewHyperlink("透支记录", nil)
	overdraftRecordLink.OnTapped = func() {
		// Fetch the updated account information
		updatedAccount := cardService.GetAccountByCardNumber(cardNumberSelect.Selected)
		// Fetch and format the overdraft records
		overdraftRecords := OverdraftService.GetOverDraftByAccountId(updatedAccount.ID)
		// Create and show the dialog_recrd

		list := widget.NewList(
			func() int {
				return len(overdraftRecords)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, item fyne.CanvasObject) {
				record := overdraftRecords[i]
				item.(*widget.Label).SetText(
					fmt.Sprintf("Record ID: %d, Amount: %.2f, Due Date: %s, Repaid: %t",
						record.ID, record.Amount, record.RepaymentDueDate.Format("2006-01-02"), record.Repaid))
			},
		)

		// Create a custom dialog_recrd to display the list
		dialog_recrd := dialog.NewCustom("透支记录", "Close", container.NewVScroll(list), w)
		dialog_recrd.Resize(fyne.NewSize(500, 300)) // Adjust the size as needed
		dialog_recrd.Show()

	}

	// 用户密码输入
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// 取款金额输入
	depositAmountEntry := widget.NewEntry()
	depositAmountEntry.SetPlaceHolder("请输入取款金额")

	// 提交按钮
	submitButton := widget.NewButton("提交", func() {
		cardType := cardTypeSelect.Selected
		cardNumber := cardNumberSelect.Selected
		password := passwordEntry.Text
		depositAmount := depositAmountEntry.Text

		if cardType == "" || cardNumber == "" || password == "" || depositAmount == "" {
			dialog.ShowInformation("错误", "所有字段均为必填项", w)
			return
		}

		amount, err := strconv.ParseFloat(depositAmount, 64)
		if err != nil || amount <= 0 {
			dialog.ShowInformation("错误", "取款金额必须是正数", w)
			return
		}
		// 取款如果超过余额则提示是否要透支，透支提示余额，超过透支余额就不让取款，不超过就生成透支表，下次透支的时候需要判断是否有透支行为
		// 获取当前账号信息
		if account := cardService.GetAccountByCardNumber(cardNumber); account.IsOverdraftLimitReached == true {
			dialog.ShowInformation("失败", "超时未还款，禁止透支", w)
			return
		} else {
			//不存在透支行为可以进行透支
			if account.Balance < amount {
				dialog.NewConfirm("提示", "余额不足，是否要透支？", func(b bool) {
					if b {
						if account.Balance+account.OverdraftLimit < amount {
							dialog.ShowInformation("失败", "透支额度不足", w)
							return
						} else {
							// 透支 余额为0 透支额度减去取款金额 透支记录表增加一条记录
							amount = amount - account.Balance
							err := cardService.WithDrawOverDraft(cardNumber, password, amount, selectedCardTypeID)
							_ = OverdraftService.AddOverdraft(account.ID, amount)
							if err != nil {
								dialog.ShowInformation("失败", err.Error(), w)
							} else {
								dialog.ShowInformation("成功", "透支成功", w)
								updatedAccount := cardService.GetAccountByCardNumber(cardNumber)
								// Update the labels with the new account information
								balanceLabel.SetText("当前余额: " + strconv.FormatFloat(updatedAccount.Balance, 'f', 2, 64))
								OverdraftLimitLabel.SetText("当前可透支额度: " + strconv.FormatFloat(updatedAccount.OverdraftLimit, 'f', 2, 64))
								creditLevelLabel.SetText("信誉等级:" + string(rune(updatedAccount.CreditRating)))
								// You need to implement a method to fetch the overdraft record
								//overdraftRecord := OverdraftService.GetOverdraftRecord(updatedAccount.ID)
								//overdraftRecordLabel.SetText("透支记录: " + overdraftRecord)

							}
						}
					}
				}, w).Show()
			} else {
				err = cardService.WithDraw(cardNumber, password, amount, selectedCardTypeID)
				if err != nil {
					dialog.ShowInformation("失败", err.Error(), w)
				} else {
					dialog.ShowInformation("成功", "取款成功", w)
					updatedAccount := cardService.GetAccountByCardNumber(cardNumber)
					// Update the labels with the new account information
					balanceLabel.SetText("当前余额: " + strconv.FormatFloat(updatedAccount.Balance, 'f', 2, 64))
					OverdraftLimitLabel.SetText("当前可透支额度: " + strconv.FormatFloat(updatedAccount.OverdraftLimit, 'f', 2, 64))
					creditLevelLabel.SetText("信誉等级:" + string(rune(updatedAccount.CreditRating)))
					// You need to implement a method to fetch the overdraft record
					//overdraftRecord := OverdraftService.GetOverdraftRecord(updatedAccount.ID)
					//overdraftRecordLabel.SetText("透支记录: " + overdraftRecord)
				}
			}
		}

	})
	top := container.NewHBox(
		balanceLabel,
		creditLevelLabel,
		OverdraftLimitLabel,
		overdraftRecordLink,
	)
	// 布局
	form := container.NewVBox(
		top,
		widget.NewLabel("选择卡类别"),
		cardTypeSelect,
		widget.NewLabel("选择卡号"),
		cardNumberSelect,
		widget.NewLabel("输入密码"),
		passwordEntry,
		widget.NewLabel("输入取款金额"),
		depositAmountEntry,
		submitButton,
	)

	return form
}
