package game

import "github.com/rishu/design/poker/internal/deck"

type Game struct {
	Players []Player
	Deck    deck.Deck
}

func NewGame() *Game {
	game := &Game{
		Deck: deck.NewDeck(),
	}
	return game
}

func (g *Game) AddPlayer(name string) {
	player := Player{
		Name: name,
	}
	g.Players = append(g.Players, player)
}

func (g *Game) DealCards(numCards int) {
	for i := 0; i < numCards; i++ {
		for j := range g.Players {
			g.Players[j].Hand = append(g.Players[j].Hand, g.Deck[0])
			g.Deck = g.Deck[1:]
		}
	}
}
