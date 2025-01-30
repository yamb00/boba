package datepicker

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	CursorUp    key.Binding
	CursorDown  key.Binding
	CursorLeft  key.Binding
	CursorRight key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down"),
		),
		CursorLeft: key.NewBinding(
			key.WithKeys("left"),
		),
		CursorRight: key.NewBinding(
			key.WithKeys("right"),
		),
	}
}
