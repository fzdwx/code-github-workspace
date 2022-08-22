package list

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fzdwx/code-github-workspace/gh"
	"github.com/spf13/cobra"
)

type Model struct {
	Keymap *Keymap
}

func New() *Model {
	return &Model{
		Keymap: DefaultKeyMap(),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Quit):
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m *Model) View() string {
	return "hello"
}

func test() {
	ctx := context.Background()
	// list all repositories for the authenticated user
	repos, _, err := gh.GetAuth(ctx).Repositories.List(ctx, "", nil)

	cobra.CheckErr(err)

	for _, repo := range repos {
		fmt.Println(repo)
	}
	fmt.Println(repos)
}
