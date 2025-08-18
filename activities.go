package main

import (
	"math/rand/v2"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/mluna-again/luna/luna"
	"github.com/mluna-again/mlssh/repo"
)

// make this a flag
var DEBUG = false

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

func (m model) scheduleActivityChange(skipSleep bool) tea.Cmd {
	return func() tea.Msg {
		if !skipSleep {
			t := time.Minute
			if DEBUG {
				t = time.Second * 3
			}
			time.Sleep(t)
		}
		// user is not loaded yet (or doesnt have a pk somehow), i should probably handle this better but whatever
		if m.user.publicKey == "" {
			return activityTick{ready: false}
		}

		ctx, cancel := aLittleBit()
		defer cancel()

		user, err := m.queries.GetUser(ctx, m.user.publicKey)
		if err != nil {
			log.Error(err)
			return activityTick{ready: false}
		}

		t := time.Unix(user.NextActivityChangeAt, 0)
		ready := time.Now().After(t)
		if ready {
			nextDate := randDateInTheFuture()
			log.Infof("%s user's pet is scheduled for a change at: %d", user.Name, nextDate)
			_, err := m.queries.UpdateUser(ctx, repo.UpdateUserParams{NextActivityChangeAt: nextDate})
			if err != nil {
				log.Error(err)
				return activityTick{ready: false}
			}

			return activityTick{ready: true, next: m.randActivity()}
		}

		return activityTick{ready: false}
	}
}
