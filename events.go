package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type quitSlowlyMsg struct{}

func (m model) quitSlowly() tea.Msg {
	time.Sleep(time.Second * 3)
	return quitSlowlyMsg{}
}

func (m model) quitWithError(err error) (tea.Model, tea.Cmd) {
	m.quitting = true
	m.err = err
	log.Error(err)
	return m, tea.Batch(nextErrViewTick, m.quitSlowly)
}
