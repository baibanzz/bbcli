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
	a.initPanels()

	// 创建右侧内容容器
	a.rightPanel = fyne.NewContainer(a.rightPanels[0].Object)

	// 创建左侧菜单
	leftNav := a.initLeft()
	// 创建左右分栏布局
	split := container.NewHSplit(leftNav, a.rightPanel)
	split.Offset = 0.2

	a.w.SetContent(split)
	a.w.Resize(fyne.NewSize(900, 600))
	a.w.ShowAndRun()
}

func (a *App) switchTab(uid string, w fyne.Window) {
	//var content fyne.CanvasObject
	//switch uid {
	//case "配置":
	//	content = a.configPanel
	//case "路由":
	//	content = a.routerPanel
	//case "生成":
	//	content = a.createPanel
	//default:
	//	content = a.configPanel
	//}
	log.Println(uid)
	// 替换右侧内容
	a.rightPanel.Objects = []fyne.CanvasObject{a.rightPanels.Title(uid)}
}
