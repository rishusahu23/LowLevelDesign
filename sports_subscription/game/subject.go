package game

type Subject interface {
	RegisterObserver(o Observer)
	RemoveObserver(o Observer)
	NotifyObservers()
}

type SubjectGame struct {
	Game      Game
	Observers map[Observer]struct{}
}

func NewSubjectGame(game Game) *SubjectGame {
	return &SubjectGame{
		Game:      game,
		Observers: make(map[Observer]struct{}),
	}
}

func (gs *SubjectGame) RegisterObserver(o Observer) {
	gs.Observers[o] = struct{}{}
}

func (gs *SubjectGame) RemoveObserver(o Observer) {
	delete(gs.Observers, o)
}

func (gs *SubjectGame) NotifyObservers() {
	for o := range gs.Observers {
		o.Update(gs.Game)
	}
}
