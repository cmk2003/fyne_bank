package page

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"strconv"
)

// CalculateLoanAmount 根据信用等级计算可贷款金额
func CalculateLoanAmount(creditRating int) float64 {
	switch {
	case creditRating >= 1 && creditRating <= 3:
		return 10000.0
	case creditRating >= 4 && creditRating <= 6:
		return 50000.0
	case creditRating >= 7 && creditRating <= 8:
		return 200000.0
	case creditRating == 9:
		return 300000.0
	case creditRating == 10:
		return 500000.0
	default:
		return 0.0 // 不合法的信用等级
	}
}
func MakeLoanUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {

	info := cardService.GetAccountInfo(userInfo.ID)
	fmt.Println(info)
	// 信誉等级决定可以贷款的多少
	creditLevelLabel := widget.NewLabel("信誉等级:" + strconv.Itoa(info.CreditRating))
	loanLabel := widget.NewLabel("可贷款金额：" + strconv.FormatFloat(CalculateLoanAmount(info.CreditRating), 'f', 2, 64) + "元")

	// 卡类别选择 根据用户id选择已经开户的卡
	cardTypeList := cardService.GetCardType(userInfo.ID)

	//卡的类型
	cardTypeOptions := make([]string, 0)
	//卡号
	cardNumberOptions := make([]string, 0)
	cardNumberSelect := widget.NewSelect(cardNumberOptions, func(value string) {
		// 在回调拿到卡号
		fmt.Println("Selected Card Number:", value)
		//根据卡号获取到余额
		//account := cardService.GetAccountByCardNumber(value)
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
	// 用户密码输入
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	submitButton := widget.NewButton("提交", func() {
		cardType := cardTypeSelect.Selected
		cardNumber := cardNumberSelect.Selected
		password := passwordEntry.Text

		if cardType == "" || cardNumber == "" || password == "" {
			dialog.ShowInformation("错误", "所有字段均为必填项", w)
			return
		}

	})
	overdraftRecordLink := widget.NewHyperlink("贷款记录", nil)
	overdraftRecordLink.OnTapped = func() {
		//卡号 TODO 贷款记录
		updatedAccount := cardService.GetAccountByCardNumber(cardNumberSelect.Selected)
		LoadRecords := OverdraftService.GetOverDraftByAccountId(updatedAccount.ID)

		list := widget.NewList(
			func() int {
				return len(LoadRecords)
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("template")
			},
			func(i widget.ListItemID, item fyne.CanvasObject) {
				record := LoadRecords[i]
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

	// 布局
	top := container.NewHBox(
		creditLevelLabel,
		loanLabel,
	)
	form := container.NewVBox(
		top,
		widget.NewLabel("选择卡类别"),
		cardTypeSelect,
		widget.NewLabel("选择卡号"),
		cardNumberSelect,
		widget.NewLabel("输入密码"),
		passwordEntry,
		submitButton,
	)
	return form
}
