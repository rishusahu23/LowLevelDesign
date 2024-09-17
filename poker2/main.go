package main

import (
	"math/rand"
	"time"
)

type ICard interface {
	GetSuit() string
	GetValue() string
}

type IDeck interface {
	Shuffle()
	Draw() ICard
}

type IPlayer interface {
	GetName() string
	GetHand() []ICard
	AddCard(card ICard)
	RemoveCard(card ICard)
}

type Game interface {
	AddPlayer(player IPlayer)
	RemovePlayer(player IPlayer)
	Start()
	End()
	PlayRound()
}

type Card struct {
	suit, value string
}

func (c *Card) GetSuit() string {
	return c.suit
}

func (c *Card) GetValue() string {
	return c.value
}

func NewCard(suit string, value string) *Card {
	return &Card{suit: suit, value: value}
}

type Deck struct {
	cards []*Card
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Draw() ICard {
	if len(d.cards) == 0 {
		return nil
	}
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card
}

func NewDeck() *Deck {
	cards := make([]*Card, 0)
	suits := []string{"a", "b", "c", "d"}
	values := []string{"2", "3", "3", "4"}
	for _, suit := range suits {
		for _, value := range values {
			cards = append(cards, &Card{
				suit:  suit,
				value: value,
			})
		}
	}
	return &Deck{
		cards: cards,
	}
}
