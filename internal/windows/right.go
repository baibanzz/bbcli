package windows

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/baibanzz/bbcli/internal/global"
)

func (a *App) initRight() {
	a.rightPanels.Push("配置", a.initConfig())
	a.rightPanels.Push("路由", a.initRouter())
	a.rightPanels.Push("生成", a.initCreate())
	a.rightPanels.Push("关于", a.initAbout())
}

func (a *App) initConfig() fyne.CanvasObject {
	// 项目名称
	nameEntry := widget.NewEntry()
	nameEntry.SetText(global.Config.Name)
	nameEntry.OnChanged = func(s string) {
		global.Config.Name = s
	}

	// 数据库选项
	checkMySQL := widget.NewCheck("MySQL", func(b bool) {
		global.Config.UseMySQL = b
	})
	checkMySQL.Checked = global.Config.UseMySQL

	checkPostgreSQL := widget.NewCheck("PostgreSQL", func(b bool) {
		global.Config.UsePostgreSQL = b
	})
	checkPostgreSQL.Checked = global.Config.UsePostgreSQL

	checkSQLite3 := widget.NewCheck("SQLite3", func(b bool) {
		global.Config.UseSQLite3 = b
	})
	checkSQLite3.Checked = global.Config.UseSQLite3

	checkRedis := widget.NewCheck("Redis", func(b bool) {
		global.Config.UseRedis = b
	})
	checkRedis.Checked = global.Config.UseRedis

	// 中间件选项
	checkJWT := widget.NewCheck("JWT", func(b bool) {
		global.Config.UseJwt = b
	})
	checkJWT.Checked = global.Config.UseJwt

	checkCore := widget.NewCheck("跨域", func(b bool) {
		global.Config.UseCore = b
	})
	checkCore.Checked = global.Config.UseCore

	// 扩展组件
	checkEtcd := widget.NewCheck("Etcd", func(b bool) {
		global.Config.UseEtcd = b
	})
	checkEtcd.Checked = global.Config.UseEtcd

	vbox := container.NewVBox(
		container.NewBorder(nil, nil, widget.NewLabel("项目名称:"), nil, nameEntry),
		widget.NewSeparator(),
		widget.NewLabel("数据库"),
		container.NewHBox(checkMySQL, checkPostgreSQL, checkSQLite3, checkRedis),
		widget.NewSeparator(),
		widget.NewLabel("中间件"),
		container.NewHBox(checkJWT, checkCore),
		widget.NewSeparator(),
		widget.NewLabel("扩展组件"),
		checkEtcd,
	)
	top := widget.NewLabel("项目配置")
	top.TextStyle = fyne.TextStyle{Bold: true, TabWidth: 20}
	// 添加内边距
	pad := widget.NewLabel("")
	paddedContent := container.NewBorder(top, nil, pad, pad, vbox)
	scroll := container.NewScroll(paddedContent)
	scroll.SetMinSize(fyne.NewSize(400, 500))
	return scroll
}

func (a *App) initRouter() fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.SetText(a.path)
	entry.OnChanged = func(s string) {
		log.Printf("%s\n", s)
	}
	return entry
}

func (a *App) initCreate() fyne.CanvasObject {
	return widget.NewButton("生成代码", func() {
		log.Printf("开始生成代码...")
	})
}

func (a *App) initAbout() fyne.CanvasObject {
	return widget.NewLabel("关于面板")
}
