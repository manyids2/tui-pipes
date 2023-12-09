package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Statusbar for gui
type Statusbar struct {
	*tview.Box
	Text string

	// Appearance
	Height    int
	Indicator string
}

func NewStatusbar(text string) *Statusbar {
	s := Statusbar{
		Box:       tview.NewBox(),
		Text:      text,
		Height:    1,
		Indicator: "  ",
	}
	s.Box.SetBorder(true)
	return &s
}

func (s *Statusbar) Draw(screen tcell.Screen) {
	s.Box.DrawForSubclass(screen, s)
	x, y, width, height := s.GetInnerRect()
	if height < 1 {
		return
	}
	tview.Print(screen, fmt.Sprintf("%s %s", s.Indicator, s.Text), x+3, y, width, tview.AlignLeft, tcell.ColorWhite)
}
