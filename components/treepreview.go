package components

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Single tree+preview pair
type TreePreview struct {
	*tview.Grid

	// Gui
	App       *tview.Application // ref to app for stop, focus, etc.
	Navbar    *Navbar            // navbar to adjust state
	Statusbar *Statusbar         // statusbar for debug output, confirmations, etc.
	Tree      *tview.TreeView    // list obtained from running CmdStr
	Preview   *tview.TextView    // preview window
	Ansi      io.Writer

	// Layouts
	ShowNavbar  bool
	GridColumns []int
	Focused     int
	FocusCycle  []tview.Primitive // Interface, no need for '*'

	// Config
	Config Config
	Output chan string
	Done   bool
}

func NewTreePreview(config Config, app *tview.Application) *TreePreview {
	lp := TreePreview{
		App:    app,
		Config: config,

		// Navbar, Statusbar
		Navbar:    NewNavbar(config.Path, config.Tree.Command, config.Tree.Args),
		Statusbar: NewStatusbar(config.Path),

		// Tree
		Tree: tview.NewTreeView(),

		// Preview
		Preview: tview.NewTextView(),

		// Layout
		Grid:        tview.NewGrid(),
		ShowNavbar:  true,
		GridColumns: []int{-1, -3},
	}

	// Appearance for focus
	lp.Tree.SetBorder(true)
	lp.Preview.SetBorder(true)
	lp.Preview.SetDynamicColors(true)
	lp.FocusCycle = []tview.Primitive{lp.Tree, lp.Preview}
	lp.Focused = 0

	// Hover func of list
	lp.Tree.SetChangedFunc(func(node *tview.TreeNode) {
		lp.LoadPreview()
	})

	// Selected func of list
	lp.Tree.SetSelectedFunc(func(node *tview.TreeNode) {
		lp.Focused = 1
		lp.App.SetFocus(lp.FocusCycle[lp.Focused])
	})

	// Cancel func of list
	lp.Tree.SetDoneFunc(func(tcell.Key) {
		lp.App.Stop()
	})

	// Keymaps for list
	lp.Tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	// Keymaps for preview
	lp.Preview.SetDoneFunc(func(tcell.Key) {
		lp.Focused = 0
		lp.App.SetFocus(lp.FocusCycle[lp.Focused])
	})

	// Render the grid
	lp.Render()

	// Ansi
	lp.Ansi = tview.ANSIWriter(lp.Preview)
	return &lp
}

func (lp *TreePreview) Render() {
	lp.Grid.Clear()
	if lp.ShowNavbar {
		// navbar, content
		lp.Grid.SetRows(lp.Navbar.Height, -1, 1).SetColumns(lp.GridColumns...)
		lp.Grid.AddItem(lp.Navbar, 0, 0, 1, 2, 0, 0, false)
		lp.Grid.AddItem(lp.Tree, 1, 0, 1, 1, 0, 0, true)
		lp.Grid.AddItem(lp.Preview, 1, 1, 1, 1, 0, 0, false)
		lp.Grid.AddItem(lp.Statusbar, 2, 0, 1, 2, 0, 0, false)
	} else {
		// content
		lp.Grid.SetRows(-1, 1).SetColumns(lp.GridColumns...)
		lp.Grid.AddItem(lp.Tree, 0, 0, 1, 1, 0, 0, true)
		lp.Grid.AddItem(lp.Preview, 0, 1, 1, 1, 0, 0, false)
		lp.Grid.AddItem(lp.Statusbar, 1, 0, 1, 2, 0, 0, false)
	}
	lp.FocusCycle = []tview.Primitive{lp.Tree, lp.Preview}
}

func (lp *TreePreview) ToggleNavbar() {
	lp.ShowNavbar = !lp.ShowNavbar
	lp.Render()
}

func (lp *TreePreview) Focus(delegate func(p tview.Primitive)) {
	if lp.Focused < 0 {
		delegate(lp.Box)
	} else {
		delegate(lp.FocusCycle[lp.Focused])
	}
}

func (lp *TreePreview) HasFocus() bool {
	if lp.Focused < 0 {
		return lp.Box.HasFocus()
	} else {
		return lp.FocusCycle[lp.Focused].HasFocus()
	}
}

func (lp *TreePreview) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return lp.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// Global, before handling to childPrimitive
		switch event.Rune() {
		case 'q':
			lp.App.Stop()
		}

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
		case tcell.KeyCtrlR:
			lp.LoadTree()
		case tcell.KeyEscape:
			lp.App.Stop()
		}
	})
}

func (lp *TreePreview) LoadTree() {
	cmd := exec.Command(lp.Config.Tree.Command, lp.Config.Tree.Args...)

	// Start the command
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(lp.Preview, "%s\n", "Command failed")
		return
	}

	// Scanner
	lp.Preview.Clear()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)

	root := tview.NewTreeNode("root")
	children := []*tview.TreeNode{}

	for scanner.Scan() {
		text := scanner.Text()
		children = append(children, tview.NewTreeNode(text))
		// fmt.Fprintf(lp.Ansi, "%s\n", text)
	}
	if scanner.Err() != nil {
		fmt.Fprint(lp.Ansi, fmt.Sprintf("%s\n", "Command failed"))
	}
	cmd.Wait()
	root.SetChildren(children)
	lp.Tree.SetRoot(root)
	lp.Tree.SetCurrentNode(root)

	lp.Done = true
	lp.Navbar.Loaded = true
	lp.Render()
}

func (lp *TreePreview) LoadPreview() {
	lp.Preview.Clear()

	// Inject selected
	text := lp.Tree.GetCurrentNode().GetText()
	args := []string{}
	for _, a := range lp.Config.Preview.Args {
		if a == "selected" {
			args = append(args, text)
		} else {
			args = append(args, a)
		}
	}
	cmd := exec.Command(lp.Config.Preview.Command, args...)

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
}
