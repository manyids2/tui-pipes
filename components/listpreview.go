package components

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Define transform using go for convinience over bash?
type transform_fn func(string) string

// Single list+preview pair
type ListPreviewOptions struct {
	Separator       string       // lines, words, delimiter, etc.
	TransformView   transform_fn // function for display
	TransformSelect transform_fn // function for selection
}

// Single list+preview pair
type ListPreview struct {
	*tview.Grid

	// Gui
	App       *tview.Application // ref to app for stop, focus, etc.
	Navbar    *Navbar            // navbar to adjust state
	Statusbar *Statusbar         // statusbar for debug output, confirmations, etc.
	List      *tview.List        // list obtained from running CmdStr
	Preview   *tview.TextView    // preview window
	Ansi      io.Writer

	// Layouts
	ShowNavbar  bool
	GridColumns []int
	Focused     int
	FocusCycle  []tview.Primitive // Interface, no need for '*'

	// Arguments for cli command
	Command string
	Args    []string
	Limit   int
	Output  chan string
	Loaded  bool

	// Options
	Options ListPreviewOptions
}

func NewListPreview(name, command string, args []string, app *tview.Application) *ListPreview {
	lp := ListPreview{
		App: app,

		// Navbar
		Command: command,
		Args:    args,
		Limit:   128,

		// Navbar, Statusbar
		Navbar:    NewNavbar(name, command, args),
		Statusbar: NewStatusbar(name),

		// List
		List: tview.NewList(),

		// Preview
		Preview: tview.NewTextView(),

		// Layout
		Grid:        tview.NewGrid(),
		ShowNavbar:  true,
		GridColumns: []int{-1, -3},
	}

	// Appearance for focus
	lp.List.SetBorder(true)
	lp.Preview.SetBorder(true)
	lp.Preview.SetDynamicColors(true)
	lp.FocusCycle = []tview.Primitive{lp.Navbar, lp.List, lp.Preview}

	// Cancel func of navbar
	lp.Navbar.SetCancelFunc(func() {
		lp.Focused = -1
		app.SetFocus(lp.Box)
	})

	// Cancel func of list
	lp.List.SetDoneFunc(func() {
		lp.Focused = -1
		app.SetFocus(lp.Box)
	})

	// Cancel func of TextView
	lp.Preview.SetDoneFunc(func(tcell.Key) {
		lp.Focused = -1
		app.SetFocus(lp.Box)
	})

	// Render the grid
	lp.Render()

	// Ansi
	lp.Ansi = tview.ANSIWriter(lp.Preview)
	return &lp
}

func (lp *ListPreview) Render() {
	lp.Grid.Clear()

	if lp.ShowNavbar {
		// navbar, content
		lp.Grid.SetRows(lp.Navbar.Height, -1, 1).SetColumns(lp.GridColumns...)
		lp.Grid.AddItem(lp.Navbar, 0, 0, 1, 2, 0, 0, true)
		lp.Grid.AddItem(lp.List, 1, 0, 1, 1, 0, 0, false)
		lp.Grid.AddItem(lp.Preview, 1, 1, 1, 1, 0, 0, false)
		lp.Grid.AddItem(lp.Statusbar, 2, 0, 1, 2, 0, 0, false)
		lp.FocusCycle = []tview.Primitive{lp.Navbar, lp.List, lp.Preview}
	} else {
		// content
		lp.Grid.SetRows(-1, 1).SetColumns(lp.GridColumns...)
		lp.Grid.AddItem(lp.List, 0, 0, 1, 1, 0, 0, true)
		lp.Grid.AddItem(lp.Preview, 0, 1, 1, 1, 0, 0, false)
		lp.Grid.AddItem(lp.Statusbar, 1, 0, 1, 2, 0, 0, false)
		if lp.FocusCycle[lp.Focused] == lp.Navbar {
			lp.Focused = 0
		} else {
			lp.Focused = lp.Focused - 1 // Navbar was first element
		}
		lp.FocusCycle = []tview.Primitive{lp.List, lp.Preview}
	}
}

func (lp *ListPreview) ToggleNavbar() {
	lp.ShowNavbar = !lp.ShowNavbar
	if lp.FocusCycle[lp.Focused] == lp.Navbar {
		lp.Focused = 0
	}
	lp.Render()
}

func (lp *ListPreview) Focus(delegate func(p tview.Primitive)) {
	if lp.Focused < 0 {
		delegate(lp.Box)
	} else {
		delegate(lp.FocusCycle[lp.Focused])
	}
}

func (lp *ListPreview) HasFocus() bool {
	if lp.Focused < 0 {
		return lp.Box.HasFocus()
	} else {
		return lp.FocusCycle[lp.Focused].HasFocus()
	}
}

func (lp *ListPreview) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return lp.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if lp.Focused >= 0 {
			childPrimitive := lp.FocusCycle[lp.Focused]
			if childPrimitive.HasFocus() {
				if handler := childPrimitive.InputHandler(); handler != nil {
					handler(event, setFocus)
				}
			}
			return
		}

		// Grid level shortcuts
		switch event.Key() {
		case tcell.KeyTAB:
			lp.Focused = (lp.Focused + 1) % len(lp.FocusCycle)
			setFocus(lp.FocusCycle[lp.Focused])
		case tcell.KeyCtrlSpace:
			lp.ToggleNavbar()
		case tcell.KeyCtrlT:
			lp.Focused = -1
			setFocus(lp.Box)
		case tcell.KeyEnter:
			lp.Run()
		case tcell.KeyEscape:
			lp.App.Stop()
		}
		switch event.Rune() {
		case 'q':
			lp.App.Stop()
		}
	})
}

func (lp *ListPreview) Run() {
	cmd := exec.Command(lp.Command, lp.Args...)

	// Start the command
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(lp.Preview, "%s\n", "Command failed")
		return
	}

	// Scanner
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintf(lp.Ansi, "%s\n", text)
	}
	if scanner.Err() != nil {
		fmt.Fprint(lp.Ansi, fmt.Sprintf("%s\n", "Command failed"))
	}
	cmd.Wait()

	lp.Loaded = true
	lp.Navbar.Loaded = true
	lp.Render()
}
