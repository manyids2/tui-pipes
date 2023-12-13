package cmd

import (
	"log"

	"github.com/manyids2/tui-pipes/components"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// tilesCmd represents the tiles command
var tilesCmd = &cobra.Command{
	Use:   "tiles",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := components.ReadConfig("./configs/find_bat.json")
		if err != nil {
			log.Fatal(err)
		}

		// Create ListTiles
		app := tview.NewApplication()
		lt := components.NewListTiles(config, app)
		lt.LoadList()
		app.SetRoot(lt, true).SetFocus(lt)

		// Run app
		err = app.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tilesCmd)
}
