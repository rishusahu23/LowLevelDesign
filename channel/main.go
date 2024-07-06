package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)

	go func() {
		//time.Sleep(3 * time.Second)
		ch1 <- "ch1"
	}()

	select {
	case msg := <-ch1:
		fmt.Println(msg)
	case <-time.After(2 * time.Second):
		fmt.Println("timeout")
	default:
		fmt.Println("default")
	}
	fmt.Println("done")
	//ch2 := make(chan string)
	//go func() {
	//	time.Sleep(2 * time.Second)
	//	ch2 <- "ch2"
	//}()
	//for i := 0; i < 2; i++ {
	//	select {
	//	case msg := <-ch1:
	//		fmt.Println(msg)
	//	case msg := <-ch2:
	//		fmt.Println(msg)
	//
	//	}
	//}
}
