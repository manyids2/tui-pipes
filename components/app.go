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

func NewApp(config Config) *App {
	app := App{
		Pages:   tview.NewPages(),
		Sidebar: NewSidebar([]string{}),
		App:     tview.NewApplication(),
	}

	// // Create ListPreview
	// lp := NewListPreview(config, app.App)
	// lp.LoadList()

	// Create ListPreview
	tp := NewTreePreview(config, app.App)
	tp.LoadTree()

	// Add it to page and display
	app.Pages.AddPage("home", tp, true, true)
	app.Pages.SwitchToPage("home")
	app.App.SetRoot(app.Pages, true).SetFocus(tp)

	return &app
}
