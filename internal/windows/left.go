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
			a.switchTab(v.Title)
		}))
	}
	// 创建左侧导航容器
	leftNav := container.NewVBox(btns...)
	return leftNav
}
