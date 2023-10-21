package store

import "github.com/rishu/design/sports_subscription/game"

var gameInstance *game.Game

func SetGame(g *game.Game) {
	if gameInstance == nil {
		gameInstance = g
	}
}

func GetGameInstance() *game.Game {
	return gameInstance
}
