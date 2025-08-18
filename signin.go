package main

import tea "github.com/charmbracelet/bubbletea"

type signinScreen struct{}

func newSigninScreen() signinScreen {
	return signinScreen{}
}

func (s signinScreen) Init() tea.Cmd {
	return nil
}

func (s signinScreen) Update(msg tea.Msg) (signinScreen, tea.Cmd) {
	return s, nil
}

func (s signinScreen) View() string {
	return "hello friend, hello, friend."
}
