package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/page/admin"
)

func MakeAdminMainUI(a fyne.App, userInfo model.User) {
	w := a.NewWindow("后台管理")
	//创建用户管理界面
	userManageUI := admin.MakeUserManageUI(w, userInfo)
	//创建银行卡管理界面
	cardManageUI := admin.MakeCardManageUI(w, userInfo)
	// 创建账号管理
	accountManageUI := admin.MakeAccountManageUI(w, userInfo)

	// 创建贷款管理
	loanManageUI := admin.MakeLoanManager(w, userInfo)

	//退出系统按钮
	exitButton := widget.NewButton("退出系统", func() {
		w.Close()
		MakeLoginUI(a) // 显示登录窗口
	})

	//使用Tabs容器来组织不同的功能界面
	tabs := container.NewAppTabs(
		container.NewTabItem("用户管理", userManageUI),
		container.NewTabItem("银行卡管理", cardManageUI),
		container.NewTabItem("账号管理", accountManageUI),
		container.NewTabItem("贷款管理", loanManageUI),
	)
	// 创建一个顶部的容器放置退出按钮
	topBar := container.NewHBox(layout.NewSpacer(), exitButton) // 使用Spacer让按钮保持在右侧

	// 将标签和顶部的退出按钮放在一个垂直的布局中
	mainContent := container.NewBorder(topBar, nil, nil, nil, tabs)

	w.SetContent(mainContent)
	w.Resize(fyne.NewSize(800, 600))
	w.Show()
}
