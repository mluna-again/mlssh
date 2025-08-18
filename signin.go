package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type signinScreen struct {
	nameInput  textinput.Model
	nameInputS lipgloss.Style
	renderer   *lipgloss.Renderer
}

func newSigninScreen(r *lipgloss.Renderer) signinScreen {
	nameInputS := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(fg)

	nameInput := textinput.New()
	nameInput.Placeholder = "Hatchling"
	nameInput.Prompt = ""
	nameInput.Focus()
	nameInput.CharLimit = 30
	nameInput.Width = 28

	return signinScreen{
		nameInput:  nameInput,
		nameInputS: nameInputS,
		renderer:   r,
	}
}

func (s signinScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s signinScreen) Update(msg tea.Msg) (signinScreen, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	s.nameInput, cmd = s.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s signinScreen) View() string {
	greetings := "hello friend, hello, friend.\n"
	ni := s.nameInputS.Render(s.nameInput.View())
	ni = lipgloss.JoinVertical(lipgloss.Top, "What name do you want to give it?", ni)

	return lipgloss.JoinVertical(lipgloss.Top, greetings, ni)
}
