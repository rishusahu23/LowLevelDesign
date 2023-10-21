package database

import "github.com/rishu/design/sports_subscription/game"

type InMemoryDB struct {
	Games map[string]*game.SubjectGame
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		Games: make(map[string]*game.SubjectGame),
	}
}

func (db *InMemoryDB) AddGame(g game.Game) {
	subject := game.NewSubjectGame(g)
	db.Games[g.ID] = subject
}

func (db *InMemoryDB) GetGame(id string) *game.SubjectGame {
	return db.Games[id]
}
