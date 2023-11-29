package components

import (
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
	Navbar  *Navbar         // navbar to adjust state
	List    *tview.List     // list obtained from running CmdStr
	Preview *tview.TextView // preview window

	// Layouts
	ShowNavbar  bool
	GridColumns []int
	Focused     int
	FocusCycle  []tview.Primitive // Interface, no need for '*'

	// Arguments for cli command
	CmdArgs []string
	Limit   int

	// Options
	Options ListPreviewOptions
}

func NewListPreview(cmdArgs []string) *ListPreview {
	lp := ListPreview{
		// Navbar
		Navbar: NewNavbar(cmdArgs),
		// List
		List: tview.NewList(),
		// Preview
		Preview: tview.NewTextView(),
		// Layout
		Grid:        tview.NewGrid(),
		ShowNavbar:  true,
		GridColumns: []int{-1, -1},
		// Inputs
		CmdArgs: cmdArgs,
	}
	lp.List.SetBorder(true)
	lp.Preview.SetBorder(true)
	lp.FocusCycle = []tview.Primitive{lp.Navbar, lp.List, lp.Preview}
	lp.Render()
	return &lp
}

func (lp *ListPreview) Render() {
	lp.Grid.Clear()

	if lp.ShowNavbar {
		// navbar, content
		lp.Grid.SetRows(3, -1).SetColumns(lp.GridColumns...)
		lp.Grid.AddItem(lp.Navbar, 0, 0, 1, 2, 0, 0, true)
		lp.Grid.AddItem(lp.List, 1, 0, 1, 1, 0, 0, false)
		lp.Grid.AddItem(lp.Preview, 1, 1, 1, 1, 0, 0, false)
		lp.FocusCycle = []tview.Primitive{lp.Navbar, lp.List, lp.Preview}
	} else {
		// content
		lp.Grid.SetRows(-1).SetColumns(lp.GridColumns...)
		lp.Grid.AddItem(lp.List, 0, 0, 1, 1, 0, 0, true)
		lp.Grid.AddItem(lp.Preview, 0, 1, 1, 1, 0, 0, false)
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
	lp.Render()
}

func (lp *ListPreview) SetKeymaps(A *tview.Application) {
	A.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {

		case tcell.KeyEsc:
			A.Stop()
			return nil

		case tcell.KeyTAB:
			lp.Focused = (lp.Focused + 1) % len(lp.FocusCycle)
			A.SetFocus(lp.FocusCycle[lp.Focused])
			return nil

		case tcell.KeyCtrlSpace:
			lp.ToggleNavbar()
			A.SetFocus(lp.FocusCycle[lp.Focused])
			return nil
		}
		return event
	})
}
