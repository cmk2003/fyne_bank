package page

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"sql_bank/model"
)

func MakeMainUI(a fyne.App, userInfo model.User) {
	w := a.NewWindow("Bank")
	// 创建存款界面
	depositUI := MakeDepositUI(w, userInfo)

	// 创建取款界面
	withdrawUI := MakeWithdrawUI(w, userInfo)

	// 创建贷款界面
	loanUI := MakeLoanUI(w, userInfo)

	// 创建转账界面
	transferUI := MakeTransferUI(w, userInfo)

	// 创建还贷界面
	repayUI := MakeRepayUI(w, userInfo)

	// 使用Tabs容器来组织不同的功能界面
	tabs := container.NewAppTabs(
		container.NewTabItem("存款", depositUI),
		container.NewTabItem("取钱", withdrawUI),
		container.NewTabItem("贷款", loanUI),
		container.NewTabItem("转账", transferUI),
		container.NewTabItem("还贷款", repayUI),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	//w.ShowAndRun()
	w.Show()
}
