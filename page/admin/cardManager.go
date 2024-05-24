package admin

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"sql_bank/model"
	"sql_bank/service"
	"strconv"
)

var (
	cardTypeService = service.CardTypeService{}
)

func MakeCardManageUI(w fyne.Window, userInfo model.User) *fyne.Container {
	var cardTypeList []model.AccountType
	cardTypeList = cardTypeService.GetCardType()
	fmt.Println(cardTypeList)
	// 卡类型列表
	var cardTypeListTable *widget.Table
	cardTypeListTable = widget.NewTable(
		func() (int, int) {
			return len(cardTypeList) + 1, 4 // 行数为卡类型数量，列数固定为3
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
					label.SetText("卡类型")
				case 1:
					label.SetText("利率")
				case 2:
					label.SetText("是否允许借贷")
				case 3:
					label.SetText("描述")

				}

			} else {
				// 其他行为数据行，注意行号需要减1因为第一行是列名
				switch tid.Col {
				case 0:
					label.SetText(cardTypeList[tid.Row-1].Name)
				case 1:
					label.SetText(fmt.Sprintf("%f", cardTypeList[tid.Row-1].InterestRate))
				case 2:
					if cardTypeList[tid.Row-1].OverdraftPolicy == true {
						label.SetText("是")
					}
					if cardTypeList[tid.Row-1].OverdraftPolicy == false {
						label.SetText("否")
					}
				case 3:
					label.SetText(cardTypeList[tid.Row-1].Description)

				}
			}

		},
	)
	// 设置表格的列宽
	cardTypeListTable.SetColumnWidth(0, 150) // 第一列宽度
	cardTypeListTable.SetColumnWidth(1, 100) // 第二列宽度
	cardTypeListTable.SetColumnWidth(2, 100) // 第二列宽度
	cardTypeListTable.SetColumnWidth(3, 100) // 第二列宽度

	// 包装 userList 在一个可滚动的容器中，并设置高度
	scrollContainer := container.NewVScroll(cardTypeListTable)
	scrollContainer.SetMinSize(fyne.NewSize(200, 300)) // 设置列表的最小尺寸，其中高度为300

	//选择
	rowSelected := -1
	cardTypeListTable.OnSelected = func(id widget.TableCellID) {
		rowSelected = id.Row - 1
	}
	//新增修改表单
	nameEntry := widget.NewEntry()
	descriptionEntry := widget.NewEntry()
	interestRateEntry := widget.NewEntry()
	overdraftPolicy := true
	overdraftPolicyEntry := widget.NewCheck("是否允许借贷", func(checked bool) {
		overdraftPolicy = checked
	})
	addOrUpdateForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "卡类型", Widget: nameEntry},
			{Text: "描述", Widget: descriptionEntry},
			{Text: "利率", Widget: interestRateEntry},
			{Text: "是否允许借贷", Widget: overdraftPolicyEntry},
		},
		OnSubmit: func() {
			interestRate, _ := strconv.ParseFloat(interestRateEntry.Text, 64)
			if rowSelected == -1 {
				// 新增
				cardType := model.AccountType{
					Name:            nameEntry.Text,
					Description:     descriptionEntry.Text,
					InterestRate:    interestRate,
					OverdraftPolicy: overdraftPolicy,
				}
				addCardType, err := cardTypeService.AddCardType(cardType)
				//err := addCardType
				if err != nil {
					dialog.ShowInformation("Error", "添加卡类型失败", w)
					return
				}
				cardTypeList = append(cardTypeList, addCardType)
				cardTypeListTable.Refresh()
				dialog.ShowInformation("Success", "新增卡类型成功", w)

			} else {
				// 修改
				cardType := cardTypeList[rowSelected]
				fmt.Println("cardType1:", cardType)
				cardType.Name = nameEntry.Text
				cardType.Description = descriptionEntry.Text
				cardType.InterestRate = interestRate
				cardType.OverdraftPolicy = overdraftPolicy
				//cardType = model.AccountType{
				//	Name:            nameEntry.Text,
				//	Description:     descriptionEntry.Text,
				//	InterestRate:    interestRate,
				//	OverdraftPolicy: overdraftPolicy,
				//}
				fmt.Println(cardType)
				err := cardTypeService.UpdateCardType(cardType)
				if err != nil {
					dialog.ShowInformation("Error", "修改卡类型失败", w)
					return
				}
				cardTypeList[rowSelected] = cardType
				cardTypeListTable.Refresh()
				dialog.ShowInformation("Success", "修改卡类型成功", w)
			}
			cardTypeList = cardTypeService.GetCardType()
			cardTypeListTable.Refresh()
		},
	}
	// 添加用户的按钮
	addCardTypeButton := widget.NewButton("新增卡类型", func() {
		dialog.ShowForm("新增卡类型", "确定", "取消", addOrUpdateForm.Items, func(b bool) {
			if b {
				addOrUpdateForm.OnSubmit()
			}
		}, w)
	})
	// 修改用户的按钮
	updateCardTypeButton := widget.NewButton("修改卡类型", func() {
		if rowSelected < 0 {
			dialog.ShowInformation("Error", "No card type selected", w)
			return
		}
		cardType := cardTypeList[rowSelected]
		nameEntry.SetText(cardType.Name)
		descriptionEntry.SetText(cardType.Description)
		interestRateEntry.SetText(fmt.Sprintf("%f", cardType.InterestRate))
		overdraftPolicyEntry.SetChecked(cardType.OverdraftPolicy)
		fmt.Println("cardType:", cardType)
		dialog.ShowForm("修改卡类型", "确定", "取消", addOrUpdateForm.Items, func(b bool) {
			if b {
				addOrUpdateForm.OnSubmit()
			}
		}, w)
	})
	// 删除用户的按钮
	deleteCardTypeButton := widget.NewButton("删除卡类型", func() {
		if rowSelected < 0 {
			dialog.ShowInformation("Error", "No card type selected", w)
			return
		}
		cardType := cardTypeList[rowSelected]
		err := cardTypeService.DeleteCardType(cardType)
		if err != nil {
			dialog.ShowInformation("Error", "删除卡类型失败", w)
			return
		}
		cardTypeList = cardTypeService.GetCardType()
		cardTypeListTable.Refresh()
		dialog.ShowInformation("Success", "删除卡类型成功", w)
	})
	//搜索
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索卡类型...")
	searchEntry.OnChanged = func(s string) {
		cardTypeList = cardTypeService.GetCardTypeByName(searchEntry.Text)
		cardTypeListTable.Refresh()
	}
	top := container.NewHBox(addCardTypeButton, updateCardTypeButton, deleteCardTypeButton)

	return container.NewVBox(
		searchEntry,
		top,
		scrollContainer,
	)

}
