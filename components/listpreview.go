package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Define transform using go for convinience over bash?
type transform_fn func(string) string

// Single list+preview pair
type ListPreview struct {
	*tview.Grid

	// Gui
	G *tview.Grid     // root of this ListPreview
	N *Navbar         // navbar to adjust state
	L *tview.List     // list obtained from running CmdStr
	T *tview.TextView // preview window

	// Purity
	ShowNavbar bool

	// State
	CmdArgs         []string     // args for cli command
	Separator       string       // lines, words
	TransformView   transform_fn // function for display
	TransformSelect transform_fn // function for selection
}

func NewListPreview(cmdArgs []string) *ListPreview {
	lp := ListPreview{
		// Default ui
		G: tview.NewGrid(),
		N: NewNavbar(cmdArgs),
		L: tview.NewList(),
		T: tview.NewTextView().SetText("TextView"),
		// Purity
		ShowNavbar: true,
		// State
		CmdArgs:         cmdArgs,
		Separator:       "lines",
		TransformView:   func(s string) string { return s },
		TransformSelect: func(s string) string { return s },
	}
	lp.L.SetBorder(true)
	lp.T.SetBorder(true)
	lp.Render()
	return &lp
}

func (lp *ListPreview) Render() {
	lp.G.Clear()

	if lp.ShowNavbar {
		// navbar, content
		lp.G.SetRows(3, -1).SetColumns(-1, -2)
		lp.G.AddItem(lp.N, 0, 0, 1, 2, 0, 0, false)
		lp.G.AddItem(lp.L, 1, 0, 1, 1, 0, 0, true)
		lp.G.AddItem(lp.T, 1, 1, 1, 1, 0, 0, false)
	} else {
		// content
		lp.G.SetRows(-1).SetColumns(-1, -2)
		lp.G.AddItem(lp.L, 0, 0, 1, 1, 0, 0, true)
		lp.G.AddItem(lp.T, 0, 1, 1, 1, 0, 0, false)
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

		case tcell.KeyCtrlSpace:
			lp.ToggleNavbar()
			return nil
		}
		return event
	})
}
