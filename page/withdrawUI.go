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

func MakeWithdrawUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {
	//获取用户账户信息
	info := cardService.GetAccountInfo(userInfo.ID)

	// 用户信息标签
	balanceLabel := widget.NewLabel("当前余额: N/A")
	creditLevelLabel := widget.NewLabel("信誉等级:" + string(rune(info.CreditRating)))
	overdraftRecordLabel := widget.NewLabel("透支记录: N/A")

	// 卡类别选择 根据用户id选择已经开户的卡
	cardTypeList := cardService.GetCardType(userInfo.ID)

	fmt.Println(cardTypeList)
	//卡的类型
	cardTypeOptions := make([]string, 0)
	//卡号
	cardNumberOptions := make([]string, 0)
	cardNumberSelect := widget.NewSelect(cardNumberOptions, func(value string) {
		// 在回调拿到卡号
		fmt.Println("Selected Card Number:", value)
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
		err = cardService.Saving(cardNumber, password, amount, selectedCardTypeID)
		if err != nil {
			dialog.ShowInformation("失败", err.Error(), w)
		} else {
			dialog.ShowInformation("成功", "存款成功", w)
		}
	})
	// 布局
	form := container.NewVBox(
		balanceLabel,
		creditLevelLabel,
		overdraftRecordLabel,
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
