package common

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/page/user"
	"sql_bank/utils"
	"strconv"
)

func MakeMainUI(a fyne.App, userInfo model.User) {

	w := a.NewWindow("Bank")
	// 创建存款界面
	depositUI := user.MakeDepositUI(w, userInfo)

	// 创建取款界面
	withdrawUI := user.MakeWithdrawUI(w, userInfo)

	// 创建贷款界面
	loanUI := user.MakeLoanUI(w, userInfo)

	// 创建转账界面
	transferUI := user.MakeTransferUI(w, userInfo)

	// 创建还贷界面
	repayUI := user.MakeRepayUI(w, userInfo)

	// 创建还超额度界面
	repayDraftUI := user.MakeRepayDraftUI(w, userInfo)
	// 创建修改密码
	changePasswordUI := user.MakeChangePasswordUI(w, userInfo)

	//退出系统按钮
	exitButton := widget.NewButton("退出系统", func() {
		w.Close()
		MakeLoginUI(a) // 显示登录窗口
	})
	// 判断键是否存
	fmt.Println(userInfo.ID)
	key := fmt.Sprintf("repay:%s", strconv.Itoa(int(userInfo.ID)))
	fmt.Println(key)
	if utils.CheckRedisKeyExist(key) {
		//如果存在，弹出还款窗口
		dialog.ShowInformation("还款提醒", "您有未还款的贷款，请及时还款", w)
	}

	// 使用Tabs容器来组织不同的功能界面
	tabs := container.NewAppTabs(
		container.NewTabItem("存款", depositUI),
		container.NewTabItem("取钱", withdrawUI),
		container.NewTabItem("贷款", loanUI),
		container.NewTabItem("转账", transferUI),
		container.NewTabItem("还贷款", repayUI),
		container.NewTabItem("还超额度", repayDraftUI),
		container.NewTabItem("修改密码", changePasswordUI),
	)
	// 创建一个顶部的容器放置退出按钮
	topBar := container.NewHBox(layout.NewSpacer(), exitButton) // 使用Spacer让按钮保持在右侧

	// 将标签和顶部的退出按钮放在一个垂直的布局中
	mainContent := container.NewBorder(topBar, nil, nil, nil, tabs)

	w.SetContent(mainContent)
	w.Resize(fyne.NewSize(800, 600))
	//w.ShowAndRun()
	w.Show()
}
