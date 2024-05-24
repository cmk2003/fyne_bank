package user

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/utils"
)

// 修改密码
func MakeChangePasswordUI(w fyne.Window, userInfo model.User) fyne.CanvasObject {
	// 旧密码输入框
	oldPasswordEntry := widget.NewEntry()
	oldPasswordEntry.SetPlaceHolder("请输入旧密码")
	// 新密码输入框
	newPasswordEntry := widget.NewPasswordEntry()
	newPasswordEntry.SetPlaceHolder("请输入新密码")
	// 确认新密码输入框
	confirmPasswordEntry := widget.NewPasswordEntry()
	confirmPasswordEntry.SetPlaceHolder("请再次输入新密码")
	// 修改密码按钮
	changePasswordButton := widget.NewButton("修改密码", func() {
		// 获取输入的旧密码
		oldPassword := oldPasswordEntry.Text
		// 获取输入的新密码
		newPassword := newPasswordEntry.Text
		// 获取输入的确认新密码
		confirmPassword := confirmPasswordEntry.Text
		//解析密码
		fmt.Println(userInfo.Password)
		_, salt, _, err2 := utils.ParseDjangoHash(userInfo.Password)
		if err2 != nil {
			dialog.ShowInformation("提示", "解析密码失败", w)
			return
		}
		password_old_hash := utils.EncryptPassword(oldPassword, salt)
		// 判断旧密码是否正确
		if password_old_hash != userInfo.Password {
			dialog.ShowInformation("提示", "旧密码错误", w)
			return
		}
		// 判断新密码和确认新密码是否一致
		if newPassword != confirmPassword {
			dialog.ShowInformation("提示", "两次输入的新密码不一致", w)
			return
		}
		hash := utils.DjangoHash(newPassword)
		// 更新用户密码
		userInfo.Password = hash
		err := userService.ChangePass(userInfo)
		if err != nil {
			dialog.ShowInformation("提示", "修改密码失败", w)
			return
		}
		dialog.ShowInformation("提示", "修改密码成功", w)
		//// 返回登录页面
		//w.Close()
		//// 显示登录窗口
		//MakeLoginUI(a)
	})

	// 修改密码界面布局
	changePasswordContent := container.NewVBox(
		widget.NewLabel("修改密码"),
		oldPasswordEntry,
		newPasswordEntry,
		confirmPasswordEntry,
		changePasswordButton,
	)
	return container.NewVBox(
		changePasswordContent,
	)
}
