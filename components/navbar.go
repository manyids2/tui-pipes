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
	Output  chan string

	// Appearance
	Height    int
	Prompt    string
	Indicator string

	// Focus : 0 for command, 1+ for args
	Focused    int
	FocusCycle []string

	// Status
	Loaded bool

	// An optional function which is called when the user hits Escape.
	cancel func()
}

func NewNavbar(name, command string, args []string) *Navbar {
	n := Navbar{
		Box:        tview.NewBox(),
		Name:       name,
		Command:    command,
		Args:       args,
		Loaded:     false,
		Height:     5,
		Prompt:     "  ",
		Indicator:  "  ",
		Focused:    0,
		FocusCycle: []string{"command", "args"},
	}
	n.Box.SetBorder(true)
	return &n
}

func (n *Navbar) Draw(screen tcell.Screen) {
	n.Box.DrawForSubclass(screen, n)
	x, y, width, height := n.GetInnerRect()
	// debugInfo := fmt.Sprintf("%d, %d, %d, %d", x, y, width, height)

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
	var focused string
	if n.FocusCycle[n.Focused] == "command" {
		focused = "[::b]"
	}
	line := fmt.Sprintf(`%s Cmd %s[white][::b]%s`, colorLoaded, n.Prompt, n.Command)
	tview.Print(screen, focused+line, x+3, y, width, tview.AlignLeft, tcell.ColorWhite)
	focused = ""

	// Args
	if n.FocusCycle[n.Focused] == "args" {
		focused = "[::b]"
	}
	line = colorLoaded + focused + " Args" + n.Prompt
	for index := range n.Args {
		label := n.Args[index]
		if index == 0 {
			line += fmt.Sprintf(`[white][::b]%s`, label)
		} else {
			line += fmt.Sprintf(`%s · [white][::b]%s`, colorLoaded, label)
		}
	}
	tview.Print(screen, line, x+3, y+1, width, tview.AlignLeft, tcell.ColorWhite)
	focused = ""

	// Focused
	tview.Print(screen, fmt.Sprintf("%d", n.Focused), x+3, y+2, width, tview.AlignLeft, tcell.ColorWhite)
}

func (n *Navbar) SetCancelFunc(callback func()) *Navbar {
	n.cancel = callback
	return n
}

func (n *Navbar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return n.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyEscape, tcell.KeyTAB:
			if n.cancel != nil {
				n.cancel()
			} else {
				n.Focused = 0
				setFocus(n)
			}
		case tcell.KeyDown:
			n.Focused = (n.Focused + 1) % (1 + len(n.FocusCycle))
		case tcell.KeyUp:
			n.Focused = n.Focused - 1
			if n.Focused < 0 {
				n.Focused = len(n.FocusCycle)
			}
		}
	})
}
