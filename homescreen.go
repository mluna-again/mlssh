package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type homescreen struct {
	height   int
	width    int
	renderer *lipgloss.Renderer

	bgStyle           lipgloss.Style
	tableS            lipgloss.Style
	progressS         lipgloss.Style
	progressEmptyS    lipgloss.Style
	progressRedS      lipgloss.Style
	progressRedEmptyS lipgloss.Style

	user     user
	settings settings
}

func newHomescreen(renderer *lipgloss.Renderer) homescreen {
	bg := renderer.
		NewStyle().
		Background(none).
		Foreground(fg)

	tableS := renderer.
		NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(fg)

	pbS := renderer.
		NewStyle().
		Background(fg)

	pbES := renderer.
		NewStyle().
		Background(gray)

	pbRS := renderer.
		NewStyle().
		Background(red)

	pbRES := renderer.
		NewStyle().
		Background(redDark)

	return homescreen{
		height:            0,
		width:             0,
		renderer:          renderer,
		bgStyle:           bg,
		tableS:            tableS,
		progressS:         pbS,
		progressEmptyS:    pbES,
		progressRedS:      pbRS,
		progressRedEmptyS: pbRES,
	}
}

func (h homescreen) Init() tea.Cmd {
	return nil
}

func (h homescreen) Update(msg tea.Msg) (homescreen, tea.Cmd) {
	switch msg := msg.(type) {
	case connectToDBMsg:
		if msg.err != nil {
			break
		}

		h.user = msg.user

	case newSettingsMsg:
		if msg.err != nil || msg.ignore {
			break
		}

		h.settings = settings{
			species:    getLunaPet(msg.pet),
			color:      getLunaVariant(msg.variant),
			name:       msg.name,
			readyToUse: true,
		}
	}

	return h, nil
}

func (h homescreen) View() string {
	stats := h.stats()
	thoughts := h.think()
	content := h.renderer.PlaceHorizontal(h.width-2, lipgloss.Left, lipgloss.JoinVertical(lipgloss.Top, stats, "", thoughts))

	return lipgloss.PlaceVertical(h.height, lipgloss.Top, h.tableS.Render(content))
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

func (h homescreen) stats() string {
	hunger := lipgloss.JoinHorizontal(lipgloss.Left, "Hunger ", h.progressBar(50))
	thirst := lipgloss.JoinHorizontal(lipgloss.Left, "Thirst ", h.progressBar(75))
	left := lipgloss.JoinVertical(lipgloss.Top, hunger, thirst)

	happiness := lipgloss.JoinHorizontal(lipgloss.Left, "Hapiness ", h.progressBar(21))
	boredom := lipgloss.JoinHorizontal(lipgloss.Left, "Boredom  ", h.progressBar(54))
	right := lipgloss.JoinVertical(lipgloss.Top, happiness, boredom)

	content := lipgloss.JoinHorizontal(lipgloss.Left, left, "  ", right)

	return content
}

func (h homescreen) think() string {
	thought := "💭 Thinkings about: fish"
	left := lipgloss.PlaceHorizontal(((h.width - 2) / 2), lipgloss.Left, thought)

	album := "🎶 Listeing to: funny monke gif"
	leftW := lipgloss.Width(left)
	right := lipgloss.PlaceHorizontal(h.width-2-leftW, lipgloss.Left, album)

	return lipgloss.JoinHorizontal(lipgloss.Left, left, right)
}

func (h homescreen) progressBar(progress int) string {
	w := (h.width - 19) / 2 // magic numbers: borders - leght of labels
	progressPerBlock := 100 / w
	blocks := progress / progressPerBlock
	padd := w - blocks

	blocksStr := strings.Repeat(" ", blocks)
	paddStr := strings.Repeat(" ", padd)

	fullStyle := h.progressS
	emptyStyle := h.progressEmptyS

	if blocks <= w/3 {
		fullStyle = h.progressRedS
		emptyStyle = h.progressRedEmptyS
	}

	return fullStyle.Render(blocksStr) + emptyStyle.Render(paddStr)
}
