package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Navbar struct {
	*tview.Box
	Labels  []string
	Current int
}

func NewNavbar(labels []string) *Navbar {
	return &Navbar{
		Box:    tview.NewBox(),
		Labels: labels,
	}
}

func (n *Navbar) Draw(screen tcell.Screen) {
	n.Box.DrawForSubclass(screen, n)
	x, y, width, height := n.GetInnerRect()
	// debugInfo := fmt.Sprintf("%d, %d, %d, %d", x, y, width, height)

	if height < 1 {
		return
	}
	selected := "  "
	notSelected := "  "
	line := ""
	for index := range n.Labels {
		label := n.Labels[index]
		if index == n.Current {
			line += fmt.Sprintf(`[red]%s[::b]%s`, selected, label)
		} else {
			line += fmt.Sprintf(`[white]%s%s`, notSelected, label)
		}
	}
	tview.Print(screen, line, x+3, y, width, tview.AlignLeft, tcell.ColorWhite)
}

func (n *Navbar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return n.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyLeft:
			n.Current--
			if n.Current < 0 {
				n.Current = len(n.Labels) - 1
			}
		case tcell.KeyRight:
			n.Current++
			if n.Current >= len(n.Labels) {
				n.Current = 0
			}
		}
	})
}
