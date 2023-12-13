package cmd

import (
	"log"

	"github.com/manyids2/tui-pipes/components"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// lpCmd represents the lp command
var lpCmd = &cobra.Command{
	Use:   "lp",
	Short: "list preview",
	Long:  `list preview`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := components.ReadConfig("./configs/find_bat.json")
		if err != nil {
			log.Fatal(err)
		}

		// Create ListPreview
		app := tview.NewApplication()
		lp := components.NewListPreview(config, app)
		lp.LoadList()
		app.SetRoot(lp, true).SetFocus(lp)

		// Run app
		err = app.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lpCmd)
}
