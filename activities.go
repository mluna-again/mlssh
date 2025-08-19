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
	offset := (rand.IntN(timeRangeMax) + timeRangeMin) * 60

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
	next    luna.LunaAnimation
	ready   bool
	waiting bool
	err     error
}

func (m model) scheduleActivityChange() tea.Msg {
	// user is not loaded yet (or doesnt have a pk somehow), i should probably handle this better but whatever
	if m.user.publicKey == "" {
		if DEBUG {
			log.Error("no key")
		}
		return activityTick{ready: false, waiting: true}
	}

	t := time.Millisecond * 500
	time.Sleep(t)

	ctx, cancel := aLittleBit()
	defer cancel()

	user, err := m.queries.GetUser(ctx, m.user.publicKey)
	if err != nil {
		log.Error(err)
		return activityTick{ready: false, err: err}
	}

	tn := time.Unix(user.NextActivityChangeAt, 0)
	ready := time.Now().After(tn)
	nextDate := randDateInTheFuture()
	nextDateFormatted := time.Unix(nextDate, 0).Format("2006-01-02 3:04PM")
	nowFormatted := time.Now().Format("2006-01-02 3:04PM")

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
			return activityTick{ready: false, err: err}
		}
		log.Infof("%s user's pet changed at: %s (next time: %s)", user.Name, nowFormatted, nextDateFormatted)

		return activityTick{ready: true, next: m.randActivity()}
	}

	return activityTick{ready: false}
}
