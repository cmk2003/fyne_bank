package admin

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"sql_bank/model"
)

func MakeAdminMainUI(a fyne.App, userInfo model.User) {
	w := a.NewWindow("后台管理")
	//创建用户管理界面
	userManageUI := MakeUserManageUI(w, userInfo)
	//创建银行卡管理界面
	cardManageUI := MakeCardManageUI(w, userInfo)
	// 创建账号管理
	accountManageUI := MakeAccountManageUI(w, userInfo)
	////创建贷款管理界面
	//loanManageUI := MakeLoanManageUI(w, userInfo)
	////创建转账管理界面
	//transferManageUI := MakeTransferManageUI(w, userInfo)
	////创建还贷管理界面
	//repayManageUI := MakeRepayManageUI(w, userInfo)
	//使用Tabs容器来组织不同的功能界面
	tabs := container.NewAppTabs(
		container.NewTabItem("用户管理", userManageUI),
		container.NewTabItem("银行卡管理", cardManageUI),
		container.NewTabItem("账号管理", accountManageUI),
		//container.NewTabItem("贷款管理", loanManageUI),
		//container.NewTabItem("转账管理", transferManageUI),
		//container.NewTabItem("还贷管理", repayManageUI),
	)
	w.SetContent(tabs)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
}
