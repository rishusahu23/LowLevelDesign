package main

import "fmt"

func main() {
	entities := GetInstance()
	var snake, ladder, player int
	fmt.Println("get number of snake")
	fmt.Scanln(&snake)
	for i := 0; i < snake; i++ {
		var s, e int
		fmt.Scanln(&s, &e)
		entities.SetSnake(s, e)

	}
	fmt.Println("get number of ladder")
	fmt.Scanln(&ladder)
	for i := 0; i < ladder; i++ {
		var s, e int
		fmt.Scanln(&s, &e)
		entities.SetLadder(s, e)
	}
	fmt.Println("get number of player")
	fmt.Scanln(&player)
	for i := 0; i < player; i++ {
		var name string
		fmt.Scanln(&name)
		entities.SetPlayer(i, name)
	}
	entities.PrintSnakes()
	entities.PrintLadders()
}
