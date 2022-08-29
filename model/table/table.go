package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fzdwx/x/strx"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	KeyMap       KeyMap
	EnableFilter bool
	FilterFunc   Filter
	FilterInput  textinput.Model
	Viewport     viewport.Model
	Styles       Styles
	Rows         []Row

	headers     Headers
	currentRows []Row
	totalRatio  int

	cursor int
	focus  bool
}

func NewModel(headers Headers) *Model {
	m := &Model{
		Viewport:     viewport.New(0, 0),
		Styles:       DefaultStyles(),
		KeyMap:       DefaultKeyMap(),
		EnableFilter: true,
		FilterInput:  textinput.New(),
		cursor:       0,
		FilterFunc:   DefaultFilter(),
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
			m.MoveUp(m.Viewport.Height)
		case key.Matches(msg, m.KeyMap.PageDown):
			m.MoveDown(m.Viewport.Height)
		case key.Matches(msg, m.KeyMap.HalfPageUp):
			m.MoveUp(m.Viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.HalfPageDown):
			m.MoveDown(m.Viewport.Height / 2)
		case key.Matches(msg, m.KeyMap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.KeyMap.GotoTop):
			m.GotoTop()
		case key.Matches(msg, m.KeyMap.GotoBottom):
			m.GotoBottom()
		case key.Matches(msg, m.KeyMap.ShowFilter):
			m.EnableFilter = !m.EnableFilter
			m.doFilter()
		}
	}

	if m.EnableFilter {
		input, cmd := m.FilterInput.Update(msg)
		cmds = append(cmds, cmd)
		m.FilterInput = input
		m.doFilter()
	}

	return m, tea.Batch(cmds...)
}

// View renders the component.
func (m *Model) View() string {
	s := strx.NewFluent()
	s.Str(m.headersView()).NewLine()
	s.Str(m.Viewport.View())

	if m.EnableFilter {
		s.NewLine().Str(m.FilterInput.View())
	}

	return s.String()
}

// UpdateViewport updates the list content based on the previously defined
// columns and currentRows.
func (m *Model) UpdateViewport() {
	renderedRows := make([]string, 0, len(m.currentRows))
	for i := range m.currentRows {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.Viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

// SetRows set rows
func (m *Model) SetRows(rows []Row) {
	m.Rows = rows
	m.currentRows = rows
}

// SetHeaders remember to call UpdateViewport
func (m *Model) SetHeaders(headers Headers) {
	m.headers = headers
	m.totalRatio = headers.TotalRatio()
}

// SetWidth remember to call UpdateViewport
func (m *Model) SetWidth(width int) {
	m.Viewport.Width = width
	m.refreshHeaderMaxWidth(width)
}

// SetHeight remember to call UpdateViewport
func (m *Model) SetHeight(height int) {
	m.Viewport.Height = height - 1
	if m.EnableFilter {
		m.Viewport.Height -= 1
	}
}

// Cursor returns the index of the selected row.
func (m *Model) Cursor() int {
	return m.cursor
}

// MoveUp moves the selection up by any number of row.
// It can not go above the first row.
func (m *Model) MoveUp(n int) {
	m.cursor = clamp(m.cursor-n, 0, len(m.currentRows)-1)
	m.UpdateViewport()

	if m.cursor < m.Viewport.YOffset {
		m.Viewport.SetYOffset(m.cursor)
	}
}

// MoveDown moves the selection down by any number of row.
// It can not go below the last row.
func (m *Model) MoveDown(n int) {
	m.cursor = clamp(m.cursor+n, 0, len(m.currentRows)-1)
	m.UpdateViewport()

	if m.cursor > (m.Viewport.YOffset + (m.Viewport.Height - 1)) {
		m.Viewport.SetYOffset(m.cursor - (m.Viewport.Height - 1))
	}
}

// GotoTop moves the selection to the first row.
func (m *Model) GotoTop() {
	m.MoveUp(m.cursor)
}

// GotoBottom moves the selection to the last row.
func (m *Model) GotoBottom() {
	m.MoveDown(len(m.currentRows))
}

// SelectedRow returns the selected row.
// You can cast it to your own implementation.
func (m *Model) SelectedRow() Row {
	return m.currentRows[m.cursor]
}

// Focused returns the focus state of the table.
func (m *Model) Focused() bool {
	return m.focus
}

// Focus focusses the table, allowing the user to move around the currentRows and
// interact.
func (m *Model) Focus() tea.Cmd {
	var cmd tea.Cmd
	m.focus = true
	if m.EnableFilter {
		cmd = m.FilterInput.Focus()
	}
	m.UpdateViewport()

	return cmd
}

// Blur blurs the table, preventing selection or movement.
func (m *Model) Blur() {
	m.focus = false
	if m.EnableFilter {
		m.FilterInput.Blur()
	}
	m.UpdateViewport()
}

func (m *Model) headersView() string {
	var s = make([]string, len(m.headers))
	for _, h := range m.headers {
		style := lipgloss.NewStyle().Width(h.maxWidth).MaxWidth(h.maxWidth).Inline(true)
		renderedCell := style.Render(runewidth.Truncate(h.Text, h.maxWidth, "…"))
		s = append(s, m.Styles.Header.Render(renderedCell))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, s...)
}

func (m *Model) renderRow(rowID int) string {
	var str = make([]string, len(m.headers))

	for i, value := range m.currentRows[rowID] {
		style := lipgloss.NewStyle().Width(m.headers[i].maxWidth).MaxWidth(m.headers[i].maxWidth).Inline(true)
		renderedCell := m.Styles.Cell.Render(style.Render(runewidth.Truncate(value, m.headers[i].maxWidth, "...")))
		str = append(str, renderedCell)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, str...)

	if rowID == m.cursor {
		return m.Styles.Selected.Render(row)
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

func (m *Model) doFilter() {
	value := m.FilterInput.Value()
	if !m.EnableFilter || len(value) == 0 {
		m.currentRows = m.Rows
	} else {
		filterIndexs, err := m.FilterFunc(m.Rows, value)
		if err != nil {
			return
		}

		var current []Row
		for _, index := range filterIndexs {
			current = append(current, m.Rows[index])
		}

		m.currentRows = current
	}

	if m.cursor > len(m.currentRows) {
		m.cursor = 0
	}

	m.UpdateViewport()
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
