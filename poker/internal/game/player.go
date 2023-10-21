package game

import "github.com/rishu/design/poker/pkg/cards"

type Player struct {
	Name string
	Hand []cards.Card
}
