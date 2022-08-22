package list

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/fzdwx/x/str"
	"github.com/google/go-github/v46/github"
)

type (
	Item struct {
		repo               *github.Repository
		focusedBorderStyle lipgloss.Style
		blurredBorderStyle lipgloss.Style
		viewPort           viewport.Model
	}

	Items struct {
		items []*Item
		resp  *github.Response
	}
)

func NewItem(repo *github.Repository, focusedBorderStyle lipgloss.Style, blurredBorderStyle lipgloss.Style) *Item {
	vp := viewport.New(50, 3)
	return &Item{repo: repo, focusedBorderStyle: focusedBorderStyle, blurredBorderStyle: blurredBorderStyle, viewPort: vp}
}

func (i Item) view(width int, height int) string {
	i.viewPort.Height = height
	i.viewPort.Width = width
	firstLine := str.NewFluent().Str(mapstrp(i.repo.FullName)).Space(2).Str(mapstrp(i.repo.Visibility)).String()
	descLine := str.NewFluent().Str(mapstrp(i.repo.Description) + "asd").String()
	thirdLine := str.NewFluent().
		Str(fmt.Sprintf("%s %d", "*", mapintp(i.repo.StargazersCount))).Space(2).
		Str(fmt.Sprintf("%s %d", "*", mapintp(i.repo.ForksCount))).Space(2).
		Str(fmt.Sprintf("%s %d", "*", mapintp(i.repo.OpenIssuesCount))).Space(2).
		String()

	lipgloss.NewStyle().MaxWidth(width)

	item := str.NewFluent().
		Str(firstLine).
		NewLine().
		Str(descLine).
		NewLine().
		Str(thirdLine).
		String()

	i.viewPort.SetContent(item)
	return i.focusedBorderStyle.Render(i.viewPort.View())
}

func (i Items) view(width int, height int) string {
	fluent := str.NewFluent()
	for _, item := range i.items {
		fluent.Str(item.view(width, 3)).NewLine()
	}
	return lipgloss.JoinVertical(lipgloss.Top, fluent.String())
}

func NewItems(repos []*github.Repository, resp *github.Response) *Items {
	var items []*Item
	focusedBorderStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("238")).MaxWidth(100)
	blurredBorderStyle := lipgloss.NewStyle().Border(lipgloss.HiddenBorder())

	for _, repo := range repos {
		items = append(items, NewItem(repo, focusedBorderStyle, blurredBorderStyle))
	}

	return &Items{items: items, resp: resp}
}

func mapintp(count *int) int {
	if count == nil {
		return 0
	}

	return *count
}

func mapstrp(s *string) string {
	if s == nil {
		return str.Empty
	}
	return *s
}
