package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type signinField int

const (
	signinName signinField = iota
	signinPet
	signinVariant
	signinGo
)

type signinScreen struct {
	nameInput textinput.Model

	nameInputS        lipgloss.Style
	nameInputFocusedS lipgloss.Style
	boxS              lipgloss.Style
	boxFocusedS       lipgloss.Style
	blockS            lipgloss.Style

	renderer     *lipgloss.Renderer
	focusedInput signinField
	width        int
	heigth       int
	availabePets []string
	petIndex     int
}

func newSigninScreen(r *lipgloss.Renderer) signinScreen {
	nameInputS := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(gray)

	nameInputFocusedS := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(fg)

	boxFocusedS := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(fg)

	boxS := r.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(gray)

	blockS := r.NewStyle().
		Background(fg).
		Foreground(bg)

	nameInput := textinput.New()
	nameInput.Placeholder = "Hatchling"
	nameInput.Prompt = ""
	nameInput.Focus()
	nameInput.CharLimit = 30
	nameInput.Width = 29 // this number + 1 (idk why lipgloss.Width returns it +1) has to be divisible by len(availabePets)

	return signinScreen{
		nameInput:         nameInput,
		nameInputS:        nameInputS,
		boxS:              boxS,
		boxFocusedS:       boxFocusedS,
		nameInputFocusedS: nameInputFocusedS,
		blockS:            blockS,
		renderer:          r,
		availabePets:      []string{"Cat", "Bunny", "Turtle"},
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
		case "right":
			if s.focusedInput == signinPet {
				if s.petIndex == len(s.availabePets)-1 {
					s.petIndex = 0
				} else {
					s.petIndex++
				}
			}

		case "left":
			if s.focusedInput == signinPet {
				if s.petIndex == 0 {
					s.petIndex = len(s.availabePets) - 1
				} else {
					s.petIndex--
				}
			}

		case "tab":
			switch s.focusedInput {
			case signinName:
				s.focusedInput = signinPet
			case signinPet:
				s.focusedInput = signinVariant
			case signinVariant:
				s.focusedInput = signinGo
			case signinGo:
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
	niinput := s.nameInput.View()
	ni := nis.Render(niinput)
	ni = lipgloss.JoinVertical(lipgloss.Center, "WHO IS THE USER?", ni)

	pbS := s.boxFocusedS
	if s.focusedInput != signinPet {
		pbS = s.boxS
	}
	p := strings.Builder{}
	for i, pet := range s.availabePets {
		w := lipgloss.Width(niinput) / len(s.availabePets)

		if s.petIndex == i {
			p.WriteString(s.blockS.Render(lipgloss.PlaceHorizontal(w, lipgloss.Center, pet)))
		} else {
			p.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, pet))
		}
	}
	pets := pbS.Render(p.String())

	return lipgloss.Place(s.width, s.heigth, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, ni, pets))
}

func (s *signinScreen) SetHeight(h int) {
	s.heigth = h
}

func (s *signinScreen) SetWidth(w int) {
	s.width = w
}
