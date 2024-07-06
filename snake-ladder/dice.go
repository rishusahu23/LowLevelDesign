package main

import "math/rand"

type DiceStrategy interface {
	GetDiceValue() int
}

type Strategy1 struct{}

func (s *Strategy1) GetDiceValue() int {
	return rand.Intn(6) + 1
}

type Strategy2 struct{}

func (s *Strategy2) GetDiceValue() int {
	return rand.Intn(4) + 1
}
