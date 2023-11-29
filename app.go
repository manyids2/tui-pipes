package main

import (
	"github.com/manyids2/tui-pipes/components"
	"github.com/rivo/tview"
)

// View of app
type App struct {
	App         *tview.Application
	Pages       *tview.Pages
	Sidebar     *components.Sidebar
	ShowSidebar bool
}

func NewApp() *App {
	app := App{
		Pages:   tview.NewPages(),
		Sidebar: components.NewSidebar([]string{}),
	}

	// Create fullscreen app
	app.App = tview.NewApplication()

	// Create ListPreview with focus on navbar
	lp := components.NewListPreview([]string{"ls"})
	lp.SetKeymaps(app.App)

	// Add it to page and display
	app.Pages.AddPage("home", lp, true, true)
	app.App.SetRoot(app.Pages, true).SetFocus(lp.Navbar)

	return &app
}
