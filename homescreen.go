package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type homescreen struct {
	height   int
	width    int
	renderer *lipgloss.Renderer

	bgStyle lipgloss.Style

	user user
}

func newHomescreen(u user, renderer *lipgloss.Renderer) homescreen {
	bg := renderer.
		NewStyle().
		Background(none).
		Foreground(fg)

	return homescreen{
		height:   0,
		width:    0,
		renderer: renderer,
		bgStyle:  bg,
		user:     u,
	}
}

func (h homescreen) Init() tea.Cmd {
	return nil
}

func (h homescreen) Update(msg tea.Msg) (homescreen, tea.Cmd) {
	return h, nil
}

func (h homescreen) View() string {
	msg := "welcome, traveler, come sit with me"
	if !h.user.isNew {
		msg = "nice to see you back"
	}
	content := h.renderer.Place(h.width, h.height, lipgloss.Center, lipgloss.Center, msg)

	return h.bgStyle.Render(content)
}

func (h *homescreen) SetWidth(w int) {
	h.width = w
}

func (h *homescreen) SetHeight(height int) {
	h.height = height
}

func (h *homescreen) SetUser(u user) {
	h.user = u
}
