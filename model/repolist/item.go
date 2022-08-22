package repolist

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/fzdwx/x/str"
	"github.com/google/go-github/v46/github"
	"github.com/mattn/go-runewidth"
	"strings"
	"unicode/utf8"
)

type (
	Item struct {
		repo *github.Repository
	}

	Items struct {
		currentIdx int
		items      []*Item
		resp       *github.Response
	}
)

func NewItem(repo *github.Repository) *Item {
	return &Item{repo: repo}
}

func (i Item) view(width int, repoNameStyle lipgloss.Style, itemStyle lipgloss.Style) string {
	cell := width / 20

	item := fmt.Sprintf("%s%s%s%s%s",
		repoNameStyle.Render(padding(cell*5, i.repo.GetFullName())),
		padding(cell*10, i.repo.GetDescription()),
		padding(cell*3, fmt.Sprintf(" 🌟 %d", i.repo.GetStargazersCount())),
		padding(cell*2, fmt.Sprintf("%s", i.repo.GetVisibility())),
		padding(cell*2, fmt.Sprintf(" 🎯 %d", i.repo.GetOpenIssues())),
	)
	return itemStyle.Render(item)
}

func (i Items) view(width int) string {
	fluent := str.NewFluent()
	var focusedStyle = lipgloss.NewStyle().Background(lipgloss.Color("13"))
	var blurredStyle = lipgloss.NewStyle()
	var repoNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Bold(true)

	for idx, item := range i.items {
		var itemStyle = blurredStyle

		if idx == i.currentIdx {
			itemStyle = focusedStyle
		}

		fluent.Str(item.view(width, repoNameStyle, itemStyle)).NewLine()
	}
	return lipgloss.JoinVertical(lipgloss.Top, fluent.String())
}

func NewItems(repos []*github.Repository, resp *github.Response) *Items {
	var items []*Item

	for _, repo := range repos {
		items = append(items, NewItem(repo))
	}

	return &Items{items: items, resp: resp, currentIdx: 0}
}

func padding(size int, s string) string {
	width := runewidth.StringWidth(s)
	if width > size {
		f := limit(size, s)
		return f.String()
	}

	return s + str.RepeatSpace(size-width)
}

func limit(size int, s string) *str.FluentStringBuilder {
	bytes := []byte(s)
	f := str.NewFluent()
	rwrw := 0
	for i, w := 0, 0; i < len(bytes); i += w {
		if rwrw >= size-3 {
			break
		}

		r, rw := utf8.DecodeRune(bytes[i:])
		if r == utf8.RuneError {
			panic("could not decode rune")
		}

		rwrw += runewidth.RuneWidth(r)
		f.Str(string(r))
		w = rw
	}

	return f.Str(strings.Repeat(".", size-rwrw))
}
