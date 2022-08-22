package repolist

import (
	"context"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fzdwx/code-github-workspace/gh"
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
	user   string

	spinner spinner.Model

	status status

	items  *Items
	cancel context.CancelFunc
}

func New(user string, ops *github.RepositoryListOptions) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &Model{
		Keymap:  DefaultKeyMap(),
		spinner: s,
		ops:     ops,
		user:    user,
	}
}

func (m *Model) Init() tea.Cmd {
	m.status = loading
	go func() {
		ctx, cancelFunc := context.WithCancel(context.Background())
		m.cancel = cancelFunc
		repos, resp, err := gh.GetAuth(ctx).Repositories.List(ctx, m.user, m.ops)

		cobra.CheckErr(err)

		m.items = NewItems(repos, resp)
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
	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
	}

	return m, cmd
}

func (m *Model) View() string {
	if m.status == loading {
		return m.spinnerView()
	}

	//return lipgloss.JoinVertical(lipgloss.Right, blockB) + lipgloss.JoinHorizontal(lipgloss.Top, blockA)
	return m.items.view(m.Width)
}

func (m *Model) spinnerView() string {
	return lipgloss.NewStyle().PaddingTop(m.Height / 3).Width(m.Width).Align(lipgloss.Center).Render(m.spinner.View() + " 正在加载 repository ...")
}
