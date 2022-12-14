package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
	"sort"
)

type Header struct {
	Text     string
	Ratio    int
	MinWidth int

	maxWidth int
}

func NewHeader(text string, ratio int, minWidth int) *Header {
	return &Header{text, ratio, minWidth, 0}
}

type Headers []*Header

func (h Headers) TotalRatio() int {
	total := 0

	for _, header := range h {
		total += header.Ratio
	}

	return total
}

type Row []string

// Styles contains style definitions for this list component. By default, these
// values are generated by DefaultStyles.
type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
}

// DefaultStyles returns a set of default style definitions for this table.
func DefaultStyles() Styles {
	return Styles{
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
		Cell:     lipgloss.NewStyle().Padding(0, 1),
	}
}

// KeyMap defines keybindings. It satisfies to the help.KeyMap interface, which
// is used to render the menu menu.
type KeyMap struct {
	LineUp       key.Binding
	LineDown     key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	GotoTop      key.Binding
	GotoBottom   key.Binding
	ShowFilter   key.Binding
}

// DefaultKeyMap returns a default set of keybindings.
func DefaultKeyMap() KeyMap {
	const spacebar = " "
	return KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys(tea.KeyDown.String()),
			key.WithHelp("↓", "down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", spacebar),
			key.WithHelp("pgdn", "page down"),
		),
		HalfPageUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
			key.WithHelp("^U", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("^D", "½ page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "go to end"),
		),
		ShowFilter: key.NewBinding(
			key.WithKeys("ctrl+f"),
			key.WithHelp("^F", "show filter input"),
		),
	}
}

type Filter func(rows []Row, val string) ([]int, error)

func DefaultFilter() Filter {
	return func(rows []Row, val string) ([]int, error) {
		var str []string
		for _, row := range rows {
			str = append(str, row[0])
		}
		matches := fuzzy.Find(val, str)
		sort.Stable(matches)

		var res []int
		for _, match := range matches {
			res = append(res, match.Index)
		}

		return res, nil
	}
}
