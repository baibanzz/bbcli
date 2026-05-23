package windows

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	path                  string
	fyneApp               fyne.App
	config, router, creat *container.TabItem
}

func NewApp(path string) *App {
	return &App{path: path}
}

func (a *App) Run() {
	// 创建 Fyne 应用
	a.fyneApp = app.New()

	// 创建窗口
	w := a.fyneApp.NewWindow(a.path)

	b := widget.NewButton("test", func() {
		log.Printf("s")
	})
	entry := widget.NewEntry()
	entry.Resize(fyne.NewSize(400, 400))
	entry.SetText(a.path)
	entry.OnChanged = func(s string) {
		log.Printf("%s\n", s)
	}
	a.config = container.NewTabItem("配置", widget.NewLabel("项目配置面板"))
	a.router = container.NewTabItem("路由", entry)
	a.creat = container.NewTabItem("生成", b)
	//widget.NewLabel("代码生成面板"))

	// 创建标签页容器
	tabs := container.NewAppTabs(
		a.config,
		a.router,
		a.creat,
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(900, 600))
	w.ShowAndRun()
}
