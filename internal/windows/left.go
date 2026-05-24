package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) initLeft() *fyne.Container {
	var btns []fyne.CanvasObject
	for _, v := range a.rightPanels {
		btns = append(btns, widget.NewButton(v.Title, func() {
			a.switchTab(v.Title, a.w)
		}))
	}
	// 创建左侧导航按钮
	//btnConfig := widget.NewButton("配置", func() { a.switchTab("配置", a.w) })
	//btnRouter := widget.NewButton("路由", func() { a.switchTab("路由", a.w) })
	//btnCreate := widget.NewButton("生成", func() { a.switchTab("生成", a.w) })

	// 创建左侧导航容器
	leftNav := container.NewVBox(btns...)
	//btns[0].Importance = widget.HighImportance
	return leftNav
}
