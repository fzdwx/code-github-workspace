package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fzdwx/code-github-workspace/model/list"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show Github repo list",
	Run: func(cmd *cobra.Command, args []string) {

		if err := tea.NewProgram(list.New(), tea.WithAltScreen()).Start(); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
