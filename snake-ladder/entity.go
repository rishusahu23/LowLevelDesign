package main

import (
	"fmt"
	"sync"
)

type Entities struct {
	Snakes  map[int]int
	Ladders map[int]int
	Players map[int]string
}

var (
	entitiesInstance *Entities
	once             sync.Once
)

func GetInstance() *Entities {
	once.Do(func() {
		entitiesInstance = &Entities{
			Snakes:  make(map[int]int),
			Ladders: make(map[int]int),
			Players: make(map[int]string),
		}
	})
	return entitiesInstance
}

func (e *Entities) SetSnake(ss, ee int) {
	e.Snakes[ss] = ee
}

func (e *Entities) GetSnakes() map[int]int {
	return e.Snakes
}

func (e *Entities) PrintSnakes() {
	fmt.Println("Printing Snakes")
	for key, val := range e.Snakes {
		fmt.Println(key, " ", val)
	}
	fmt.Println("Done")
}

func (e *Entities) SetLadder(ss, ee int) {
	e.Ladders[ss] = ee
}

func (e *Entities) GetLadders() map[int]int {
	return e.Ladders
}

func (e *Entities) PrintLadders() {
	fmt.Println("Printing Ladders")
	for key, val := range e.Ladders {
		fmt.Println(key, " ", val)
	}
	fmt.Println("Done")
}

func (e *Entities) SetPlayer(index int, player string) {
	e.Players[index] = player
}

func (e *Entities) GetPlayers() map[int]string {
	return e.Players
}
