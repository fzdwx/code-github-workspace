package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	KeyMap     KeyMap
	headers    Headers
	rows       []Row
	totalRatio int
	styles     Styles

	viewport viewport.Model
	cursor   int
	focus    bool
}

func NewModel(headers Headers) *Model {
	m := &Model{
		viewport: viewport.New(0, 0),
		styles:   DefaultStyles(),
		KeyMap:   DefaultKeyMap(),
		cursor:   0,
	}

	m.SetHeaders(headers)
	m.SetWidth(100)
	m.SetHeight(10)
	m.UpdateViewport()
	return m
}

// Update is the Bubble Tea update loop.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.LineUp):
			m.MoveUp(1)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.PageUp):
			m.MoveUp(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.viewport.Height)
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		}
	}
	return m, tea.Batch(cmds...)
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

// AppendRow remember to call UpdateViewport
func (m *Model) AppendRow(row Row) {
	m.rows = append(m.rows, row)
}

// SetHeaders remember to call UpdateViewport
func (m *Model) SetHeaders(headers Headers) {
	m.headers = headers
	m.totalRatio = headers.TotalRatio()
}

// SetWidth remember to call UpdateViewport
func (m *Model) SetWidth(width int) {
	m.viewport.Width = width
	m.refreshHeaderMaxWidth(width)
}

// SetHeight remember to call UpdateViewport
func (m *Model) SetHeight(height int) {
	m.viewport.Height = height - 1
}

// Cursor returns the index of the selected row.
func (m *Model) Cursor() int {
	return m.cursor
}

// MoveUp moves the selection up by any number of row.
// It can not go above the first row.
func (m *Model) MoveUp(n int) {
	m.cursor = clamp(m.cursor-n, 0, len(m.rows)-1)
	m.UpdateViewport()

	if m.cursor < m.viewport.YOffset {
		m.viewport.SetYOffset(m.cursor)
	}
}

// MoveDown moves the selection down by any number of row.
// It can not go below the last row.
func (m *Model) MoveDown(n int) {
	m.cursor = clamp(m.cursor+n, 0, len(m.rows)-1)
	m.UpdateViewport()

	if m.cursor > (m.viewport.YOffset + (m.viewport.Height - 1)) {
		m.viewport.SetYOffset(m.cursor - (m.viewport.Height - 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *Model) GotoTop() {
	m.MoveUp(m.cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Model) GotoBottom() {
	m.MoveDown(len(m.rows))
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m *Model) SelectedRow() Row {
	return m.rows[m.cursor]
}

// Focused returns the focus state of the table.
func (m *Model) Focused() bool {
	return m.focus
}

// Focus focusses the table, allowing the user to move around the rows and
// interact.
func (m *Model) Focus() {
	m.focus = true
	m.UpdateViewport()
}

// Blur blurs the table, preventing selection or movement.
func (m *Model) Blur() {
	m.focus = false
	m.UpdateViewport()
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

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
