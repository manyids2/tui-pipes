package components

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Single list+tiles pair
type ListTiles struct {
	*tview.Grid

	// Gui
	App   *tview.Application // ref to app for stop, focus, etc.
	List  *tview.List        // list window
	Tiles []*tview.TextView  // preview window
	Ansis []io.Writer

	// Layouts
	GridColumns []int
	GridRows    []int
	Focused     int

	// Config
	Config Config
}

func NewListTiles(config Config, app *tview.Application) *ListTiles {
	lt := ListTiles{
		App:         app,
		Config:      config,
		List:        tview.NewList(),
		Tiles:       []*tview.TextView{},
		Grid:        tview.NewGrid(),
		GridColumns: []int{-1},
		GridRows:    []int{-1},
	}
	lt.InitTiles(lt.GridRows, lt.GridColumns)
	return &lt
}

func (lt *ListTiles) InitTiles(GridRows, GridColumns []int) {
	// Initialize base tile
	lt.GridRows = GridRows
	lt.GridColumns = GridColumns
	if (len(lt.GridColumns) <= 0) || (len(lt.GridRows) <= 0) {
		log.Fatalln("Bad params:", lt.GridRows, lt.GridColumns)
	}

	// Reset
	lt.Tiles = []*tview.TextView{}
	lt.Ansis = []io.Writer{}
	n_tiles := len(lt.GridRows) * len(lt.GridColumns)
	for i := 0; i < n_tiles; i++ {
		tv := tview.NewTextView()
		tv.SetBorder(true)
		tv.SetDoneFunc(func(key tcell.Key) {
			lt.App.Stop()
		})
		lt.Tiles = append(lt.Tiles, tv)
		lt.Ansis = append(lt.Ansis, tview.ANSIWriter(tv))
	}

	lt.Render()
}

func (lt *ListTiles) Render() {
	lt.Grid.Clear()
	GridColumns := append(lt.GridColumns, -1)
	lt.Grid.SetRows(lt.GridRows...).SetColumns(GridColumns...)
	for i := 0; i < len(lt.GridRows); i++ {
		for j := 0; j < len(lt.GridColumns); j++ {
			lt.Grid.AddItem(lt.Tiles[i*len(lt.GridColumns)+j], i, j, 1, 1, 0, 0, false)
		}
	}
}

func (lt *ListTiles) Focus(delegate func(p tview.Primitive)) {
	if lt.Focused < 0 {
		delegate(lt.List)
	} else {
		delegate(lt.Tiles[lt.Focused])
	}
}

func (lt *ListTiles) HasFocus() bool {
	if lt.Focused < 0 {
		return lt.List.HasFocus()
	} else {
		return lt.Tiles[lt.Focused].HasFocus()
	}
}

func (lt *ListTiles) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return lt.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		// Global, before handling to childPrimitive
		switch event.Rune() {
		case 'q':
			lt.App.Stop()
		}

		if lt.Focused >= 0 {
			childPrimitive := lt.Tiles[lt.Focused]
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
			lt.Focused = (lt.Focused + 1) % len(lt.Tiles)
			setFocus(lt.Tiles[lt.Focused])
		case tcell.KeyCtrlR:
			lt.LoadList()
		case tcell.KeyCtrlL:
			lt.Focused = -1
			setFocus(lt.List)
		case tcell.KeyEscape:
			lt.App.Stop()
		}
	})
}

func (lt *ListTiles) LoadList() {
	cmd := exec.Command(lt.Config.List.Command, lt.Config.List.Args...)

	// Start the command
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		fmt.Fprintf(lt.Tiles[0], "%s: %s\n", "Command failed", err)
		return
	}

	// Scanner
	lt.Tiles[0].Clear()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := scanner.Text()
		lt.List.AddItem(text, "", 0, nil)
		// fmt.Fprintf(lt.Ansis[0], "%s\n", text)
	}
	if scanner.Err() != nil {
		fmt.Fprint(lt.Ansis[0], fmt.Sprintf("%s\n", "Command failed"))
	}
	cmd.Wait()

	// lt.Render()
}

func (lt *ListTiles) LoadTile(index int) {
}
