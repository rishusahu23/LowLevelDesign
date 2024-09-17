package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Enums for Game Status
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
	GameID    string
	Player    Player
	PositionX int
	PositionY int
	MoveOrder int
	CreatedAt time.Time
}

// In-memory storage for games
type GameService struct {
	mu    sync.RWMutex
	games map[string]*Game
}

// StartGame initializes a new game between two players and stores it in memory.
func (s *GameService) StartGame(playerXID, playerXName, playerOID, playerOName string) (*Game, error) {
	playerX := Player{ID: playerXID, Name: playerXName}
	playerO := Player{ID: playerOID, Name: playerOName}

	game := &Game{
		ID:        generateUUID(),
		PlayerX:   playerX,
		PlayerO:   playerO,
		Status:    InProgress,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Board:     [3][3]string{},
	}

	// Store game in memory
	s.mu.Lock()
	s.games[game.ID] = game
	s.mu.Unlock()

	return game, nil
}

// MakeMove allows a player to make a move in the game.
func (s *GameService) MakeMove(gameID, playerID string, x, y int) (*Game, error) {
	s.mu.RLock()
	game, exists := s.games[gameID]
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

	// Identify the player's mark (X or O)
	var mark string
	if game.PlayerX.ID == playerID {
		mark = "X"
	} else if game.PlayerO.ID == playerID {
		mark = "O"
	} else {
		return nil, errors.New("invalid player")
	}

	// Update the board
	game.Board[x][y] = mark
	game.UpdatedAt = time.Now()

	move := Move{
		ID:        generateUUID(),
		GameID:    gameID,
		Player:    *game.GetPlayer(playerID),
		PositionX: x,
		PositionY: y,
		MoveOrder: len(game.Moves) + 1,
		CreatedAt: time.Now(),
	}
	game.Moves = append(game.Moves, move)

	// Check if the game is won or drawn
	if checkWin(game) {
		game.Status = XWon
		if mark == "O" {
			game.Status = OWon
		}
		game.Winner = game.GetPlayer(playerID)
	} else if checkDraw(game) {
		game.Status = Draw
	}

	// Update the game state in memory
	s.mu.Lock()
	s.games[gameID] = game
	s.mu.Unlock()

	return game, nil
}

// Helper function to get a player by ID
func (g *Game) GetPlayer(playerID string) *Player {
	if g.PlayerX.ID == playerID {
		return &g.PlayerX
	} else if g.PlayerO.ID == playerID {
		return &g.PlayerO
	}
	return nil
}

// checkWin checks if a player has won the game
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

// generateUUID generates a dummy UUID for demonstration
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// Main function
func main() {
	// Initialize GameService with an in-memory store
	gameService := &GameService{
		games: make(map[string]*Game),
	}

	// Start a new game
	game, err := gameService.StartGame("player1", "Alice", "player2", "Bob")
	if err != nil {
		fmt.Println("Error starting game:", err)
		return
	}

	fmt.Println("Game started! Game ID:", game.ID)

	// Make some moves
	_, err = gameService.MakeMove(game.ID, "player1", 0, 0) // Alice (X) moves
	if err != nil {
		fmt.Println("Error making move:", err)
	}

	_, err = gameService.MakeMove(game.ID, "player2", 1, 1) // Bob (O) moves
	if err != nil {
		fmt.Println("Error making move:", err)
	}

	_, err = gameService.MakeMove(game.ID, "player1", 0, 1) // Alice (X) moves
	if err != nil {
		fmt.Println("Error making move:", err)
	}

	// Print the current board state
	printBoard(game.Board)

	// Continue playing until the game ends...
}

// Helper function to print the board
func printBoard(board [3][3]string) {
	for _, row := range board {
		fmt.Println(row)
	}
}
