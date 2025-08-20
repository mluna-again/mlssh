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
	settings          settings

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
	errView  errorViewModel

	// screen managment
	currentScreen availableScreen

	// components
	luna         luna.LunaModel
	homescreen   homescreen
	signinscreen signinScreen
}

func newModel(s ssh.Session, db *sql.DB) (model, []error) {
	l, err := luna.NewLuna(luna.NewLunaParams{
		Animation: luna.SLEEPING,
		Pet:       luna.CAT,
		Size:      luna.SMALL,
		ResizePoints: luna.LunaResizePoints{
			HeightLarge:  40,
			WidthLarge:   40,
			HeightMedium: 30,
			WidthMedium:  30,
		},
	})
	if len(err) > 0 {
		return model{}, err
	}
	l.DisableKeys()
	l.SetAutoresize(true)
	l.ShowName()

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
		// TODO: remove user, i should fetch it from the connectToDBMsg event to avoid having stale data
		//       check signinscreen to see how to do it
		homescreen:   newHomescreen(renderer),
		signinscreen: newSigninScreen(renderer),
		renderer:     renderer,
		db:           db,
		errView:      newErrorViewModel(),
	}

	return m, nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.luna.Init(), m.connectToDB, m.scheduleActivityChange)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}

	switch msg := msg.(type) {
	case activityTick:
		if msg.err != nil {
			return m.quitWithError(msg.err)
		}

		if msg.ready {
			m.luna.SetAnimation(msg.next)
		}
		m.ready = !msg.waiting
		return m, m.scheduleActivityChange

	case quitSlowlyMsg:
		return m, tea.Quit

	case connectToDBMsg:
		if msg.err != nil {
			return m.quitWithError(msg.err)
		}
		m.db = msg.db
		m.queries = msg.queries
		m.user = msg.user

		m.homescreen.SetUser(msg.user)

		if m.user.isNew {
			m.currentScreen = signin
		}

		m.settings = msg.settings
		m.updateLuna()

	case newSettingsMsg:
		if msg.ignore {
			break
		}

		if msg.err != nil {
			return m.quitWithError(msg.err)
		}

		m.settings = settings{
			species:    getLunaPet(msg.pet),
			color:      getLunaVariant(msg.variant),
			name:       msg.name,
			readyToUse: true,
		}
		m.updateLuna()
		m.currentScreen = home

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		lunaH := lipgloss.Height(m.luna.View())

		m.homescreen.SetHeight(msg.Height - lunaH)
		m.homescreen.SetWidth(msg.Width)

		m.signinscreen.SetHeight(msg.Height)
		m.signinscreen.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, m.quitSlowly(800)
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

	m.errView, cmd = m.errView.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting && m.err != nil {
		return m.errView.View()
	}

	if m.quitting {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "byebye!")
	}

	if !m.ready {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, "loading...")
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

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, dejaVu, "you took a wrong turn"))
}

func (m *model) updateLuna() {
	if !m.settings.readyToUse {
		return
	}

	log.Infof("%s new settings loaded", m.user.name)
	m.luna.SetPet(m.settings.species)
	m.luna.SetName(m.settings.name)
	m.luna.SetVariant(m.settings.color)
}
