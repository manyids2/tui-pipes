package cmd

import (
	"log"

	"github.com/manyids2/tui-pipes/components"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// tpCmd represents the tp command
var tpCmd = &cobra.Command{
	Use:   "tp",
	Short: "tree preview",
	Long:  `tree preview`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := components.ReadConfig("./configs/find_bat.json")
		if err != nil {
			log.Fatal(err)
		}

		// Create ListPreview
		app := tview.NewApplication()
		tp := components.NewTreePreview(config, app)
		tp.LoadTree()
		app.SetRoot(tp, true).SetFocus(tp)

		// Run app
		err = app.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(tpCmd)
}
