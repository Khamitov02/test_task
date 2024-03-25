package flood

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"
)

type Container struct {
	messageTimesByUserID map[int64][]time.Time
	N                    int64 //nseconds
	K                    int64 //max calls
	mx                   sync.Mutex
}

func (cntr *Container) UpdateInfo(userID int64, t time.Time) []time.Time {
	cntr.mx.Lock()
	defer cntr.mx.Unlock()
	if _, ok := cntr.messageTimesByUserID[userID]; !ok {
		cntr.messageTimesByUserID[userID] = []time.Time{t}
	} else {
		cntr.messageTimesByUserID[userID] = append(cntr.messageTimesByUserID[userID], t)
	}
	return cntr.messageTimesByUserID[userID]
}

// Returns a new Container
func NewFlood(N int64, K int64) Container {
	return Container{
		messageTimesByUserID: make(map[int64][]time.Time),
		N:                    N,
		K:                    K,
	}
}

// Check checks if a user can make a call based on the time of their previous calls.
// It takes a context and the user ID as input and returns a boolean indicating whether the user can make a call and an error if any.
func (cntr *Container) Check(ctx context.Context, userID int64) (canMakeCall bool, err error) {
	// Defer a function that checks if the context has been cancelled and sets the error accordingly
	defer func() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			err = ctxErr
		}
	}()

	// Get the current time
	now := time.Now()

	// Add the current time for the user to the slice of times for the user
	userTimes := cntr.UpdateInfo(userID, now)

	// Sort the slice of times for the user by time
	sort.Slice(userTimes, func(i, j int) bool {
		return userTimes[i].Before(userTimes[j])
	})

	// Check if the number of calls made by the user within the last K seconds is greater than the maximum allowed calls
	switch {
	case cntr.Count(userTimes, now) > cntr.K:
		canMakeCall = false
	case cntr.Count(userTimes, now) < cntr.K:
		canMakeCall = true
	}
	return
}

func (cntr *Container) Count(usersTime []time.Time, t time.Time) int64 {
	// Binary search for the index of the first element in usersTime that is greater than t
	index, _ := slices.BinarySearchFunc(usersTime, t.Add(time.Duration(-cntr.N)*time.Second), func(a, b time.Time) int {
		if b.Before(a) {
			return 1
		} else if a == b {
			return 0
		} else {
			return -1
		}
	})
	// Return the count of elements greater than t
	return int64(len(usersTime) - index)
}
