package list

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fzdwx/gh-sp/model/repolist"
	"github.com/google/go-github/v46/github"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	// repoListCmd represents the list command
	var repoListCmd = &cobra.Command{
		Use:   "list",
		Short: `List repositories owned by user or organization`,
		Run: func(cmd *cobra.Command, args []string) {
			ops := &github.RepositoryListOptions{
				Sort:        "updated",
				Direction:   "desc",
				ListOptions: github.ListOptions{},
			}

			if err := tea.NewProgram(repolist.New(ops), tea.WithAltScreen()).Start(); err != nil {
				panic(err)
			}
		},
	}

	return repoListCmd
}
