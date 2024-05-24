package user

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/service"
	"strconv"
	"time"
)

var (
	cardTypeService = service.CardTypeService{}
	loanService     = service.LoanService{}
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
func getDepositRate(duration time.Duration) float64 {
	// 利率表（示例）
	// 3个月、6个月、1年、2年、3年、5年
	rates := map[string]float64{
		"3m": 1.463,
		"6m": 1.725,
		"1y": 2.05,
		"2y": 2.78,
		"3y": 3.538,
		"5y": 3.194,
	}

	switch {
	case duration <= 3*30*24*time.Hour:
		return rates["3m"]
	case duration <= 6*30*24*time.Hour:
		return rates["6m"]
	case duration <= 12*30*24*time.Hour:
		return rates["1y"]
	case duration <= 2*365*24*time.Hour:
		return rates["2y"]
	case duration <= 3*365*24*time.Hour:
		return rates["3y"]
	case duration <= 5*365*24*time.Hour:
		return rates["5y"]
	default:
		return rates["5y"]
	}
}
func MakeLoanUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {

	info := cardService.GetAccountInfo(userInfo.ID)
	fmt.Println(info)
	// 信誉等级决定可以贷款的多少
	creditLevelLabel := widget.NewLabel("信誉等级:N/A")
	loanLabel := widget.NewLabel("一次可贷款金额：N/A")

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
		account := cardService.GetAccountByCardNumber(value)
		//更新信誉等级
		//info = cardService.GetAccountInfo(userInfo.ID)
		creditLevelLabel.SetText("信誉等级:" + strconv.Itoa(account.CreditRating))
		loanLabel.SetText("一次可贷款金额：" + strconv.FormatFloat(CalculateLoanAmount(account.CreditRating), 'f', 2, 64))
	})
	for _, cardType := range cardTypeList {
		cardTypeOptions = append(cardTypeOptions, cardType.Name)
	}
	fmt.Println(cardTypeOptions)
	// 定义一个变量来保存所选卡类别的ID
	curRate := 0.0
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
				curRate = cardTypeService.GetCardTypeRate(selectedCardTypeID)
				break
			}
		}
	})
	// 用户密码输入
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	//overdraftRecordLink := widget.NewHyperlink("贷款记录", nil)
	//overdraftRecordLink.OnTapped = func() {
	//	//卡号
	//	updatedAccount := cardService.GetAccountByCardNumber(cardNumberSelect.Selected)
	//	LoadRecords := OverdraftService.GetOverDraftByAccountId(updatedAccount.ID)
	//
	//	list := widget.NewList(
	//		func() int {
	//			return len(LoadRecords)
	//		},
	//		func() fyne.CanvasObject {
	//			return widget.NewLabel("template")
	//		},
	//		func(i widget.ListItemID, item fyne.CanvasObject) {
	//			record := LoadRecords[i]
	//			item.(*widget.Label).SetText(
	//				fmt.Sprintf("Record ID: %d, Amount: %.2f, Due Date: %s, Repaid: %t",
	//					record.ID, record.Amount, record.RepaymentDueDate.Format("2006-01-02"), record.Repaid))
	//		},
	//	)
	//
	//	// Create a custom dialog_recrd to display the list
	//	dialog_recrd := dialog.NewCustom("透支记录", "Close", container.NewVScroll(list), w)
	//	dialog_recrd.Resize(fyne.NewSize(500, 300)) // Adjust the size as needed
	//	dialog_recrd.Show()
	//
	//}
	//时间选择器
	currentYear := time.Now().Year()
	years := []string{}
	for i := currentYear; i <= currentYear+5; i++ {
		years = append(years, fmt.Sprintf("%d", i))
	}
	yearSelect := widget.NewSelect(years, nil)
	yearSelect.PlaceHolder = "YYYY"

	months := []string{}
	for i := 1; i <= 12; i++ {
		months = append(months, fmt.Sprintf("%02d", i))
	}
	monthSelect := widget.NewSelect(months, nil)
	monthSelect.PlaceHolder = "MM"

	days := []string{}
	for i := 1; i <= 31; i++ {
		days = append(days, fmt.Sprintf("%02d", i))
	}
	daySelect := widget.NewSelect(days, nil)
	daySelect.PlaceHolder = "DD"

	selectedDate := widget.NewLabel("Selected Date: ----/--/--")
	selectedRate := widget.NewLabel("利率:" + strconv.FormatFloat(0.0, 'f', 2, 64))
	selectButton := widget.NewButton("Select Date", func() {
		year, err1 := strconv.Atoi(yearSelect.Selected)
		month, err2 := strconv.Atoi(monthSelect.Selected)
		day, err3 := strconv.Atoi(daySelect.Selected)
		if err1 == nil && err2 == nil && err3 == nil {
			selectedDateStr := fmt.Sprintf("%04d/%02d/%02d", year, month, day)
			selectedDate.SetText("Selected Date: " + selectedDateStr)
			selectedTime, _ := time.Parse("2006/01/02", selectedDateStr)
			fmt.Println(selectedDate)
			duration := selectedTime.Sub(time.Now())
			// 计算利率
			fmt.Println(duration, curRate)
			rate := getDepositRate(duration) * curRate
			selectedRate.SetText("利率:" + strconv.FormatFloat(rate, 'f', 2, 64))
		} else {
			selectedDate.SetText("Invalid date selected")
		}
	})

	datePicker := container.NewVBox(
		container.NewHBox(widget.NewLabel("Year:"),
			yearSelect, widget.NewLabel("Month:"), monthSelect,
			widget.NewLabel("Day:"), daySelect,
			selectButton,
			selectedDate,
			selectedRate,
		),
	)
	//贷款金额
	loanAmountEntry := widget.NewEntry()
	loanAmountEntry.SetPlaceHolder("请输入贷款金额")

	//提交按钮
	submitButton := widget.NewButton("提交", func() {
		cardType := cardTypeSelect.Selected
		cardNumber := cardNumberSelect.Selected
		password := passwordEntry.Text
		loanAmount := loanAmountEntry.Text

		year, _ := strconv.Atoi(yearSelect.Selected)
		month, _ := strconv.Atoi(monthSelect.Selected)
		day, _ := strconv.Atoi(daySelect.Selected)
		selectedDateStr := fmt.Sprintf("%04d/%02d/%02d", year, month, day)
		selectedTime, _ := time.Parse("2006/01/02", selectedDateStr)
		fmt.Println(selectedDate)
		duration := selectedTime.Sub(time.Now())

		if cardType == "" || cardNumber == "" || password == "" || loanAmount == "" {
			dialog.ShowInformation("错误", "所有字段均为必填项", w)
			return
		}
		// 贷款金额必须是正数
		loanAmountFloat, err := strconv.ParseFloat(loanAmount, 64)
		if err != nil || loanAmountFloat <= 0 {
			dialog.ShowInformation("错误", "贷款金额必须是正数", w)
			return
		}
		//贷款金额是否超过
		if loanAmountFloat > CalculateLoanAmount(info.CreditRating) {
			dialog.ShowInformation("错误", "贷款金额超过信用等级限制", w)
			return
		}
		//计算利息
		// 计算利息
		interest := loanAmountFloat * curRate / 100 * duration.Hours() / (365 * 24)

		//显示详细信息到一个表上
		bankName := cardTypeService.GetBankNameByCardTypeId(selectedCardTypeID)
		details := fmt.Sprintf("卡类型: %s\n卡号: %s\n贷款金额: %s\n选择的日期: %s\n当前利率 %s\n 本次贷款利息: %.2f\n 银行名称: %s\n",
			cardType, cardNumber, loanAmount, selectedDateStr, selectedRate.Text[6:], interest, bankName)

		detailLabel := widget.NewLabel(details)

		dialog_recrd := dialog.NewCustomConfirm("贷款详情", "确定", "取消", container.NewVScroll(detailLabel), func(confirm bool) {
			if confirm {
				// 确认操作
				fmt.Println("确认提交贷款信息")
				//密码是否正确
				if !cardService.VerifyPassword(cardNumber, password) {
					dialog.ShowInformation("错误", "密码错误", w)
					return
				}
				//生成贷款记录
				loanRecord := model.Loan{
					AccountID:       cardService.GetAccountByCardNumber(cardNumber).ID,
					AmountBorrowed:  loanAmountFloat,
					InterestRate:    curRate,
					LoanDate:        time.Now(),
					DueDate:         selectedTime,
					InterestAccrued: interest,
				}
				err2 := loanService.AddLoan(loanRecord)
				if err2 != nil {
					dialog.ShowInformation("失败", "贷款失败", w)
					return
				}
				//向余额里面加钱
				err2 = cardService.AddBalanceFromLoan(cardNumber, loanAmountFloat)
				if err2 != nil {
					dialog.ShowInformation("失败", "贷款失败", w)
					return
				}
				dialog.ShowInformation("成功", "贷款成功，注意查收账户余额", w)
			} else {
				// 取消操作
				fmt.Println("取消提交贷款信息")
			}
		}, w)
		dialog_recrd.Resize(fyne.NewSize(500, 300)) // Adjust the size as needed
		dialog_recrd.Show()

	})
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
		widget.NewLabel("选择时间"),
		datePicker,
		widget.NewLabel("输入贷款金额"),
		loanAmountEntry,
		submitButton,
	)
	return form
}
