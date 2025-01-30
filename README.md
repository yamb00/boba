# boba
TUI components for Bubble Tea

# Date Picker

```go
package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	datepicker "github.com/yamb00/boba/date-picker"
)

type model struct {
	d datepicker.Model
}

func initialModel() model {
	return model{
		d: datepicker.New(time.Now()),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() string {
	return m.d.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	m.d, cmd = m.d.Update(msg)
	return m, cmd
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

```
