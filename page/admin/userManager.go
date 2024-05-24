package admin

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/service"
)

var (
	userService = service.UserService{}
)

func MakeUserManageUI(w fyne.Window, userInfo model.User) *fyne.Container {
	var users []model.User
	users = userService.GetUserList()
	fmt.Println(users)
	// 使用 widget.Table 创建表格形式的用户列表
	var userTable *widget.Table
	userTable = widget.NewTable(
		func() (int, int) {
			return len(users) + 1, 3 // 行数为用户数量，列数固定为2
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
					label.SetText("用户名")
				case 1:
					label.SetText("角色")
				case 2:
					label.SetText("性别")

				}

			} else {
				// 其他行为数据行，注意行号需要减1因为第一行是列名
				switch tid.Col {
				case 0:
					label.SetText(users[tid.Row-1].Username)
				case 1:
					if users[tid.Row-1].Role == 0 {
						label.SetText("普通用户") // 0:普通用户
					}
					if users[tid.Row-1].Role == 1 {
						label.SetText("管理员") // 1:管理员
					}
				case 2:
					if users[tid.Row-1].Gender == 0 {
						label.SetText("女") // 0:女
					}
					if users[tid.Row-1].Gender == 1 {
						label.SetText("男") // 1:男
					}

				}

			}
		},
	)
	// 设置表格的列宽
	userTable.SetColumnWidth(0, 150) // 第一列宽度
	userTable.SetColumnWidth(1, 100) // 第二列宽度
	userTable.SetColumnWidth(2, 100) // 第二列宽度
	// 包装 userList 在一个可滚动的容器中，并设置高度
	scrollContainer := container.NewVScroll(userTable)
	scrollContainer.SetMinSize(fyne.NewSize(200, 300)) // 设置列表的最小尺寸，其中高度为300
	// 用户搜索
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索用户...")
	searchEntry.OnChanged = func(s string) {
		// 这里可以加上搜索逻辑，动态更新 userList 的内容
		users = userService.SearchUser(s)
		userTable.Refresh()
	}
	// 添加用户的表单
	nameEntry := widget.NewEntry()
	roleEntry := widget.NewSelect([]string{"User", "Admin"}, nil)
	genderEntry := widget.NewSelect([]string{"女", "男"}, nil)

	addForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Name", Widget: nameEntry},
			{Text: "Role", Widget: roleEntry},
			{Text: "Gender", Widget: genderEntry},
		},
		OnSubmit: func() {
			role := 0
			fmt.Println(roleEntry.Selected)
			if roleEntry.Selected == "Admin" {
				role = 1
			}

			gender := 0
			if genderEntry.Selected == "男" {
				gender = 1
			}
			fmt.Println(roleEntry, genderEntry)
			fmt.Println("Adding user", nameEntry.Text, role, gender)
			// 在这里处理添加用户的逻辑
			newUser := model.User{
				Username: nameEntry.Text,
				Role:     role,
				Password: "123456",
				Gender:   gender,
			}
			//判断用户名是否存在
			err := userService.AddUser(newUser)
			if err != nil {
				dialog.ShowInformation("Error", "用户名重复", w)
				return
			}
			users = append(users, newUser)
			//users = append(users, model.User{Username: nameEntry.Text, Role: role})
			userTable.Refresh() // 刷新列表显示新增的用户
			dialog.ShowInformation("Success", "新增用户成功", w)
		},
	}
	// 添加用户的按钮
	addUserButton := widget.NewButton("新增用户", func() {
		dialog.ShowForm("新增用户", "确定", "取消", addForm.Items, func(b bool) {
			if b {
				addForm.OnSubmit()
			}
		}, w)
	})
	userTable.ShowHeaderRow = true
	rowSelected := -1
	userTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row
	}
	deleteUserButton := widget.NewButton("删除用户", func() {
		if rowSelected < 0 {
			dialog.ShowInformation("Error", "No user selected", w)
			return
		}
		if rowSelected > 0 {
			// 删除选中的用户
			fmt.Println("删除用户", users[rowSelected-1].ID, users[rowSelected-1].Username)
			userService.DeleteUser(users[rowSelected-1].ID)
			users = append(users[:rowSelected-1], users[rowSelected:]...)
			userTable.Refresh()
			rowSelected = -1
			dialog.ShowInformation("Success", "User deleted successfully", w)
		} else {
			dialog.ShowInformation("Error", "No user selected", w)
		}

	})
	//修改表单
	updateNameEntry := widget.NewEntry()
	updateRoleEntry := widget.NewSelect([]string{"User", "Admin"}, nil)
	updateGenderEntry := widget.NewSelect([]string{"女", "男"}, nil)
	updateForm := &widget.Form{
		Items: []*widget.FormItem{
			//{Text: "Name", Widget: updateNameEntry},
			{Text: "Role", Widget: updateRoleEntry},
			{Text: "Gender", Widget: updateGenderEntry},
		},
		OnSubmit: func() {
			role := 0
			if updateRoleEntry.Selected == "Admin" {
				role = 1
			}
			gender := 0
			if updateGenderEntry.Selected == "男" {
				gender = 1
			}
			// 在这里处理修改用户的逻辑
			updateUser := model.User{
				Username: updateNameEntry.Text,
				Role:     role,
				Gender:   gender,
			}
			err := userService.UpdateUser(updateUser)
			fmt.Println(err)
			if err != nil {
				dialog.ShowInformation("Error", "修改用户失败", w)
				return
			}
			users[rowSelected-1] = updateUser
			userTable.Refresh()
			dialog.ShowInformation("Success", "修改用户成功", w)
		},
	}
	//修改按钮
	updateUserButton := widget.NewButton("修改用户", func() {
		if rowSelected < 0 {
			dialog.ShowInformation("Error", "No user selected", w)
			return
		}
		updateNameEntry.SetText(users[rowSelected-1].Username)
		if users[rowSelected-1].Role == 0 {
			updateRoleEntry.SetSelected("User")
		}
		if users[rowSelected-1].Role == 1 {
			updateRoleEntry.SetSelected("Admin")
		}
		if users[rowSelected-1].Gender == 0 {
			updateGenderEntry.SetSelected("女")
		}
		if users[rowSelected-1].Gender == 1 {
			updateGenderEntry.SetSelected("男")
		}
		dialog.ShowForm("新增用户", "确定", "取消", updateForm.Items, func(b bool) {
			if b {
				updateForm.OnSubmit()
			}
		}, w)

	})

	top := container.NewHBox(addUserButton, deleteUserButton, updateUserButton)

	return container.NewVBox(
		searchEntry,
		top,
		scrollContainer,
	)
}
