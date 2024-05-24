package admin

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/service"
	"sql_bank/utils"
	"strconv"
)

var (
	accountService = service.AccountService{}
)

// 账号管理GUI
func MakeAccountManageUI(w fyne.Window, userInfo model.User) *fyne.Container {
	//账号管理界面
	var accountList []model.Account
	accountList = accountService.GetAccountList()
	//账号列表
	var accountListTable *widget.Table
	accountListTable = widget.NewTable(
		func() (int, int) {
			return len(accountList) + 1, 6 //行数为账号数量，列数固定为6
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") //每个单元格都使用Label
		},
		func(tid widget.TableCellID, co fyne.CanvasObject) {
			label := co.(*widget.Label)
			if tid.Row == 0 {
				//第一行为列名
				switch tid.Col {
				case 0:
					label.SetText("账号")
				case 1:
					label.SetText("账户类型")
				case 2:
					label.SetText("用户姓名")
				case 3:
					label.SetText("透支限额")
				case 4:
					label.SetText("信用等级")
				case 5:
					label.SetText("余额")
				}
			} else {
				//其他行为数据行，注意行号需要减1因为第一行是列名
				switch tid.Col {
				case 0:
					label.SetText(accountList[tid.Row-1].AccountNumber)
				case 1:
					label.SetText(accountList[tid.Row-1].AccountType.Name)
				case 2:
					label.SetText(userService.GetUserNameById(accountList[tid.Row-1].UserID))
				case 3:
					label.SetText(fmt.Sprintf("%f", accountList[tid.Row-1].OverdraftLimit))
				case 4:
					label.SetText(strconv.Itoa(accountList[tid.Row-1].CreditRating))
				case 5:

					label.SetText(fmt.Sprintf("%f", accountList[tid.Row-1].Balance))
				}
			}
		})

	// 设置表格的列宽
	accountListTable.SetColumnWidth(0, 150) // 第一列宽度
	accountListTable.SetColumnWidth(1, 100) // 第二列宽度
	accountListTable.SetColumnWidth(2, 100) // 第二列宽度
	accountListTable.SetColumnWidth(3, 100) // 第二列宽度
	accountListTable.SetColumnWidth(4, 100) // 第二列宽度
	accountListTable.SetColumnWidth(5, 100) // 第二列宽度

	// 包装 userList 在一个可滚动的容器中，并设置高度
	scrollContainer := container.NewVScroll(accountListTable)
	scrollContainer.SetMinSize(fyne.NewSize(200, 300)) // 设置列表的最小尺寸，其中高度为300

	//选择
	rowSelected := -1
	accountListTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row - 1
	}
	//新增修改表单
	//选择用户
	//用户列表
	var userList []model.User
	userName := make([]string, 0)
	userList = userService.GetUserList()
	for _, user := range userList {
		userName = append(userName, user.Username)
	}
	selectedUserID := uint(0)
	UserSelect := widget.NewSelect(userName, func(value string) {
		// 在回调中查找所选项的ID
		for _, user := range userList {
			if user.Username == value {
				selectedUserID = user.ID
				fmt.Println("Selected User ID:", selectedUserID)
				break
			}
		}
	})
	//选择账户类型
	//账户类型列表
	var accountTypeList []model.AccountType
	accountTypeList = cardTypeService.GetCardType()
	accountTypeName := make([]string, 0)
	for _, accountType := range accountTypeList {
		accountTypeName = append(accountTypeName, accountType.Name)
	}
	selectedAccountTypeID := uint(0)
	AccountTypeSelect := widget.NewSelect(accountTypeName, func(value string) {
		// 在回调中查找所选项的ID
		for _, accountType := range accountTypeList {
			if accountType.Name == value {
				selectedAccountTypeID = accountType.ID
				fmt.Println("Selected Account Type ID:", selectedAccountTypeID)
				break
			}
		}
	})
	//透支限额
	OverdraftLimitEntry := widget.NewEntry()

	//新增或修改表单
	addOrUpdateForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "用户", Widget: UserSelect},
			{Text: "账户类型", Widget: AccountTypeSelect},
			{Text: "透支限额", Widget: OverdraftLimitEntry},
		},
		OnSubmit: func() {
			OverdraftLimit, _ := strconv.ParseFloat(OverdraftLimitEntry.Text, 64)
			if rowSelected == -1 {
				// 新增
				account := model.Account{
					AccountNumber:  utils.GenerateRandomAccountID(10),
					UserID:         selectedUserID,
					AccountTypeID:  selectedAccountTypeID,
					OverdraftLimit: OverdraftLimit,
					PasswordHash:   "123456",
					Balance:        0,
					CreditRating:   4,
				}
				_, err := accountService.AddAccount(account)
				if err != nil {
					dialog.ShowInformation("Error", "添加账户失败", w)
					return
				}
				//accountList = append(accountList, addAccount)
				accountList = accountService.GetAccountList()
				accountListTable.Refresh()
				dialog.ShowInformation("Success", "新增账户成功", w)

			} else {
				// 修改
				account := accountList[rowSelected]
				account.UserID = selectedUserID
				account.AccountTypeID = selectedAccountTypeID
				account.OverdraftLimit = OverdraftLimit
				err := accountService.UpdateAccount(account)
				if err != nil {
					dialog.ShowInformation("Error", "更新账户失败", w)
					return
				}
				accountList = accountService.GetAccountList()
				accountListTable.Refresh()
				dialog.ShowInformation("Success", "更新账户成功", w)
			}
		},
	}

	//新增按钮
	addAccountButton := widget.NewButton("新增账户", func() {
		rowSelected = -1
		UserSelect.SetSelected("")
		AccountTypeSelect.SetSelected("")
		OverdraftLimitEntry.SetText("")

		userList = userService.GetUserList()

		userName = make([]string, 0)
		for _, user := range userList {
			userName = append(userName, user.Username)
		}
		UserSelect.Options = userName
		//更新用户列表
		UserSelect.Refresh()
		dialog.ShowForm("新增账户", "确定", "取消", addOrUpdateForm.Items, func(b bool) {
			if b {
				addOrUpdateForm.OnSubmit()
			}
		}, w)
	})
	//修改按钮
	updateAccountButton := widget.NewButton("修改账户", func() {
		if rowSelected < 0 {
			dialog.ShowInformation("Error", "No account selected", w)
			return
		}
		account := accountList[rowSelected]
		userNameById := userService.GetUserNameById(account.UserID)
		UserSelect.SetSelected(userNameById)
		AccountTypeSelect.SetSelected(account.AccountType.Name)
		OverdraftLimitEntry.SetText(fmt.Sprintf("%f", account.OverdraftLimit))
		dialog.ShowForm("修改账户", "确定", "取消", addOrUpdateForm.Items, func(b bool) {
			if b {
				addOrUpdateForm.OnSubmit()
			}
		}, w)
	})
	//删除按钮
	deleteAccountButton := widget.NewButton("删除账户", func() {

		if rowSelected < 0 {
			dialog.ShowInformation("Error", "No account selected", w)
			return
		}
		//余额不为0不能删除
		if accountList[rowSelected].Balance != 0 {
			dialog.ShowInformation("Error", "账户余额不为0，不能删除", w)
			return
		}

		err := accountService.DeleteAccount(accountList[rowSelected])
		if err != nil {
			dialog.ShowInformation("Error", "删除账户失败", w)
			return
		}
		accountList = accountService.GetAccountList()
		accountListTable.Refresh()
		dialog.ShowInformation("Success", "删除账户成功", w)
	})
	//搜索
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索...")
	searchEntry.OnChanged = func(s string) {
		accountList = accountService.SearchAccount(searchEntry.Text)
		accountListTable.Refresh()
	}

	//布局
	top := container.NewHBox(addAccountButton, updateAccountButton, deleteAccountButton)
	return container.NewVBox(searchEntry, top, scrollContainer)

}
