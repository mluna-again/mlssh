package main

import (
	"database/sql"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mluna-again/luna/luna"
	"github.com/mluna-again/mlssh/repo"
)

type user struct {
	publicKey string
	name      string
	isNew     bool
}

type model struct {
	term              string
	profile           string
	width             int
	height            int
	bg                string
	txtStyle          lipgloss.Style
	quitStyle         lipgloss.Style
	luna              luna.LunaModel
	originalUsername  string
	originalPublicKey string
	remoteAddr        string
	db                *sql.DB
	queries           *repo.Queries
	quitting          bool
	err               error
	user              user
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.luna.Init(), m.connectToDB)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case quitSlowlyMsg:
		// TODO: how do i gracefully close the db?
		if m.db != nil {
			_ = m.db.Close()
		}
		return m, tea.Quit

	case connectToDBMsg:
		if msg.err != nil {
			m.quitting = true
			m.err = msg.err
			log.Error(m.err)
			return m, quitSlowly
		}
		m.db = msg.db
		m.queries = msg.queries
		m.user = msg.user

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var lunaCmd tea.Cmd
	m.luna, lunaCmd = m.luna.Update(msg)
	cmds = append(cmds, lunaCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting && m.err != nil {
		return "sorry, something went wrong! see you!"
	}

	if m.quitting {
		return "bye!"
	}

	backMsg := ""
	if m.user.isNew {
		backMsg = "welcome!"
	} else {
		backMsg = "nice to see you back :)"
	}

	content := lipgloss.JoinVertical(lipgloss.Center, fmt.Sprintf("Hi, %s. nice ip (%s) ;)", m.user.name, m.remoteAddr), m.luna.View(), backMsg)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
