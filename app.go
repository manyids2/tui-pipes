package main

import (
	"github.com/manyids2/tui-pipes/components"
	"github.com/rivo/tview"
)

// View of app
type App struct {
	A           *tview.Application
	P           *tview.Pages
	S           *components.Sidebar // sidebar to choose from various ListPreviews
	ShowSidebar bool
}

func NewApp() *App {
	app := App{
		P: tview.NewPages(),
		S: components.NewSidebar([]string{}),
	}

	// Create fullscreen app
	app.A = tview.NewApplication()

	// Create ListPreview with focus on navbar
	listpreview := components.NewListPreview([]string{"ls"})
	listpreview.SetKeymaps(app.A)

	// Add it to page and display
	app.P.AddPage("home", listpreview.G, true, true)
	app.A.SetRoot(app.P, true).SetFocus(app.P)

	return &app
}
