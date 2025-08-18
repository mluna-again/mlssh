package main

import (
	"database/sql"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/mluna-again/luna/luna"
	"github.com/mluna-again/mlssh/repo"
)

type settings struct {
	species    luna.LunaPet
	name       string
	color      luna.LunaVariant
	readyToUse bool
}

type connectToDBMsg struct {
	db       *sql.DB
	queries  *repo.Queries
	err      error
	user     user
	settings settings
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

	ctx, cancel := aLittleBit()
	defer cancel()
	userWithSettings, err := queries.GetUser(ctx, m.originalPublicKey)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Error(err)
		return connectToDBMsg{err: err}
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		log.Info(fmt.Sprintf("creating a new user: %s", m.originalUsername))
		u, err := queries.CreateUser(ctx, repo.CreateUserParams{
			Name:      m.originalUsername,
			PublicKey: m.originalPublicKey,
		})
		if err != nil {
			log.Error(err)
			return connectToDBMsg{err: err}
		}

		log.Info(fmt.Sprintf("new user %s logged", u.Name))
		return connectToDBMsg{
			db:      db,
			err:     nil,
			queries: queries,
			user: user{
				// use the provided username, it's whatever
				name:      m.originalUsername,
				publicKey: u.PublicKey,
				isNew:     true,
			},
		}
	}

	log.Info(fmt.Sprintf("user %s logged", userWithSettings.Name))

	settings := settings{readyToUse: false}
	if userWithSettings.InsertedAt.Valid {
		settings.color = getLunaVariant(userWithSettings.PetColor.String)
		settings.species = getLunaPet(userWithSettings.PetSpecies.String)
		settings.name = userWithSettings.PetName.String
		settings.readyToUse = true
	}

	return connectToDBMsg{
		db:      db,
		err:     nil,
		queries: queries,
		user: user{
			// use the provided username, it's whatever
			name:      m.originalUsername,
			publicKey: userWithSettings.PublicKey,
			isNew:     !userWithSettings.InsertedAt.Valid,
		},
		settings: settings,
	}
}
