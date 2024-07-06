package main

import "fmt"

type PlaySnakeAndLadder struct {
	playerHistory   map[string]PairPosition
	playerLatestPos map[string]int
	entities        *Entities
	dice            DiceStrategy
}

func NewPlaySnakeAndLadder(strategy DiceStrategy) *PlaySnakeAndLadder {
	return &PlaySnakeAndLadder{
		playerHistory:   make(map[string]PairPosition),
		playerLatestPos: make(map[string]int),
		entities:        GetInstance(),
		dice:            strategy,
	}
}

func (s *PlaySnakeAndLadder) isValidPos(endPos int) bool {
	return endPos <= 100
}

func (s *PlaySnakeAndLadder) PlayGame() string {
	i := -1
	for {
		i = (i + 1) % len(s.entities.GetPlayers())
		playerName := s.entities.GetPlayers()[i]
		str := playerName
		diceValue := s.dice.GetDiceValue()
		endPos := s.playerLatestPos[playerName] + diceValue
		var s1 string
		if s.isValidPos(endPos) {
			str += fmt.Sprintf("rolled a dice %d", diceValue)
			str += fmt.Sprintf("and moved from %d", s.playerLatestPos[playerName])
			if val, ok := s.entities.GetSnakes()[endPos]; ok {
				s1 = " after Snake dinner"
				s.playerLatestPos[playerName] = val
			} else if val, ok := s.entities.GetLadders()[endPos]; ok {
				s1 = " after Ladder ride "
				s.playerLatestPos[playerName] = val
			} else {
				s.playerLatestPos[playerName] = endPos
			}
			str += fmt.Sprintf("to %d", s.playerLatestPos[playerName])
			str += s1
		}

	}
}

func (s *PlaySnakeAndLadder) initialise() {
	players := s.entities.GetPlayers()
	for i := 0; i < len(players); i++ {
		s.playerLatestPos[players[i]] = 0
	}
}
