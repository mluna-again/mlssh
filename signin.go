package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type signinField int

const (
	signinName signinField = iota
	signinPet
	signinVariant
)

type signinScreen struct {
	nameInput         textinput.Model
	nameInputS        lipgloss.Style
	nameInputFocusedS lipgloss.Style
	renderer          *lipgloss.Renderer
	focusedInput      signinField
	width             int
	heigth            int
}

func newSigninScreen(r *lipgloss.Renderer) signinScreen {
	nameInputS := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(gray)

	nameInputFocusedS := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(fg)

	nameInput := textinput.New()
	nameInput.Placeholder = "Hatchling"
	nameInput.Prompt = ""
	nameInput.Focus()
	nameInput.CharLimit = 30
	nameInput.Width = 28

	return signinScreen{
		nameInput:         nameInput,
		nameInputS:        nameInputS,
		nameInputFocusedS: nameInputFocusedS,
		renderer:          r,
	}
}

func (s signinScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s signinScreen) Update(msg tea.Msg) (signinScreen, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			switch s.focusedInput {
			case signinName:
				s.focusedInput = signinPet
			case signinPet:
				s.focusedInput = signinVariant
			case signinVariant:
				s.focusedInput = signinName
			}
			return s, nil
		}
	}

	if s.focusedInput == signinName {
		s.nameInput, cmd = s.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

func (s signinScreen) View() string {
	nis := s.nameInputFocusedS
	if s.focusedInput != signinName {
		nis = s.nameInputS
	}
	ni := nis.Render(s.nameInput.View())
	ni = lipgloss.JoinVertical(lipgloss.Top, "What name do you want to give it?", ni)

	return lipgloss.Place(s.width, s.heigth, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, ni))
}

func (s *signinScreen) SetHeight(h int) {
	s.heigth = h
}

func (s *signinScreen) SetWidth(w int) {
	s.width = w
}
