package main

import (
	"math/rand/v2"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/mluna-again/luna/luna"
	"github.com/mluna-again/mlssh/repo"
)

var activities []luna.LunaAnimation = []luna.LunaAnimation{
	luna.IDLE,
	luna.SLEEPING,
}

func randDateInTheFuture() int64 {
	offset := (rand.IntN(90) + 15) * 60

	if DEBUG {
		return time.Now().Unix() + 10
	}

	return time.Now().Unix() + int64(offset)
}

func (m model) randActivity() luna.LunaAnimation {
	last := m.luna.GetAnimation()
	for range 10 {
		next := activities[rand.IntN(len(activities))]
		if m.luna.GetAnimation() != next {
			return next
		}
	}

	// an attempt was made
	return last
}

type activityTick struct {
	next  luna.LunaAnimation
	ready bool
}

func (m model) scheduleActivityChange(skipMissingKeyCheck bool) tea.Cmd {
	return func() tea.Msg {
		// user is not loaded yet (or doesnt have a pk somehow), i should probably handle this better but whatever
		if m.user.publicKey == "" {
			return activityTick{ready: false}
		}

		t := time.Minute
		if DEBUG {
			t = time.Second * 3
		}
		time.Sleep(t)

		ctx, cancel := aLittleBit()
		defer cancel()

		user, err := m.queries.GetUser(ctx, m.user.publicKey)
		if err != nil {
			log.Error(err)
			return activityTick{ready: false}
		}

		tn := time.Unix(user.NextActivityChangeAt, 0)
		ready := time.Now().After(tn)
		nextDate := randDateInTheFuture()
		nextDateFormatted := time.Unix(nextDate, 0).Format("2006-01-02 3:04PM")
		currentDateFormatted := time.Unix(user.NextActivityChangeAt, 0).Format("2006-01-02 3:04PM")
		nowFormatted := time.Now().Format("2006-01-02 3:04PM")
		log.Infof("%s user's pet is scheduled for a change at: %s (current time: %s)", user.Name, currentDateFormatted, nowFormatted)

		if ready {
			// TODO: implement partial update (i don't need to update the name here, but if i don't it sets it to an empty string)
			// ok, there *has* to be a way to make partial updates with sqlc, but im too lazy
			// to look it up
			_, err := m.queries.UpdateUser(ctx, repo.UpdateUserParams{
				NextActivityChangeAt: nextDate,
				Name:                 user.Name,
				PublicKey:            user.PublicKey,
			})
			if err != nil {
				log.Error(err)
				return activityTick{ready: false}
			}
			log.Infof("%s user's pet changed at: %s (next time: %s)", user.Name, nowFormatted, nextDateFormatted)

			return activityTick{ready: true, next: m.randActivity()}
		}

		return activityTick{ready: false}
	}
}
