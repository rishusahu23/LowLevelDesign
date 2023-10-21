package deck

import (
	"github.com/rishu/design/poker/pkg/cards"
	"math/rand"
	"time"
)

type Deck []cards.Card

func NewDeck() Deck {
	var deck Deck
	for _, suit := range []cards.Suit{cards.Spades, cards.Hearts, cards.Diamonds, cards.Clubs} {
		for _, rank := range []cards.Rank{cards.Two, cards.Three, cards.Four, cards.Five, cards.Six, cards.Seven, cards.Eight, cards.Nine, cards.Ten, cards.Jack, cards.Queen, cards.King, cards.Ace} {
			deck = append(deck, cards.Card{Suit: suit, Rank: rank})
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
