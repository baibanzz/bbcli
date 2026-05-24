package core

import "fyne.io/fyne/v2"

type Tabs struct {
	Object fyne.CanvasObject
	Title  string
}
type Switchs map[int]*Tabs

func NewSwitchs() Switchs {
	return Switchs{}
}

func (sws Switchs) Push(title string, Object fyne.CanvasObject) {
	sws[len(sws)] = &Tabs{
		Object: Object,
		Title:  title,
	}
}

func (sws Switchs) Title(s string) fyne.CanvasObject {
	for _, t := range sws {
		if t.Title == s {
			return t.Object
		}
	}
	return nil
}
