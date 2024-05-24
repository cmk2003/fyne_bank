package user

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"strconv"
	"time"
)

func MakeRepayDraftUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {
	//获取用户账户信息
	info := cardService.GetAccountInfo(userInfo.ID)
	fmt.Println(info)
	// 用户信息标签
	balanceLabel := widget.NewLabel("当前余额: N/A")
	// 根据用户id选择已经开户的卡
	cardTypeList := cardService.GetCardType(userInfo.ID)

	//给每张卡展示还款信息
	var draftList []model.Overdraft
	//loanList = loanService.GetLoanList(userInfo.ID)
	var draftTable *widget.Table
	draftTable = widget.NewTable(
		func() (int, int) {
			// 卡号 + 还款金额 + 还款时间 + 还款状态 + 应还款金额
			return len(draftList) + 1, 4 // 行数为用户数量，列数固定为2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // 每个单元格都使用 Label
		},
		func(tid widget.TableCellID, co fyne.CanvasObject) {
			// 卡号 + 还款金额 + 还款时间 + 还款状态 + 应还款金额
			label := co.(*widget.Label)
			if tid.Row == 0 {
				// 第一行为列名
				switch tid.Col {
				case 0:
					label.SetText("卡号")
				case 1:
					label.SetText("透支金额")
				case 2:
					label.SetText("应还款时间")
				case 3:
					label.SetText("还款状态")
				}

			} else {
				// 其他行为数据行，注意行号需要减1因为第一行是列名
				switch tid.Col {
				case 0:
					label.SetText(draftList[tid.Row-1].Account.AccountNumber)
				case 1:
					label.SetText(fmt.Sprintf("%f", draftList[tid.Row-1].Amount))
				case 2:
					label.SetText(draftList[tid.Row-1].RepaymentDueDate.String())
				case 3:
					if draftList[tid.Row-1].Repaid == false {
						//超时还款
						if draftList[tid.Row-1].RepaymentDueDate.Before(time.Now()) {
							label.SetText("超时未还款")
						} else {
							label.SetText("未还款")
						}
					} else {
						label.SetText("已还款")
					}
				}
			}
		},
	)
	// 设置表格的列宽
	draftTable.SetColumnWidth(0, 100)
	draftTable.SetColumnWidth(1, 100)
	draftTable.SetColumnWidth(2, 220)
	draftTable.SetColumnWidth(3, 100)
	// 包装 userList 在一个可滚动的容器中，并设置高度
	scrollContainer := container.NewVScroll(draftTable)
	scrollContainer.SetMinSize(fyne.NewSize(500, 300)) // 设置列表的最小尺寸，其中高度为300

	rowSelected := -1
	draftTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row - 1
	}
	// 点击事件
	draftTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row - 1
		if rowSelected >= 0 { // 确保选中的是有效的数据行，不包括标题行
			// 获取选中的卡号
			cardNumber := draftList[rowSelected].Account.AccountNumber
			//获取卡的详细信息
			cardInfo := cardService.GetAccountByCardNumber(cardNumber)
			//更新余额
			balanceLabel.SetText("当前余额: " + strconv.FormatFloat(cardInfo.Balance, 'f', 2, 64))

		}
	}
	fmt.Println("cardTypeList", cardTypeList)
	//卡的类型
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
				//获取根据卡类别选择卡号userid和typeid 获取account信息 在一个类型下可能有多个卡
				cardInfoList := cardService.GetCardNumber(userInfo.ID, selectedCardTypeID)
				var accountIds []uint
				for _, cardInfo := range cardInfoList {
					accountIds = append(accountIds, cardInfo.ID)
				}
				fmt.Println(accountIds)
				//获取透支信息
				draftList = OverdraftService.GetOverdraftList(accountIds)

				fmt.Println(draftList)
				//更新表格
				draftTable.Refresh()
				break
			}
		}
	})

	// 用户密码输入
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("请输入密码")

	// 还款按钮
	repayButton := widget.NewButton("还款", func() {
		if rowSelected == -1 {
			dialog.ShowInformation("Error", "No user selected", w)
			return
		}
		if passwordEntry.Text == "" {
			dialog.ShowInformation("Error", "请输入密码", w)
			return
		}
		fmt.Println("rowSelected", rowSelected)
		// 密码确认
		cardNumber := draftList[rowSelected].Account.AccountNumber
		if cardService.VerifyPassword(cardNumber, passwordEntry.Text) == false {
			dialog.ShowInformation("Error", "密码错误", w)
			return
		}
		//确认框
		dialog.NewConfirm("确认还款", "确认还款吗", func(b bool) {

			if b {
				//检测余额
				if cardService.GetAccountByCardNumber(cardNumber).Balance < draftList[rowSelected].Amount {
					dialog.ShowInformation("Error", "余额不足", w)
					return
				}

				//还款
				transaction := model.Transaction{
					Amount:          draftList[rowSelected].Amount,
					TransactionType: 7,
					AccountID:       draftList[rowSelected].Account.ID,
					TransactionDate: time.Now(),
					CardNumber:      cardNumber,
					ToCardNumber:    "00000000",
					Status:          "fail",
				}
				amount := draftList[rowSelected].Amount
				amountStr := strconv.FormatFloat(amount, 'f', 2, 64)
				err := cardService.Transfer(cardNumber, "00000000", amountStr, transaction)
				//标记还款成功
				draftList[rowSelected].Repaid = true
				//更新透支表
				err = OverdraftService.UpdateOverdraft(draftList[rowSelected])
				//透支额度加上去
				cardInfo := cardService.GetAccountByCardNumber(cardNumber)
				//还款时间在规定时间之前
				if draftList[rowSelected].RepaymentDueDate.After(time.Now()) {
					cardInfo.OverdraftLimit += amount * 1.1
				} else {
					cardInfo.OverdraftLimit += amount * 0.9
				}
				cardService.UpdateAccount(cardInfo)
				if err != nil {
					dialog.ShowInformation("Error", err.Error(), w)
					return
				}
				//更新余额
				balanceLabel.SetText("当前余额: " + strconv.FormatFloat(cardService.GetAccountByCardNumber(cardNumber).Balance, 'f', 2, 64))
				//更新表格
				draftList = OverdraftService.GetOverdraftList([]uint{draftList[rowSelected].Account.ID})
				draftTable.Refresh()
				dialog.ShowInformation("Success", "还款成功", w)
			}
		}, w).Show()
	})
	searchNoPayButton := widget.NewButton("查看未还款", func() {
		//获取透支信息
		//获取根据卡类别选择卡号userid和typeid 获取account信息 在一个类型下可能有多个卡
		cardInfoList := cardService.GetCardNumber(userInfo.ID, selectedCardTypeID)
		var accountIds []uint
		for _, cardInfo := range cardInfoList {
			accountIds = append(accountIds, cardInfo.ID)
		}
		fmt.Println(accountIds)
		//获取透支信息
		draftList = OverdraftService.GetOverdraftListWithNoPay(accountIds)

		fmt.Println(draftList)
		//更新表格
		draftTable.Refresh()
	})
	top := container.NewHBox(
		cardTypeSelect, repayButton, searchNoPayButton,
	)

	return container.NewVBox(
		balanceLabel,
		top,
		container.NewVBox(
			widget.NewLabel("输入当前选择卡号的密码:"),
			passwordEntry,
		),
		scrollContainer,
	)

}
