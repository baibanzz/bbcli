package windows

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/baibanzz/bbcli/internal/core"
)

func (a *App) initPanels() {
	// 配置面板
	a.rightPanels.Push(&core.Tabs{
		Object: widget.NewLabel("项目配置面板"),
		Title:  "配置",
	})

	// 路由面板
	entry := widget.NewEntry()
	entry.Resize(fyne.NewSize(400, 400))
	entry.SetText(a.path)
	entry.OnChanged = func(s string) {
		log.Printf("%s\n", s)
	}
	a.rightPanels.Push(&core.Tabs{
		Object: entry,
		Title:  "路由",
	})

	// 生成面板
	btn := widget.NewButton("生成代码", func() {
		log.Printf("开始生成代码...")
	})
	a.rightPanels.Push(&core.Tabs{
		Object: btn,
		Title:  "生成",
	})
}
