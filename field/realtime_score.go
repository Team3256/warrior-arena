// Copyright 2017 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)
//
// Model representing the current state of the score during a match.

package field

import "github.com/Team3256/warrior-arena/game"

type RealtimeScore struct {
	CurrentScore              game.Score
	Cards                     map[string]string
	FoulsCommitted            bool
	AmplifiedTimeRemainingSec int
}

func NewRealtimeScore() *RealtimeScore {
	return &RealtimeScore{Cards: make(map[string]string)}
}
