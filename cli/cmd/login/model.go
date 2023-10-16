package login

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#148BF2"))
	cursorStyle  = focusedStyle.Copy()
)

type model struct {
	heading string
	input   *textinput.Model
}

func getEmailInput() textinput.Model {
	input := getNewInput()
	input.Placeholder = "Email"
	input.CharLimit = 32
	/* 	input.Validate = func(text string) error {
	   		_, err := mail.ParseAddress(text)
	   		return err
	   	}
	*/
	return input
}

func getPasswordInput() textinput.Model {
	input := getNewInput()
	input.Placeholder = "Password"
	input.CharLimit = 64
	input.EchoMode = textinput.EchoPassword
	input.EchoCharacter = '•'
	return input
}

func getOTPInput() textinput.Model {
	input := getNewInput()
	input.Placeholder = "MFA OTP"
	input.CharLimit = 6
	input.Validate = func(text string) error {
		//	TODO: Validate the OTP as a number.
		return nil
	}
	return input
}

func getNewInput() textinput.Model {
	input := textinput.New()
	input.Focus()
	input.PromptStyle = focusedStyle
	input.Cursor.Style = cursorStyle
	return input
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	*m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s",
		m.input.View(),
	) + "\n"
}
