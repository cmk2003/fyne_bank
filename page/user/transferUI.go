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
	transferService = service.TransferService{}
	userService     = service.UserService{}
)

func MakeTransferUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {
	info := cardService.GetAccountInfo(userInfo.ID)
	fmt.Println(info)
	balanceLabel := widget.NewLabel("当前余额: N/A")
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

	// 转账金额输入
	depositAmountEntry := widget.NewEntry()
	depositAmountEntry.SetPlaceHolder("请输入转账金额")

	//输入对方卡号
	transferAccountEntry := widget.NewEntry()
	transferAccountEntry.SetPlaceHolder("请输入对方卡号")

	//提交
	submitButton := widget.NewButton("转账", func() {
		fmt.Println("转账金额:", depositAmountEntry.Text)
		fmt.Println("对方卡号:", transferAccountEntry.Text)
		fmt.Println("密码:", passwordEntry.Text)
		//自己的卡号
		fmt.Println("自己的卡号:", cardNumberSelect.Selected)
		//类别
		fmt.Println("类别:", cardTypeSelect.Selected)
		// 验证字段
		if cardNumberSelect.Selected == "" || cardTypeSelect.Selected == "" || depositAmountEntry.Text == "" || transferAccountEntry.Text == "" || passwordEntry.Text == "" {
			dialog.ShowInformation("错误", "请填写完整信息", w)
			return
		}
		//判断对方账户是否存在
		account := cardService.GetAccountByCardNumber(transferAccountEntry.Text)
		if account.ID == 0 {
			dialog.ShowInformation("错误", "对方账户不存在", w)
			return
		}
		// 验证密码
		if !cardService.VerifyPassword(cardNumberSelect.Selected, passwordEntry.Text) {
			// 转账
			dialog.ShowInformation("错误", "密码错误", w)
			return
		}
		// 验证金额
		if depositAmountEntry.Text <= "0" {
			dialog.ShowInformation("错误", "请输入转账金额", w)
			return
		}
		// 验证余额
		depositAmount, _ := strconv.ParseFloat(depositAmountEntry.Text, 64)
		if info.Balance < depositAmount {
			dialog.ShowInformation("错误", "余额不足1", w)
			return
		}
		//获取对方账户的用户名
		transferAccount := cardService.GetAccountByCardNumber(transferAccountEntry.Text)
		fmt.Println("transferAccount", userService.GetUserNameById(transferAccount.UserID))
		//获取账户id
		accountId := cardService.GetAccountByCardNumber(cardNumberSelect.Selected).ID
		dialog.NewConfirm("转账确认", "确定转账给"+userService.GetUserNameById(transferAccount.UserID), func(b bool) {
			if b {
				// 转账
				//根据卡号获取账户类型
				accountType1 := cardService.GetAccountByCardNumber(cardNumberSelect.Selected).AccountType
				accountType2 := cardService.GetAccountByCardNumber(transferAccountEntry.Text).AccountType
				//提取比例
				tem_bili := 0.01
				if accountType1 == accountType2 {
					tem_bili = 0.002
				} else {
					tem_bili = 0.05
				}
				//判断余额
				if info.Balance < depositAmount+depositAmount*tem_bili {
					dialog.ShowInformation("错误", "余额不足", w)
					return
				}
				//交易表增加一条记录
				transaction := model.Transaction{
					CardNumber:      cardNumberSelect.Selected,
					ToCardNumber:    transferAccountEntry.Text,
					Amount:          depositAmount,
					TransactionType: 1,
					AccountID:       accountId,
					Status:          "fail",
					TransactionDate: time.Now(),
				}
				//开启事务
				err := cardService.Transfer(cardNumberSelect.Selected, transferAccountEntry.Text, depositAmountEntry.Text, transaction)
				if err != nil {
					dialog.ShowInformation("错误", err.Error(), w)
				} else {
					dialog.ShowInformation("成功", "转账成功", w)
					// 更新余额
					info = cardService.GetAccountInfo(userInfo.ID)
					// 刷新余额
					balanceLabel.SetText("当前余额: " + strconv.FormatFloat(info.Balance, 'f', 2, 64))

				}
				//收取交易费用
				transactionGiveToBank := model.Transaction{
					CardNumber:      cardNumberSelect.Selected,
					ToCardNumber:    "00000000",
					Amount:          depositAmount * tem_bili,
					TransactionType: 1,
					AccountID:       accountId,
					Status:          "fail",
					TransactionDate: time.Now(),
				}
				err = cardService.Transfer(cardNumberSelect.Selected, "00000000", strconv.FormatFloat(depositAmount*tem_bili, 'f', 2, 64), transactionGiveToBank)
			}
		}, w).Show()

	})
	top := container.NewHBox(
		balanceLabel,
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
		widget.NewLabel("输入对方卡号"),
		transferAccountEntry,
		widget.NewLabel("输入转账金额"),
		depositAmountEntry,
		submitButton,
	)
	return form
}
