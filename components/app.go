package components

import (
	"github.com/rivo/tview"
)

// View of app
type App struct {
	App         *tview.Application
	Pages       *tview.Pages
	Sidebar     *Sidebar
	ShowSidebar bool
}

func NewApp() *App {
	app := App{
		Pages:   tview.NewPages(),
		Sidebar: NewSidebar([]string{}),
		App:     tview.NewApplication(),
	}

	// Create ListPreview with focus on navbar
	lp := NewListPreview("List files",
		"exa", []string{".", "-T", "--icons", "--color=always"},
		app.App,
	)

	// Add it to page and display
	app.Pages.AddPage("home", lp, true, true)
	app.App.SetRoot(app.Pages, true).SetFocus(lp.Navbar)

	return &app
}
