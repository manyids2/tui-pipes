package components

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Navbar for gui
type Navbar struct {
	*tview.Box
	Name    string
	Command string
	Args    []string

	// Appearance
	Height    int
	Prompt    string
	Indicator string
	Loaded    bool
}

func NewNavbar(name, command string, args []string) *Navbar {
	n := Navbar{
		Box:       tview.NewBox(),
		Name:      name,
		Command:   command,
		Args:      args,
		Loaded:    false,
		Height:    2,
		Prompt:    "  ",
		Indicator: "  ",
	}
	return &n
}

func (n *Navbar) Draw(screen tcell.Screen) {
	n.Box.DrawForSubclass(screen, n)
	x, y, width, height := n.GetInnerRect()

	if height < 1 {
		return
	}
	colorLoaded := "[red]"
	if n.Loaded {
		colorLoaded = "[green]"
	}
	// Indicator
	tview.Print(screen, colorLoaded+n.Indicator, x, y, width, tview.AlignLeft, tcell.ColorWhite)

	// Command
	line := fmt.Sprintf(`%s Cmd %s[white][::b]%s`, colorLoaded, n.Prompt, n.Command)
	tview.Print(screen, line, x+3, y, width, tview.AlignLeft, tcell.ColorWhite)

	line = colorLoaded + " Args" + n.Prompt
	for index := range n.Args {
		label := n.Args[index]
		if index == 0 {
			line += fmt.Sprintf(`[white][::b]%s`, label)
		} else {
			line += fmt.Sprintf(`[gray]·[white][::b]%s`, label)
		}
	}
	tview.Print(screen, line, x+3, y+1, width, tview.AlignLeft, tcell.ColorWhite)
}
