package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Dice interface {
	Roll() int
}

type Player interface {
	GetName() string
	GetPosition() int
	Move(steps int)
}

type Board interface {
	GetFinalPosition(position int) int
}

type Game interface {
	Start()
}

type SimpleDice struct {
}

func NewSimpleDice() *SimpleDice {
	return &SimpleDice{}
}

func (d *SimpleDice) Roll() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(6) + 1
}

type GamePlayer struct {
	name     string
	position int
}

func NewGamePlayer(name string) *GamePlayer {
	return &GamePlayer{
		name:     name,
		position: 0,
	}
}

func (p *GamePlayer) GetName() string {
	return p.name
}

func (p *GamePlayer) GetPosition() int {
	return p.position
}

func (p *GamePlayer) Move(steps int) {
	p.position += steps
}

type GameBoard struct {
	snakes  map[int]int
	ladders map[int]int
}

func NewGameBoard(snakes, ladders map[int]int) *GameBoard {
	return &GameBoard{
		snakes:  snakes,
		ladders: ladders,
	}
}

func (b *GameBoard) GetFinalPosition(position int) int {
	if newPosition, exists := b.snakes[position]; exists {
		return newPosition
	}
	if newPosition, exists := b.ladders[position]; exists {
		return newPosition
	}
	return position
}

type SnakeGame struct {
	players []Player
	board   Board
	dice    Dice
}

func NewSnakeGame(players []Player, board Board, dice Dice) *SnakeGame {
	return &SnakeGame{
		players: players,
		board:   board,
		dice:    dice,
	}
}

func (s *SnakeGame) Start() {
	for {
		for _, player := range s.players {
			steps := s.dice.Roll()
			fmt.Printf("dice rolled %v for player %v \n", steps, player.GetName())
			player.Move(steps)
			finalPos := s.board.GetFinalPosition(player.GetPosition())
			player.Move(finalPos - player.GetPosition())
			fmt.Printf("new position %v for player %v\n", player.GetPosition(), player.GetName())
			if player.GetPosition() >= 100 {
				fmt.Printf("winner is %v \n", player.GetName())
				return
			}
		}
	}
}

func main() {
	player1 := NewGamePlayer("player1")
	player2 := NewGamePlayer("player2")

	snakes := map[int]int{16: 6, 47: 26, 49: 11, 56: 53, 62: 19, 64: 60, 87: 24, 93: 73, 95: 75}
	ladders := map[int]int{1: 38, 4: 14, 9: 31, 21: 42, 28: 84, 36: 44, 51: 67, 71: 91, 80: 100}

	board := NewGameBoard(snakes, ladders)

	dice := NewSimpleDice()

	game := NewSnakeGame([]Player{player1, player2}, board, dice)
	game.Start()

}
