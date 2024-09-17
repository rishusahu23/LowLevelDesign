package main

import "fmt"

func main() {
	msg := make(chan string, 1)

	msg <- "hello"

	fmt.Println(<-msg)
}
