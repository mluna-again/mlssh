package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/mluna-again/luna/luna"
	"github.com/mluna-again/mlssh/repo"

	_ "modernc.org/sqlite"
)

func aLittleBit() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	return ctx, cancel
}

func main() {
	host := os.Getenv("MLSSH_HOST")
	port := os.Getenv("MLSSH_PORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "23234"
	}

	s, err := wish.NewServer(
		wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithBanner("come, sit with me, my fellow traveler. let’s sit together and watch the stars die.\n"),
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	l := luna.NewLuna(luna.NewLunaParams{Animation: "sleeping", Pet: "turtle"})
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
		remoteAddr:        s.RemoteAddr().String(),
		user:              user{publicKey: "", name: s.User()},
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

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

type connectToDBMsg struct {
	db      *sql.DB
	queries *repo.Queries
	err     error
	user    user
}

func (m model) connectToDB() tea.Msg {
	if m.originalPublicKey == "" {
		return connectToDBMsg{err: errors.New("sorry, you need to use public key to access the server")}
	}

	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		log.Error(err)
		return connectToDBMsg{err: err}
	}

	queries := repo.New(db)

	newUser := false
	ctx, cancel := aLittleBit()
	defer cancel()
	u, err := queries.GetUser(ctx, m.originalPublicKey)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error(err)
		return connectToDBMsg{err: err}
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Info(fmt.Sprintf("creating a new user: %s", m.originalUsername))
		u, err = queries.CreateUser(ctx, repo.CreateUserParams{
			Name:      m.originalUsername,
			PublicKey: m.originalPublicKey,
		})
		newUser = true
		if err != nil {
			log.Error(err)
			return connectToDBMsg{err: err}
		}
	}

	log.Info(fmt.Sprintf("user %s logged", u.Name))
	return connectToDBMsg{
		db:      db,
		err:     nil,
		queries: queries,
		user: user{
			// use the provided username, it's whatever
			name:      m.originalUsername,
			publicKey: u.PublicKey,
			isNew:     newUser,
		},
	}
}

type quitSlowlyMsg struct{}

func quitSlowly() tea.Msg {
	time.Sleep(time.Second * 3)
	return quitSlowlyMsg{}
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

	content := lipgloss.JoinVertical(lipgloss.Center, fmt.Sprintf("Hi, %s. nice %s ip", m.user.name, m.remoteAddr), m.luna.View(), backMsg)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
