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

var cardService service.CardService

func MakeDepositUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {
	// 卡类别选择 根据用户id选择已经开户的卡
	cardTypeList := cardService.GetCardType(userInfo.ID)
	fmt.Println(cardTypeList)
	cardTypeOptions := make([]string, 0)
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
				selectedCardTypeID = cardType.ID
				fmt.Println("Selected Card Type ID:", selectedCardTypeID)
				break
			}
		}
	})
	// 卡号输入
	cardNumberEntry := widget.NewEntry()
	cardNumberEntry.SetPlaceHolder("请输入卡号")

	// 用户密码输入
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// 存款金额输入
	depositAmountEntry := widget.NewEntry()
	depositAmountEntry.SetPlaceHolder("请输入存款金额")

	// 提交按钮
	submitButton := widget.NewButton("提交", func() {
		cardType := cardTypeSelect.Selected
		cardNumber := cardNumberEntry.Text
		password := passwordEntry.Text
		depositAmount := depositAmountEntry.Text

		if cardType == "" || cardNumber == "" || password == "" || depositAmount == "" {
			dialog.ShowInformation("错误", "所有字段均为必填项", w)
			return
		}

		amount, err := strconv.ParseFloat(depositAmount, 64)
		if err != nil || amount <= 0 {
			dialog.ShowInformation("错误", "存款金额必须是正数", w)
			return
		}
		err = cardService.Saving(cardNumber, password, amount, selectedCardTypeID)
		//交易表插入数据
		if err != nil {
			dialog.ShowInformation("失败", err.Error(), w)
		} else {
			dialog.ShowInformation("成功", "存款成功", w)
		}
	})
	// 布局
	form := container.NewVBox(
		widget.NewLabel("选择卡类别"),
		cardTypeSelect,
		widget.NewLabel("输入卡号"),
		cardNumberEntry,
		widget.NewLabel("输入密码"),
		passwordEntry,
		widget.NewLabel("输入存款金额"),
		depositAmountEntry,
		submitButton,
	)

	return form
}
