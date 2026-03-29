package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerModel is a simple spinner that displays a message while processing.
type SpinnerModel struct {
	spinner  spinner.Model
	message  string
	done     bool
	duration time.Duration
	result   string
}

type spinnerDoneMsg struct {
	result string
}

// NewSpinner creates a spinner with the given message.
func NewSpinner(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	return SpinnerModel{
		spinner: s,
		message: message,
	}
}

func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case spinnerDoneMsg:
		m.done = true
		m.result = msg.result
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.done {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// RunWithSpinner runs a function while displaying a spinner.
// Falls back to simple text output when no TTY is available.
func RunWithSpinner(message string, fn func() (string, error)) (string, error) {
	// Try TUI spinner first
	var result string
	var fnErr error

	m := NewSpinner(message)
	p := tea.NewProgram(m)

	go func() {
		result, fnErr = fn()
		p.Send(spinnerDoneMsg{result: result})
	}()

	if _, err := p.Run(); err != nil {
		// TTY not available — fall back to simple output
		fmt.Printf("  %s\n", message)
		result, fnErr = fn()
		return result, fnErr
	}

	return result, fnErr
}
