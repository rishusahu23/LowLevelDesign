package main

import (
	"fmt"
	"time"
)

func worker(ch chan bool) {
	time.Sleep(1 * time.Second)
	fmt.Println("hello")
	ch <- true
}

func main() {
	msg1 := make(chan string)
	msg2 := make(chan string)
	go func() {
		time.Sleep(1 * time.Second)
		msg1 <- "one"
	}()
	go func() {
		//time.Sleep(2 * time.Second)
		msg2 <- "one"
	}()
	for i := 0; i < 2; i++ {
		select {
		case val := <-msg1:
			fmt.Println(val)
		case val := <-msg2:
			fmt.Println(val)
		}
	}
	fmt.Println("hi")
}
