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

	renderer          *lipgloss.Renderer
	focusedInput      signinField
	width             int
	heigth            int
	availablePets     []string
	availableVariants []string
	petIndex          int
	variantIndex      int
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
		availablePets:     []string{"Cat", "Bunny", "Turtle"},
		availableVariants: []string{"Ragdoll", "Black"},
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
				if s.petIndex == len(s.availablePets)-1 {
					s.petIndex = 0
				} else {
					s.petIndex++
				}
			}

			if s.focusedInput == signinVariant {
				if s.variantIndex == len(s.availableVariants)-1 {
					s.variantIndex = 0
				} else {
					s.variantIndex++
				}
			}

		case "left":
			if s.focusedInput == signinPet {
				if s.petIndex == 0 {
					s.petIndex = len(s.availablePets) - 1
				} else {
					s.petIndex--
				}
			}

			if s.focusedInput == signinVariant {
				if s.variantIndex == 0 {
					s.variantIndex = len(s.availableVariants) - 1
				} else {
					s.variantIndex--
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
	// NAME INPUT
	nis := s.nameInputFocusedS
	if s.focusedInput != signinName {
		nis = s.nameInputS
	}
	niinput := s.nameInput.View()
	niinputW := lipgloss.Width(niinput)
	ni := nis.Render(niinput)
	ni = lipgloss.JoinVertical(lipgloss.Center, "WHO IS THE USER?", ni)

	// PETS CHOICES
	pbS := s.boxFocusedS
	if s.focusedInput != signinPet {
		pbS = s.boxS
	}
	p := strings.Builder{}
	for i, pet := range s.availablePets {
		w := niinputW / len(s.availablePets)

		if s.petIndex == i {
			p.WriteString(s.blockS.Render(lipgloss.PlaceHorizontal(w, lipgloss.Center, pet)))
		} else {
			p.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, pet))
		}
	}
	pets := pbS.Render(p.String())

	// PET VARIANT CHOICES
	vS := s.boxS
	if s.focusedInput == signinVariant {
		vS = s.boxFocusedS
	}
	variants := vS.Render(lipgloss.PlaceHorizontal(niinputW, lipgloss.Center, "No variants for this species"))
	if s.availablePets[s.petIndex] == "Cat" {
		p := strings.Builder{}
		for i, variant := range s.availableVariants {
			w := niinputW / len(s.availableVariants)

			if s.variantIndex == i {
				p.WriteString(s.blockS.Render(lipgloss.PlaceHorizontal(w, lipgloss.Center, variant)))
			} else {
				p.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, variant))
			}
		}

		variants = vS.Render(p.String())
	}

	// SAVE BUTTON

	return lipgloss.Place(s.width, s.heigth, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Top, ni, pets, variants))
}

func (s *signinScreen) SetHeight(h int) {
	s.heigth = h
}

func (s *signinScreen) SetWidth(w int) {
	s.width = w
}
