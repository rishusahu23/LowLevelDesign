package main

import (
	"github.com/rishu/design/sports_subscription/database"
	"github.com/rishu/design/sports_subscription/game"
	"github.com/rishu/design/sports_subscription/subscriber"
)

func main() {
	db := database.NewInMemoryDB()
	g := game.Game{
		ID:      "1",
		Score:   "0-0",
		Status:  "Not started",
		History: []string{},
	}

	db.AddGame(g)

	subscriber1 := &subscriber.GameSubscriber{
		ID: "Subscriber1",
	}

	db.GetGame("1").RegisterObserver(subscriber1)
	g.Score = "1-0"
	g.Status = "In Progress"
	g.History = append(g.History, "Team A Scored")
	db.GetGame("1").Game = g
	db.GetGame("1").NotifyObservers()
}
