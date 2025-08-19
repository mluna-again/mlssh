package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrations embed.FS

var DEBUG bool = false

var timeRangeMin int = 15
var timeRangeMax int = 90

func main() {
	flag.BoolVar(&DEBUG, "debug", false, "turn debug mode on (time passes fater, more logs, etc)")
	flag.IntVar(&timeRangeMin, "wait-min", 15, "minimum time to wait before activity change (minutes)")
	flag.IntVar(&timeRangeMax, "wait-max", 90, "minimum time to wait before activity change (minutes)")
	flag.Parse()
	log.Info(timeRangeMax)
	log.Info(timeRangeMin)

	if DEBUG {
		log.Info("DEBUG mode on")
	}

	log.Info("Running migrations... ")
	db, err := migrateDatabase()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Migrations ran.")

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
		wish.WithBanner(bannerTXT),
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				return teaHandler(sess, db)
			}),
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
	log.Info("Closing DB connection")
	err = db.Close()
	if err != nil {
		log.Error(err)
	} else {
		log.Info("DB closed")
	}

	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session, db *sql.DB) (tea.Model, []tea.ProgramOption) {
	m, err := newModel(s, db)
	if err != nil {
		log.Fatal(err)
	}
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func migrateDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		return nil, err
	}
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrations)
	err = goose.SetDialect("sqlite")
	if err != nil {
		return nil, err
	}
	err = goose.Up(db, "migrations")
	if err != nil {
		return nil, err
	}

	return db, nil
}
