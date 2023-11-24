package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Sidebar struct {
	*tview.Box
	Labels   []string
	Listview *tview.List
	Current  int
}

func NewSidebar(labels []string) *Sidebar {
	return &Sidebar{
		Box:    tview.NewBox().SetBorder(true),
		Labels: labels,
	}
}

func (s *Sidebar) Draw(screen tcell.Screen) {
	s.Box.DrawForSubclass(screen, s)
	x, y, width, height := s.GetInnerRect()

	selected := "  "
	notSelected := "  "
	for index := range s.Labels {
		if y+index > height {
			return
		}
		line := ""
		label := s.Labels[index]
		if index == s.Current {
			line += fmt.Sprintf(`[red]%s[::b]%s`, selected, label)
		} else {
			line += fmt.Sprintf(`[white]%s%s`, notSelected, label)
		}
		line += "\n"
		tview.Print(screen, line, x, y+index, width, tview.AlignLeft, tcell.ColorWhite)
	}
}

func (s *Sidebar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return s.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyLeft:
			s.Current--
			if s.Current < 0 {
				s.Current = len(s.Labels) - 1
			}
		case tcell.KeyRight:
			s.Current++
			if s.Current >= len(s.Labels) {
				s.Current = 0
			}
		}
	})
}
