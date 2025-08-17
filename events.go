package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type quitSlowlyMsg struct{}

func quitSlowly() tea.Msg {
	time.Sleep(time.Second * 3)
	return quitSlowlyMsg{}
}
