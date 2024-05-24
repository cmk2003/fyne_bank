package page

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/page/admin"
	"sql_bank/service"
)

var userService = service.UserService{}

func MakeLoginUI(a fyne.App) {
	w := a.NewWindow("银行登陆")

	username := widget.NewEntry()
	username.SetPlaceHolder("username")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("pass")

	var loginButton *widget.Button // 先声明变量
	loginButton = widget.NewButton("login", func() {
		fmt.Println(username.Text, password.Text)
		username.Text = "root"
		password.Text = "123456"
		if isCor, userInfo := userService.LoginSys(username.Text, password.Text); isCor {
			fmt.Println("success")
			//获取用户id
			//如果是用户，跳转到用户界面
			if userInfo.Role == 0 {
				MakeMainUI(a, userInfo)
			} else if userInfo.Role == 1 {
				admin.MakeAdminMainUI(a, userInfo)
			}

			//如果是管理员调转到管理员界面
			w.Close() // 关闭登录窗口
		} else {
			fmt.Println("fail")
			dialog.ShowInformation("失败", "用户名或密码错误", w)
		}
	})
	w.SetContent(container.NewVBox(
		username,
		password,
		loginButton,
	))
	w.Resize(fyne.NewSize(400, 300))
	w.Show()
}
