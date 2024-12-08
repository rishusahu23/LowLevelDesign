package main

import (
	"math/rand"
	"time"
)

type Position struct {
	x, y int
}

type Direction string

const (
	Up    Direction = "UP"
	Down  Direction = "DOWN"
	Left  Direction = "LEFT"
	Right Direction = "RIGHT"
)

type GameBoard struct {
	width  int
	height int
	food   Position
}

func NewGameBoard(h, w int) *GameBoard {
	return &GameBoard{
		width:  w,
		height: h,
	}
}

type Snake struct {
	Body []Position
	Direction
}

func NewSnake(start Position, len int) *Snake {
	body := []Position{start}
	for i := 1; i < len; i++ {
		body = append(body, Position{
			x: start.x - i,
			y: start.y,
		})
	}
	return &Snake{
		Body:      body,
		Direction: Right,
	}
}

func (s *Snake) Move() Position {
	head := s.Body[0]
	var newHead Position

	switch s.Direction {
	case Up:
		newHead = Position{
			x: head.x,
			y: head.y - 1,
		}
	case Down:
		newHead = Position{
			x: head.x,
			y: head.y + 1,
		}
	case Left:
		newHead = Position{
			x: head.x - 1,
			y: head.y,
		}
	case Right:
		newHead = Position{
			x: head.x + 1,
			y: head.y,
		}
	}
	s.Body = append([]Position{newHead}, s.Body[:len(s.Body)-1]...)
	return newHead
}

func (s *Snake) Grow() {
	tail := s.Body[len(s.Body)-1]
	secondTail := s.Body[len(s.Body)-2]

	newTail := Position{
		x: tail.x + (tail.x - secondTail.x),
		y: tail.y + (tail.y - secondTail.y),
	}

	s.Body = append(s.Body, newTail)
}

func (s *Snake) HasCollided(board *GameBoard) bool {
	head := s.Body[0]
	if head.x < 0 || head.y < 0 || head.x >= board.width || head.y >= board.height {
		return true
	}
	for _, pos := range s.Body[1:] {
		if pos == board.food {
			return true
		}
	}
	return false
}

type Food struct {
	pos Position
}

func NewFood(board *GameBoard, snake *Snake) *Food {
	rand.Seed(time.Now().UnixNano())
	for {
		pos := Position{
			x: rand.Intn(board.width),
			y: rand.Intn(board.height),
		}

		isOnSnake := false
		for _, seg := range snake.Body {
			if seg == pos {
				isOnSnake = true
				break
			}
		}
		if !isOnSnake {
			return &Food{
				pos: pos,
			}
		}
	}
	return nil
}
