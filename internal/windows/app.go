package windows

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/baibanzz/bbcli/internal/core"
)

type App struct {
	path        string
	fyneApp     fyne.App
	w           fyne.Window
	rightPanel  *fyne.Container
	rightPanels core.Switchs
	split       *container.Split
	//configPanel fyne.CanvasObject
	//routerPanel fyne.CanvasObject
	//createPanel fyne.CanvasObject
}

func NewApp(path string) *App {
	return &App{path: path, rightPanels: core.NewSwitchs()}
}

func (a *App) Run() {
	// 创建 Fyne 应用
	a.fyneApp = app.New()

	// 创建窗口
	a.w = a.fyneApp.NewWindow(a.path)

	// 创建右侧内容面板
	a.initRight()

	// 创建左右分栏布局（直接使用配置内容，不用 Container 包装）
	a.split = container.NewHSplit(a.initLeft(), a.rightPanels.Title("配置"))
	a.split.Offset = 0.2

	a.w.SetContent(a.split)
	a.w.Resize(fyne.NewSize(900, 600))
	a.w.ShowAndRun()
}

func (a *App) switchTab(uid string) {
	log.Println(uid)
	// 重新创建分栏布局
	//split := container.NewHSplit(a.initLeft(), a.rightPanels.Title(uid))
	a.split.Trailing = a.rightPanels.Title(uid)
	//split.Offset = 0.2
	a.w.SetContent(a.split)
}
