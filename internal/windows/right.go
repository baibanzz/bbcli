package windows

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func (a *App) initPanels() {
	// 配置面板
	a.rightPanels.Push("配置", a.initConfig())

	// 路由面板
	entry := widget.NewEntry()
	entry.Resize(fyne.NewSize(400, 400))
	entry.SetText(a.path)
	entry.OnChanged = func(s string) {
		log.Printf("%s\n", s)
	}
	a.rightPanels.Push("路由", entry)

	// 生成面板
	btn := widget.NewButton("生成代码", func() {
		log.Printf("开始生成代码...")
	})
	a.rightPanels.Push("生成", btn)
	a.rightPanels.Push("关于", btn)
}

func (a *App) initConfig() fyne.CanvasObject {
	return widget.NewLabel("项目配置面板")
}
