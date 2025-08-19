package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errorViewTick struct{}

func nextErrViewTick() tea.Msg {
	time.Sleep(time.Millisecond * time.Duration(800))

	return errorViewTick{}
}

type errorViewModel struct {
	frames      []string
	activeFrame string
	index       int
	width       int
	height      int
}

func newErrorViewModel() errorViewModel {
	return errorViewModel{
		frames:      []string{mimir0, mimir1},
		index:       0,
		activeFrame: mimir0,
	}
}

func (m errorViewModel) Init() tea.Cmd {
	return nil
}

func (m errorViewModel) Update(msg tea.Msg) (errorViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case errorViewTick:
		nextIndex := m.index + 1
		if nextIndex > len(m.frames) {
			nextIndex = 0
		}

		m.index = nextIndex
		m.activeFrame = m.frames[nextIndex]
		return m, nil
	}

	return m, nil
}

func (m errorViewModel) View() string {
	content := m.activeFrame + "sorry, something went wrong! see you!"

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
