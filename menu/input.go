package menu

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) HandleInput(msg tea.KeyMsg) (tea.Cmd, error) {
	var cmd tea.Cmd
	switch msg.String() {
	case "i", "I":
		if m.Alt {
			m.Alt = false
			m.AutoAlt = false
			return tea.ExitAltScreen, nil
		} else {
			m.Alt = true
			m.AutoAlt = false
			return tea.EnterAltScreen, nil
		}
	case "h", "left":
		m.Leap(-1)
	case "l", "right":
		m.Leap(1)
	case "up", "k":
		m.Index--
	case "down", "j":
		m.Index++
	case "q", "esc", "ctrl+c":
		m.Cancel = true
		return tea.Quit, nil
	case "space", " ", "enter", "ctrl+d":
		m.Done = true
		return tea.Quit, nil
	}

	m.Validate()

	return cmd, nil
}

func (m *Model) Leap(direction int) {
	for {
		source, target := m.Result()
		m.Index += direction
		m.Validate()
		newSource, newTarget := m.Result()
		if target == newTarget {
			return
		}

		if source != newSource {
			return
		}
	}
}
