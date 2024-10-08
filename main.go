package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type GameStatus string

const (
	InProgress GameStatus = "IN_PROGRESS"
	XWon       GameStatus = "X_WON"
	OWon       GameStatus = "O_WON"
	Draw       GameStatus = "DRAW"
)

type Player struct {
	ID   string
	Name string
}

type Game struct {
	ID        string
	PlayerX   Player
	PlayerO   Player
	Status    GameStatus
	Board     [3][3]string
	CreatedAt time.Time
	UpdatedAt time.Time
	Winner    *Player
	Moves     []Move
}

type Move struct {
	ID        string
	GameId    string
	Player    Player
	PositionX int
	PositionY int
	MoveOrder int
	CreatedAt time.Time
}

type GameService struct {
	games map[string]*Game
	mu    sync.RWMutex
}

// generateUUID generates a dummy UUID for demonstration
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (s *GameService) StartGame(playerXId, playerXName, playerOId, playerOName string) (*Game, error) {
	playerX := Player{
		ID:   playerXId,
		Name: playerXName,
	}
	playerO := Player{
		ID:   playerOId,
		Name: playerOName,
	}

	game := &Game{
		ID:        generateUUID(),
		PlayerX:   playerX,
		PlayerO:   playerO,
		Status:    InProgress,
		Board:     [3][3]string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.mu.Lock()
	s.games[game.ID] = game
	s.mu.Unlock()

	return game, nil
}

func (g *Game) GetPlayer(playerID string) *Player {
	if g.PlayerX.ID == playerID {
		return &g.PlayerX
	} else if g.PlayerO.ID == playerID {
		return &g.PlayerO
	}
	return nil
}

func (s *GameService) MakeMove(gameId, playerId string, x, y int) (*Game, error) {
	s.mu.RLock()
	game, exists := s.games[gameId]
	s.mu.RUnlock()
	if !exists {
		return nil, errors.New("game not found")
	}
	if game.Status != InProgress {
		return nil, errors.New("game is not in progress")
	}
	if game.Board[x][y] != "" {
		return nil, errors.New("position already taken")
	}
	var mark string
	if game.PlayerX.ID == playerId {
		mark = "X"
	} else if game.PlayerO.ID == playerId {
		mark = "O"
	} else {
		return nil, errors.New("invalid player")
	}
	game.Board[x][y] = mark
	game.UpdatedAt = time.Now()

	move := Move{
		ID:        generateUUID(),
		GameId:    gameId,
		Player:    *game.GetPlayer(playerId),
		PositionX: x,
		PositionY: y,
		MoveOrder: len(game.Moves) + 1,
		CreatedAt: time.Now(),
	}
	game.Moves = append(game.Moves, move)

	if checkWin(game) {
		game.Status = XWon
		if mark == "O" {
			game.Status = OWon
		}
		game.Winner = game.GetPlayer(playerId)
	} else if checkDraw(game) {
		game.Status = Draw
	}
	s.mu.Lock()
	s.games[gameId] = game
	s.mu.Unlock()

	return game, nil
}

func checkWin(game *Game) bool {
	// Logic to check rows, columns, diagonals for a win
	return false
}

// checkDraw checks if the game is a draw
func checkDraw(game *Game) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if game.Board[i][j] == "" {
				return false
			}
		}
	}
	return true
}

func main()  {
	gameService := &GameService{
		games: make(map[string]*Game),
	}
	game, err := gameService.StartGame("player1", "Alice", "player2", "Bob")
	if err != nil {
		fmt.Println("Error starting game:", err)
		return
	}
	fmt.Println("Game started! Game ID:", game.ID)

}
