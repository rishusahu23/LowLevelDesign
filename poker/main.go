package main

import (
	"fmt"
	game2 "github.com/rishu/design/poker/internal/game"
)

func main() {
	game := game2.NewGame()

	game.AddPlayer("Player 1")
	game.AddPlayer("Player 2")

	game.DealCards(2)

	for _, player := range game.Players {
		fmt.Printf("%s's hand: %v\n", player.Name, player.Hand)
	}
}
