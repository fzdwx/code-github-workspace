package search

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	// searchCmd search repo
	var searchCmd = &cobra.Command{
		Use:   "search <command>",
		Short: "Search for repositories, issues, and pull requests",
		Long:  "Search across all of GitHub.",
	}

	return searchCmd
}
