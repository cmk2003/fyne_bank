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

func MakeRepayUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {
	//获取用户账户信息
	info := cardService.GetAccountInfo(userInfo.ID)
	fmt.Println(info)
	// 用户信息标签
	balanceLabel := widget.NewLabel("当前余额: N/A")
	// 根据用户id选择已经开户的卡
	cardTypeList := cardService.GetCardType(userInfo.ID)

	//给每张卡展示还款信息
	var loanList []model.Loan
	//loanList = loanService.GetLoanList(userInfo.ID)
	var loanTable *widget.Table
	loanTable = widget.NewTable(
		func() (int, int) {
			// 卡号 + 还款金额 + 还款时间 + 还款状态 + 应还款金额
			return len(loanList) + 1, 5 // 行数为用户数量，列数固定为2
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
					label.SetText("借贷金额")
				case 2:
					label.SetText("还款时间")
				case 3:
					label.SetText("还款状态")
				case 4:
					label.SetText("应还款金额")
				}

			} else {
				// 其他行为数据行，注意行号需要减1因为第一行是列名
				switch tid.Col {
				case 0:
					label.SetText(loanList[tid.Row-1].Account.AccountNumber)
				case 1:
					label.SetText(fmt.Sprintf("%f", loanList[tid.Row-1].AmountBorrowed))
				case 2:
					label.SetText(loanList[tid.Row-1].DueDate.String())
				case 3:
					if loanList[tid.Row-1].Status == false {
						//超时还款
						if loanList[tid.Row-1].DueDate.Before(time.Now()) {
							label.SetText("超时未还款")
						} else {
							label.SetText("未还款")
						}
					} else {
						label.SetText("已还款")
					}
				case 4:
					label.SetText(fmt.Sprintf("%f", loanList[tid.Row-1].AmountBorrowed+loanList[tid.Row-1].InterestAccrued))
				}
			}
		},
	)
	// 设置表格的列宽
	loanTable.SetColumnWidth(0, 100)
	loanTable.SetColumnWidth(1, 100)
	loanTable.SetColumnWidth(2, 200)
	loanTable.SetColumnWidth(3, 100)
	loanTable.SetColumnWidth(4, 100)
	// 包装 userList 在一个可滚动的容器中，并设置高度
	scrollContainer := container.NewVScroll(loanTable)
	scrollContainer.SetMinSize(fyne.NewSize(500, 300)) // 设置列表的最小尺寸，其中高度为300

	rowSelected := -1
	loanTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row - 1
	}
	// 点击事件
	loanTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row - 1
		if rowSelected >= 0 { // 确保选中的是有效的数据行，不包括标题行
			// 获取选中的卡号
			cardNumber := loanList[rowSelected].Account.AccountNumber
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
				//获取贷款信息
				loanList = loanService.GetLoanList(accountIds)
				fmt.Println(loanList)
				//更新表格
				loanTable.Refresh()
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
		//确认框
		dialog.NewConfirm("确认还款", "确认还款吗", func(b bool) {
			if b {
				info := cardService.GetAccountInfo(userInfo.ID)
				//还款 验证余额是否足够
				if info.Balance < loanList[rowSelected].AmountBorrowed+loanList[rowSelected].InterestAccrued {
					dialog.ShowInformation("Error", "余额不足", w)
					return
				}
				//验证密码
				if !cardService.VerifyPassword(loanList[rowSelected].Account.AccountNumber, passwordEntry.Text) {
					// 转账
					dialog.ShowInformation("错误", "密码错误", w)
					return
				}
				//当前行是否还款
				if loanList[rowSelected].Status == true {
					dialog.ShowInformation("Error", "已还款", w)
					return
				}
				//还款
				transaction := model.Transaction{
					AccountID:       loanList[rowSelected].AccountID,
					Amount:          loanList[rowSelected].AmountBorrowed + loanList[rowSelected].InterestAccrued,
					CardNumber:      loanList[rowSelected].Account.AccountNumber,
					TransactionType: 5,
					TransactionDate: time.Now(),
					Status:          "fail",
					ToCardNumber:    "00000000",
				}

				//还款
				amount := loanList[rowSelected].AmountBorrowed + loanList[rowSelected].InterestAccrued
				amountStr := strconv.FormatFloat(amount, 'f', 2, 64)

				err := cardService.Transfer(loanList[rowSelected].Account.AccountNumber, "00000000", amountStr, transaction)

				//更新还款状态
				loanList[rowSelected].Status = true
				loanService.UpdateLoan(loanList[rowSelected])
				//还款时间在规定时间之前的一个月 增加等级
				if loanList[rowSelected].DueDate.AddDate(0, -1, 0).After(time.Now()) {
					//增加信用等级
					cardInfo := cardService.GetAccountByCardNumber(loanList[rowSelected].Account.AccountNumber)
					cardInfo.CreditRating++
					cardService.UpdateAccount(cardInfo)
				}
				if err != nil {
					dialog.ShowInformation("Error", err.Error(), w)
				} else {
					dialog.ShowInformation("Success", "还款成功", w)
					// 更新余额
					info = cardService.GetAccountInfo(userInfo.ID)
					// 刷新余额
					balanceLabel.SetText("当前余额: " + strconv.FormatFloat(info.Balance, 'f', 2, 64))
					//更新表格
					loanList = loanService.GetLoanList([]uint{loanList[rowSelected].AccountID})
					loanTable.Refresh()

				}
			}
		}, w).Show()
	})

	top := container.NewHBox(
		cardTypeSelect, repayButton,
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
