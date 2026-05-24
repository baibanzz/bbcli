package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) initLeft() *fyne.Container {
	var btns []fyne.CanvasObject
	for i := 0; i < len(a.rightPanels); i++ {
		btns = append(btns, widget.NewButton(a.rightPanels[i].Title, func() {
			a.switchTab(a.rightPanels[i].Title)
		}))
	}
	// 创建左侧导航容器
	leftNav := container.NewVBox(btns...)
	return leftNav
}
