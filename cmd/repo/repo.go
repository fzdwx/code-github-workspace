package repo

import (
	"github.com/fzdwx/gh-sp/cmd/repo/list"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	// repoCmd represents the list command
	var repoCmd = &cobra.Command{
		Use:   "repo <command>",
		Short: "About github repos api",
	}

	repoCmd.AddCommand(list.New())

	return repoCmd
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// repoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// repoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
