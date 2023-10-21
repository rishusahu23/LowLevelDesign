package subscriber

import (
	"fmt"
	"github.com/rishu/design/sports_subscription/game"
)

type GameSubscriber struct {
	ID string
}

func (gs *GameSubscriber) Update(game game.Game) {
	fmt.Printf("Subscriber %s: Game %s - Score: %s, Status: %s, History: %v\n",
		gs.ID, game.ID, game.Score, game.Status, game.History)
}
