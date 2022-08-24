package table

import "github.com/charmbracelet/lipgloss"

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