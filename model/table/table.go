package table

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	headers    Headers
	rows       []Row
	totalRatio int
	styles     Styles

	viewport viewport.Model
	cursor   int
}

func NewModel(headers Headers) *Model {
	m := &Model{
		viewport: viewport.New(0, 0),
		styles:   DefaultStyles(),
		cursor:   0,
	}

	m.SetHeaders(headers)
	m.SetWidth(100)
	m.SetHeight(10)

	return m
}

// AppendRow remember to call UpdateViewport
func (m *Model) AppendRow(row Row) {
	m.rows = append(m.rows, row)
}

// View renders the component.
func (m *Model) View() string {
	return m.headersView() + "\n" + m.viewport.View()
}

// UpdateViewport updates the list content based on the previously defined
// columns and rows.
func (m *Model) UpdateViewport() {
	renderedRows := make([]string, 0, len(m.rows))
	for i := range m.rows {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m *Model) headersView() string {
	var s = make([]string, len(m.headers))
	for _, h := range m.headers {
		style := lipgloss.NewStyle().Width(h.maxWidth).MaxWidth(h.maxWidth).Inline(true)
		renderedCell := style.Render(runewidth.Truncate(h.Text, h.maxWidth, "…"))
		s = append(s, m.styles.Header.Render(renderedCell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *Model) renderRow(rowID int) string {
	var str = make([]string, len(m.headers))

	for i, value := range m.rows[rowID] {
		style := lipgloss.NewStyle().Width(m.headers[i].maxWidth).MaxWidth(m.headers[i].maxWidth).Inline(true)
		renderedCell := m.styles.Cell.Render(style.Render(runewidth.Truncate(value, m.headers[i].maxWidth, "...")))
		str = append(str, renderedCell)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, str...)

	if rowID == m.cursor {
		return m.styles.Selected.Render(row)
	}

	return row
}

func (m *Model) SetHeaders(headers Headers) {
	m.headers = headers
	m.totalRatio = headers.TotalRatio()
	m.UpdateViewport()
}

func (m *Model) SetWidth(width int) {
	m.viewport.Width = width
	m.refreshHeaderMaxWidth(width)
	m.UpdateViewport()
}

func (m *Model) SetHeight(height int) {
	m.viewport.Height = height - 1
	m.UpdateViewport()
}

func (m *Model) refreshHeaderMaxWidth(width int) {
	n := width / m.totalRatio
	for _, header := range m.headers {
		max := n * header.Ratio

		if header.MinWidth != 0 && header.MinWidth > max {
			max = header.MinWidth
		}

		header.maxWidth = max
	}
}
