package main

import (
	"fyne.io/fyne/v2/app"
	"os"
	"sql_bank/initialize"
	"sql_bank/page"
)

func main() {
	initialize.InitConfig()
	initialize.InitDB()
	//恢复数据库未处理的操作
	initialize.RecoverData()

	//设置中文
	err := os.Setenv("FYNE_FONT", "Front/STFANGSO.TTF")
	if err != nil {
		return
	}
	a := app.New()
	page.MakeLoginUI(a)
	a.Run()
}
