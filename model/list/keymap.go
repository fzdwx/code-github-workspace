package list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Keymap struct {
	Quit key.Binding
}

func DefaultKeyMap() *Keymap {
	return &Keymap{
		Quit: key.NewBinding(
			key.WithKeys(tea.KeyCtrlC.String()),
		),
	}
}
