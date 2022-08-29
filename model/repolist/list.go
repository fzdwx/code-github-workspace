package repolist

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fzdwx/gh-sp/api"
	"github.com/fzdwx/gh-sp/model/divid"
	"github.com/fzdwx/gh-sp/model/table"
	"github.com/fzdwx/x/strx"
	"github.com/google/go-github/v46/github"
	"github.com/mattn/go-runewidth"
	"time"
)

type status int

const (
	loading status = iota
	loaded
	heightNotEnough
)

var (
	noticeHeight = 7
)

type Model struct {
	Keymap *Keymap
	Width  int
	Height int
	Ops    *github.RepositoryListOptions

	spinner spinner.Model

	status status

	table  *table.Model
	cancel context.CancelFunc
	repos  []*github.Repository
	err    error
}

func New(ops *github.RepositoryListOptions) *Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &Model{
		Keymap:  DefaultKeyMap(),
		spinner: s,
		Ops:     ops,
		table:   table.NewModel(headers()),
	}
}

func (m *Model) Init() tea.Cmd {
	m.status = loading
	go m.fetchRepoList(m.Ops)

	return tea.Batch(m.table.Focus(), spinner.Tick)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		return m, tea.Quit
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.Quit):
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys(tea.KeyCtrlA.String()))):
			err := api.Browse(m.repos[m.table.Cursor()].GetHTMLURL())
			m.err = err
			return m, nil
		}
	case tea.WindowSizeMsg:
		if msg.Height < noticeHeight+2 {
			m.status = heightNotEnough
			return m, tea.Quit
		}

		m.Width = msg.Width
		m.Height = msg.Height
	case spinner.TickMsg:
		model, cmd := m.spinner.Update(msg)
		m.spinner = model
		cmds = append(cmds, cmd)
	}

	model, tCmd := m.table.Update(msg)
	m.table = model
	cmds = append(cmds, tCmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.err != nil {
		return "Error:" + m.err.Error() + strx.NewLine
	}

	if m.status == loading {
		return m.spinnerView()
	}

	if m.status == heightNotEnough {
		api.Error().Msgf("Repo list最小需要%d格的高度", noticeHeight+2)
		return ""
	}

	s := strx.NewFluent().
		WriteFunc(m.tableView).
		NewLine().Str(divid.Line(m.Width)).
		NewLine().Str(m.currentRepoView())

	return s.String()
}

func (m *Model) spinnerView() string {
	return m.centerStyle().Render(m.spinner.View() + " 正在加载 repository ...")
}

func (m *Model) centerStyle() lipgloss.Style {
	return lipgloss.NewStyle().PaddingTop(m.Height / 3).Width(m.Width).Align(lipgloss.Center)
}

func (m *Model) addRows(repos []*github.Repository) {
	var rows []table.Row
	for _, repo := range repos {
		rows = append(rows, table.Row{
			repo.GetFullName(),
			repo.GetDescription(),
			fmt.Sprintf("🌟 %d", repo.GetStargazersCount()),
			repo.GetVisibility(),
			fmt.Sprintf("🎯 %d", repo.GetOpenIssuesCount()),
		})
	}
	m.table.SetRows(rows)
	m.table.UpdateViewport()
}

func (m *Model) SetRepos(repos []*github.Repository) {
	m.repos = repos
	m.addRows(repos)
}

func (m *Model) getTableHeight() int {
	if m.Height >= m.minHeight() {
		return len(m.repos) + 1
	}

	return m.Height - noticeHeight
}

func (m *Model) minHeight() int {
	return len(m.repos) + 1 + noticeHeight
}

func (m *Model) currentRepoView() string {
	repo := m.repos[m.table.Cursor()]

	repoStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("231"))
	baseStyle := lipgloss.NewStyle().Padding(0, 1)
	descStyle := baseStyle.Copy().MaxHeight(2).Width(m.Width).Italic(true)

	repoText := baseStyle.Render("\uF401") + strx.Space + repoStyle.Render(repo.GetFullName())
	repoTime := fuzzTime(m.Width, runewidth.StringWidth(repoText), repo)

	s := strx.NewFluent().
		Str(repoText + repoTime).
		NewLine().Str(descStyle.Render(repo.GetDescription()))
	return s.String()
}

func (m *Model) fetchRepoList(ops *github.RepositoryListOptions) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	repos, _, err := api.Get(ctx).Repositories.List(ctx, "", ops)

	m.err = err
	m.cancel = cancelFunc
	m.SetRepos(repos)
	m.status = loaded
}

func (m *Model) tableView(s *strx.FluentStringBuilder) {
	m.table.SetWidth(m.Width)
	m.table.SetHeight(m.getTableHeight())
	m.table.UpdateViewport()
	s.Str(m.table.View())
}

func headers() table.Headers {
	return []*table.Header{
		{Text: "Repo Name", Ratio: 6},
		{Text: "Description", Ratio: 11, MinWidth: 20},
		{Text: "Stars", Ratio: 2, MinWidth: 5},
		{Text: "Visible", Ratio: 2, MinWidth: 7},
		{Text: "Issues", Ratio: 2, MinWidth: 10},
	}
}

func fuzzTime(width int, stringWidth int, repo *github.Repository) string {
	t := repo.GetPushedAt()
	if repo.PushedAt == nil {
		t = repo.GetCreatedAt()
	}
	return strx.RepeatSpace(width-stringWidth) + strx.FuzzyAgoAbbr(time.Now(), t.Time)
}
