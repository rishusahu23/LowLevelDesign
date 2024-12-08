package main

import "fmt"

const (
	EMPTY = "-"
	X     = "X"
	O     = "O"
)

type Board interface {
	PrintBoard()
	MakeMove(row, col int, symbol string) bool
	UseBooster(row, col int) bool
	CheckWinner() string
	IsDraw() bool
}

type MoveStrategy interface {
	MakeMove(b Board, symbol string) (int, int)
}

type HumanStrategy struct{}

func (h *HumanStrategy) MakeMove(b Board, symbol string) (int, int) {
	var row, col int
	fmt.Printf("Player %s, enter your move (row and column): ", symbol)
	fmt.Scanf("%d %d", &row, &col)
	return row, col
}

type BoosterCommand interface {
	Execute(b Board, row, col int) bool
}

// Concrete Command: EmptyCellBooster
type EmptyCellBooster struct{}

func (e *EmptyCellBooster) Execute(b Board, row, col int) bool {
	return b.UseBooster(row, col)
}

type Player interface {
	GetSymbol() string
	HasBooster() bool
	UseBooster()
	MakeMove(b Board) (int, int)
	UseBoosterCommand(b Board, command BoosterCommand)
}

type TicTacToeBoard struct {
	board [][]string
	moves int
}

func NewTicTacBoard() *TicTacToeBoard {
	return &TicTacToeBoard{
		board: [][]string{
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
			{EMPTY, EMPTY, EMPTY},
		},
		moves: 0,
	}
}

func (b *TicTacToeBoard) PrintBoard() {
	fmt.Println("Current Board:")
	for _, row := range b.board {
		fmt.Println(row)
	}
	fmt.Println()
}

func (b *TicTacToeBoard) MakeMove(row, col int, symbol string) bool {
	if row < 0 || row >= 3 || col < 0 || col >= 3 || b.board[row][col] != EMPTY {
		return false
	}
	b.board[row][col] = symbol
	b.moves++
	return true
}

func (b *TicTacToeBoard) UseBooster(row, col int) bool {
	if row < 0 || row >= 3 || col < 0 || col >= 3 || b.board[row][col] == EMPTY {
		return false
	}
	b.board[row][col] = EMPTY
	return true
}

func (b *TicTacToeBoard) CheckWinner() string {
	for i := 0; i < 3; i++ {
		if b.board[i][0] != EMPTY && b.board[i][0] == b.board[i][1] && b.board[i][1] == b.board[i][2] {
			return b.board[i][0]
		}
		if b.board[0][i] != EMPTY && b.board[0][i] == b.board[1][i] && b.board[1][i] == b.board[2][i] {
			return b.board[0][i]
		}
	}
	if b.board[0][0] != EMPTY && b.board[0][0] == b.board[1][1] && b.board[1][1] == b.board[2][2] {
		return b.board[0][0]
	}
	if b.board[0][2] != EMPTY && b.board[0][2] == b.board[1][1] && b.board[1][1] == b.board[2][0] {
		return b.board[0][2]
	}
	return EMPTY
}

func (b *TicTacToeBoard) IsDraw() bool {
	return b.moves >= 9
}

type TicTacToePlayer struct {
	symbol     string
	hasBooster bool
	strategy   MoveStrategy
}

func NewTicTacToePlayer(symbol string, strategy MoveStrategy) *TicTacToePlayer {
	return &TicTacToePlayer{
		symbol:     symbol,
		hasBooster: true,
		strategy:   strategy,
	}
}

func (p *TicTacToePlayer) GetSymbol() string {
	return p.symbol
}

func (p *TicTacToePlayer) HasBooster() bool {
	return p.hasBooster
}

func (p *TicTacToePlayer) UseBooster() {
	p.hasBooster = false
}

func (p *TicTacToePlayer) MakeMove(b Board) (int, int) {
	return p.strategy.MakeMove(b, p.GetSymbol())
}

func (p *TicTacToePlayer) UseBoosterCommand(b Board, command BoosterCommand) {
	if p.HasBooster() {
		var row, col int
		fmt.Printf("Player %s, enter the row and column to use booster: ", p.GetSymbol())
		fmt.Scanf("%d %d", &row, &col)
		if command.Execute(b, row, col) {
			fmt.Println("Booster used successfully!")
			p.UseBooster()
		} else {
			fmt.Println("Booster failed! Invalid cell.")
		}
	}
}

type Game interface {
	PlayGame()
}

// --- Concrete Game Implementation ---
type TicTacToeGame struct {
	board         Board
	player1       Player
	player2       Player
	currentPlayer Player
}

func NewTicTacToeGame(board Board, player1, player2 Player) *TicTacToeGame {
	return &TicTacToeGame{
		board:         board,
		player1:       player1,
		player2:       player2,
		currentPlayer: player1,
	}
}

func (g *TicTacToeGame) switchPlayer() {
	if g.currentPlayer == g.player1 {
		g.currentPlayer = g.player2
	} else {
		g.currentPlayer = g.player1
	}
}

func (g *TicTacToeGame) PlayGame() {
	boosterCommand := &EmptyCellBooster{}

	for {
		g.board.PrintBoard()

		row, col := g.currentPlayer.MakeMove(g.board)
		if !g.board.MakeMove(row, col, g.currentPlayer.GetSymbol()) {
			fmt.Println("Invalid move! Try again.")
			continue
		}

		if winner := g.board.CheckWinner(); winner != EMPTY {
			g.board.PrintBoard()
			fmt.Printf("Player %s wins!\n", winner)
			break
		}

		if g.board.IsDraw() {
			g.board.PrintBoard()
			fmt.Println("It's a draw!")
			break
		}

		// Ask for booster usage
		if g.currentPlayer.HasBooster() {
			var useBooster string
			fmt.Printf("Player %s, do you want to use your booster? (y/n): ", g.currentPlayer.GetSymbol())
			fmt.Scanf("%s", &useBooster)
			if useBooster == "y" {
				g.currentPlayer.UseBoosterCommand(g.board, boosterCommand)
			}
		}

		g.switchPlayer()
	}
}

// --- Main ---
func main() {
	board := NewTicTacBoard()

	// Factory Method for creating players with different strategies
	player1 := NewTicTacToePlayer(X, &HumanStrategy{})
	player2 := NewTicTacToePlayer(O, &HumanStrategy{})

	game := NewTicTacToeGame(board, player1, player2)
	game.PlayGame()
}
