package main

import (
	"database/sql"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/mluna-again/luna/luna"
	"github.com/mluna-again/mlssh/repo"
)

type availableScreen int

const (
	home availableScreen = iota
	signin
)

type user struct {
	publicKey string
	name      string
	isNew     bool
}

type model struct {
	// user stuff
	originalUsername  string
	originalPublicKey string
	remoteAddr        string
	user              user

	// styles
	txtStyle  lipgloss.Style
	quitStyle lipgloss.Style
	renderer  *lipgloss.Renderer

	// terminal stuff
	profile string
	width   int
	height  int
	term    string
	bg      string

	// app stuff
	db       *sql.DB
	queries  *repo.Queries
	quitting bool
	err      error
	ready    bool

	// screen managment
	currentScreen availableScreen

	// components
	luna         luna.LunaModel
	homescreen   homescreen
	signinscreen signinScreen
}

func newModel(s ssh.Session) (model, []error) {
	l, err := luna.NewLuna(luna.NewLunaParams{Animation: luna.SLEEPING, Pet: luna.CAT, Size: luna.SMALL})
	if len(err) > 0 {
		return model{}, err
	}
	l.DisableKeys()

	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)
	txtStyle := renderer.NewStyle().Foreground(lipgloss.Color("10"))
	quitStyle := renderer.NewStyle().Foreground(lipgloss.Color("8"))

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	pk := s.PublicKey()
	pkStr := ""
	if pk != nil {
		pkStr = string(pk.Marshal())
	}

	u := user{publicKey: "", name: s.User()}
	m := model{
		term:              pty.Term,
		profile:           renderer.ColorProfile().Name(),
		width:             pty.Window.Width,
		height:            pty.Window.Height,
		bg:                bg,
		txtStyle:          txtStyle,
		quitStyle:         quitStyle,
		luna:              l,
		originalUsername:  s.User(),
		originalPublicKey: pkStr,
		remoteAddr:        removePort(s.RemoteAddr().String()),
		user:              u,
		homescreen:        newHomescreen(u, renderer),
		signinscreen:      newSigninScreen(),
		renderer:          renderer,
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.luna.Init(), m.connectToDB, m.scheduleActivityChange(true))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case activityTick:
		if msg.ready {
			m.luna.SetAnimation(msg.next)
		}
		return m, m.scheduleActivityChange(false)

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

		m.ready = true
		m.homescreen.SetUser(msg.user)

		if m.user.isNew {
			m.currentScreen = signin
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		lunaH := lipgloss.Height(m.luna.View())

		m.homescreen.SetHeight(msg.Height - lunaH)
		m.homescreen.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentScreen {
	case home:
		m.homescreen, cmd = m.homescreen.Update(msg)
		cmds = append(cmds, cmd)

	case signin:
		m.signinscreen, cmd = m.signinscreen.Update(msg)
		cmds = append(cmds, cmd)
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

	switch m.currentScreen {
	case home:
		screen := m.homescreen.View()
		petCentered := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, m.luna.View())
		return lipgloss.JoinVertical(lipgloss.Top, screen, petCentered)

	case signin:
		screen := m.signinscreen.View()
		return screen
	}

	return "you should not be able to see this"
}
