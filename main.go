package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/jasonlvhit/gocron"
	"os"
	"sql_bank/initialize"
	"sql_bank/page/common"
	"sql_bank/task"
)

func main() {
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitRedis()
	//恢复数据库未处理的操作
	initialize.RecoverData()
	//定时任务
	s := gocron.NewScheduler()
	err := s.Every(1).Wednesday().Do(task.RepayOverDraftsTask)
	if err != nil {
		return
	}
	s.Start() // 启动调度器
	//设置中文
	err = os.Setenv("FYNE_FONT", "Front/STFANGSO.TTF")
	if err != nil {
		return
	}
	a := app.New()
	common.MakeLoginUI(a)
	a.Run()
}
