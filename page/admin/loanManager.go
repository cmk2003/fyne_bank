package admin

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/global"
	"sql_bank/model"
	"sql_bank/service"
	"strconv"
)

var (
	loanService = service.LoanService{}
)
var (
	currentPage = 1
	pageSize    = 10
)

func MakeLoanManager(w fyne.Window, userInfo model.User) *fyne.Container {

	var loanList []model.Loan
	loanList = loanService.SearchLoansByPage("", currentPage, pageSize)
	// 使用 widget.Table 创建表格形式的用户列表
	var loanTable *widget.Table
	loanTable = widget.NewTable(
		func() (int, int) {
			return len(loanList) + 1, 10 // 行数为用户数量，列数固定为2
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("") // 每个单元格都使用 Label
		},
		func(tid widget.TableCellID, co fyne.CanvasObject) {
			label := co.(*widget.Label)
			if tid.Row == 0 {
				// 第一行为列名
				switch tid.Col {
				case 0:
					label.SetText("账户ID")
				case 1:
					label.SetText("贷款金额")
				case 2:
					label.SetText("贷款日期")
				case 3:
					label.SetText("还款日期")

				case 4:
					label.SetText("卡号")
				case 5:
					label.SetText("账户所属人")
				case 6:
					label.SetText("账户类型")
				case 7:
					label.SetText("状态")
				case 8:
					label.SetText("应还款金额")
				case 9:
					label.SetText("利率")

				}
			} else {
				// 其他行为数据行，注意行号需要减1因为第一行是列名
				switch tid.Col {
				case 0:
					label.SetText(strconv.Itoa(int(loanList[tid.Row-1].AccountID)))
				case 1:
					amountStr := strconv.FormatFloat(loanList[tid.Row-1].AmountBorrowed, 'f', 2, 64)
					label.SetText(amountStr)
				case 2:
					// Set text to only year-month-day format
					label.SetText(loanList[tid.Row-1].LoanDate.Format("2006-01-02"))
				case 3:
					// Set text to only year-month-day format
					label.SetText(loanList[tid.Row-1].DueDate.Format("2006-01-02"))

				case 4:
					label.SetText(loanList[tid.Row-1].Account.AccountNumber)
				case 5:
					label.SetText(loanList[tid.Row-1].Account.User.RealName)
				case 6:
					label.SetText(loanList[tid.Row-1].Account.AccountType.Name)
				case 7:
					if loanList[tid.Row-1].Status == true {
						label.SetText("已还款")
					} else {
						label.SetText("未还款")
					}
				case 8:
					interestStr := strconv.FormatFloat(loanList[tid.Row-1].InterestAccrued+loanList[tid.Row-1].AmountBorrowed, 'f', 2, 64)
					label.SetText(interestStr)
				case 9:
					rateStr := strconv.FormatFloat(loanList[tid.Row-1].InterestRate, 'f', 2, 64)
					label.SetText(rateStr)

				}
			}
		},
	)
	// 设置表格的列宽
	loanTable.SetColumnWidth(0, 100)
	loanTable.SetColumnWidth(1, 100)
	loanTable.SetColumnWidth(2, 200)
	loanTable.SetColumnWidth(3, 200)
	loanTable.SetColumnWidth(4, 100)
	loanTable.SetColumnWidth(5, 100)
	loanTable.SetColumnWidth(6, 100)
	loanTable.SetColumnWidth(7, 100)
	loanTable.SetColumnWidth(8, 100)
	loanTable.SetColumnWidth(9, 100)
	// 包装 userList 在一个可滚动的容器中，并设置高度
	scrollContainer := container.NewVScroll(loanTable)
	scrollContainer.SetMinSize(fyne.NewSize(200, 500)) // 设置列表的最小尺寸，其中高度为300

	rowSelected := -1
	loanTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row - 1
	}
	//搜索
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索...")
	searchEntry.OnChanged = func(s string) {
		// 这里可以加上搜索逻辑，动态更新 userList 的内容
		loanList = loanService.SearchLoansByPage(s, currentPage, pageSize)
	}
	updateTable := func(page int) {
		loanList = loanService.SearchLoansByPage(searchEntry.Text, currentPage, pageSize)
		loanTable.Refresh() // 刷新表格显示新的数据
	}
	// 分页按钮
	nextPageButton := widget.NewButton("Next", func() {
		currentPage++
		updateTable(currentPage)
	})
	previousPageButton := widget.NewButton("Previous", func() {
		if currentPage > 1 {
			currentPage--
			updateTable(currentPage)
		}
	})
	//给账户id提示还款
	repayButton := widget.NewButton("还款", func() {
		if rowSelected < 0 {
			dialog.ShowInformation("Error", "No user selected", w)
			return
		}
		if rowSelected > 0 {
			//打印用户id
			loanID := loanList[rowSelected].Account.User.ID
			//还款
			fmt.Println("提醒还款", loanID)
			// 给useid提示还款 redis设置可以
			global.RDB.Set(global.Ctx, "repay:"+strconv.Itoa(int(loanID)), "1", -1)
			dialog.ShowInformation("Success", "提醒还款成功", w)
		}
	})

	//top := container.NewHBox(searchEntry)

	return container.NewVBox(
		//top,
		searchEntry,
		container.NewHBox(
			previousPageButton,
			nextPageButton,
			repayButton),
		scrollContainer,
	)
	// 组装顶部控件
}
