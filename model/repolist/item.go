package repolist

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/fzdwx/x/str"
	"github.com/google/go-github/v46/github"
	"github.com/mattn/go-runewidth"
	"unicode/utf8"
)

type (
	Item struct {
		repo               *github.Repository
		focusedBorderStyle lipgloss.Style
		blurredBorderStyle lipgloss.Style
	}

	Items struct {
		items []*Item
		resp  *github.Response
	}
)

func NewItem(repo *github.Repository, focusedBorderStyle lipgloss.Style, blurredBorderStyle lipgloss.Style) *Item {
	return &Item{repo: repo, focusedBorderStyle: focusedBorderStyle, blurredBorderStyle: blurredBorderStyle}
}

func (i Item) view(width int) string {
	cell := width / 20
	return fmt.Sprintf("%s%s%s%s",
		padding(cell*5, mapstrp(i.repo.FullName)),
		padding(cell*7, mapstrp(i.repo.Description)),
		padding(cell, str.Empty),
		padding(cell*2, fmt.Sprintf("🌟 %d", mapintp(i.repo.StargazersCount))),
	)
}

func (i Items) view(width int) string {
	fluent := str.NewFluent()
	for _, item := range i.items {
		fluent.Str(item.view(width)).NewLine()
	}
	return lipgloss.JoinVertical(lipgloss.Top, fluent.String())
}

func NewItems(repos []*github.Repository, resp *github.Response) *Items {
	var items []*Item
	focusedBorderStyle := lipgloss.NewStyle().Background(lipgloss.Color("12"))
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

	return f.Str("...")
}
