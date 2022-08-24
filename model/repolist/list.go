package repolist

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fzdwx/gh-sp/api"
	"github.com/fzdwx/gh-sp/model/table"
	"github.com/google/go-github/v46/github"
	"github.com/spf13/cobra"
)

type status int

const (
	loading status = iota
	loaded
)

type Model struct {
	Keymap *Keymap
	Width  int
	Height int
	ops    *github.RepositoryListOptions

	spinner spinner.Model

	status status

	table  *table.Model
	cancel context.CancelFunc
}

func headers() table.Headers {
	return []*table.Header{
		{Text: "repo name", Ratio: 5},
		{Text: "description", Ratio: 12, MinWidth: 20},
		{Text: "start", Ratio: 2, MinWidth: 5},
		{Text: "\uF707", Ratio: 2, MinWidth: 7},
		{Text: "issues", Ratio: 2, MinWidth: 10},
	}
}

func New(ops *github.RepositoryListOptions) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &Model{
		Keymap:  DefaultKeyMap(),
		spinner: s,
		ops:     ops,
		table:   table.NewModel(headers()),
	}
}

func (m *Model) Init() tea.Cmd {
	m.status = loading
	go func() {
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancel = cancelFunc
		repos, _, err := api.Get(ctx).Repositories.List(ctx, "", m.ops)

		cobra.CheckErr(err)

		m.addRows(repos)
		m.status = loaded
	}()
	return spinner.Tick
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Quit):
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.table.SetWidth(m.Width)
		m.table.SetHeight(m.Height)
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	if m.status == loading {
		return m.spinnerView()
	}
	return m.table.View()
}

func (m *Model) spinnerView() string {
	return lipgloss.NewStyle().PaddingTop(m.Height / 3).Width(m.Width).Align(lipgloss.Center).Render(m.spinner.View() + " 正在加载 repository ...")
}

func (m *Model) addRows(repos []*github.Repository) {
	for _, repo := range repos {
		m.table.AppendRow(table.Row{
			repo.GetFullName(),
			repo.GetDescription(),
			fmt.Sprintf("🌟 %d", repo.GetStargazersCount()),
			repo.GetVisibility(),
			fmt.Sprintf("🎯 %d", repo.GetOpenIssuesCount()),
		})
	}
	m.table.UpdateViewport()
}
