package main

import (
	"errors"
	"fmt"
	"sync"
)

type Cell interface {
	GetMark() string
	SetMark(mark string)
}

type CellImpl struct {
	row, col int
	mark     string
}

func NewCellImpl(row int, col int) *CellImpl {
	return &CellImpl{row: row, col: col}
}

func (c *CellImpl) GetMark() string {
	return c.mark
}

func (c *CellImpl) SetMark(mark string) {
	c.mark = mark
}

type Board interface {
	Display()
	IsFull() bool
	IsWinner(mark string) bool
	GetCell(row, col int) Cell
	SetCell(row, col int, mark string) error
}

type BoardImpl struct {
	cells [3][3]Cell
}

func (b *BoardImpl) Display() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.GetCell(i, j).GetMark() != "" {
				fmt.Print(b.GetCell(i, j).GetMark())
			} else {
				fmt.Print("-")
			}
		}
		fmt.Println()
	}
}

func (b *BoardImpl) IsFull() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.GetCell(i, j).GetMark() == "" {
				return false
			}
		}
	}
	return true
}

func (b *BoardImpl) IsWinner(mark string) bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if b.cells[i][0].GetMark() == mark && b.cells[i][1].GetMark() == mark && b.cells[i][2].GetMark() == mark {
			return true
		}
	}
	// Check columns
	for j := 0; j < 3; j++ {
		if b.cells[0][j].GetMark() == mark && b.cells[1][j].GetMark() == mark && b.cells[2][j].GetMark() == mark {
			return true
		}
	}
	// Check diagonals
	if b.cells[0][0].GetMark() == mark && b.cells[1][1].GetMark() == mark && b.cells[2][2].GetMark() == mark {
		return true
	}
	if b.cells[0][2].GetMark() == mark && b.cells[1][1].GetMark() == mark && b.cells[2][0].GetMark() == mark {
		return true
	}
	return false
}

func (b *BoardImpl) GetCell(row, col int) Cell {
	return b.cells[row][col]
}

func (b *BoardImpl) SetCell(row, col int, mark string) error {
	if row < 0 {
		return errors.New("error")
	}
	b.cells[row][col].SetMark(mark)
	return nil
}

func NewBoardImpl() *BoardImpl {
	b := &BoardImpl{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			b.cells[i][j] = NewCellImpl(i, j)
		}
	}
	return b
}

type Player interface {
	GetMove(board Board) (int, int)
	GetMark() string
}

type Human struct {
	mark string
}

func (h *Human) GetMove(impl Board) (int, int) {
	var col, row int
	fmt.Println("enter row and column")
	fmt.Scan(&row, &col)
	return row, col
}

func (h *Human) GetMark() string {
	return h.mark
}

type Computer struct {
	mark string
}

func (c *Computer) GetMove(impl Board) (int, int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if impl.GetCell(i, j).GetMark() == "" {
				return i, j
			}
		}
	}
	return -1, -1
}

func (c *Computer) GetMark() string {
	return c.mark
}

type PlayerType int

const (
	HumanType PlayerType = iota
	ComputerType
)

type PlayerFactory interface {
	CreatePlayer(playerType PlayerType, mark string) Player
}

type PlayerFactImpl struct {
}

func (f *PlayerFactImpl) CreatePlayer(playerType PlayerType, mark string) Player {
	switch playerType {
	case HumanType:
		return &Human{
			mark: mark,
		}
	case ComputerType:
		return &Computer{
			mark: mark,
		}
	}
	return nil
}

type Game interface {
	Play()
	SwitchPlayer()
}

type GameImpl struct {
	board         Board
	player1       Player
	player2       Player
	currentPlayer Player
}

func (g *GameImpl) Play() {
	for {
		g.board.Display()
		row, col := g.currentPlayer.GetMove(g.board)
		g.board.SetCell(row, col, g.currentPlayer.GetMark())
		if g.board.IsWinner(g.currentPlayer.GetMark()) {
			g.board.Display()
			fmt.Println("player won ", g.currentPlayer.GetMark())
			return
		}
		if g.board.IsFull() {
			g.board.Display()
			fmt.Println("full board")
			return
		}
		g.SwitchPlayer()
	}
}

func (g *GameImpl) SwitchPlayer() {
	if g.currentPlayer == g.player1 {
		g.currentPlayer = g.player2
	} else {
		g.currentPlayer = g.player1
	}
}

func NewGameImpl(player1 Player, player2 Player, currentPlayer Player) *GameImpl {
	once.Do(func() {
		instance = &GameImpl{
			board:         NewBoardImpl(),
			player1:       player1,
			player2:       player2,
			currentPlayer: currentPlayer,
		}
	})
	return instance
}

var instance *GameImpl
var once sync.Once
