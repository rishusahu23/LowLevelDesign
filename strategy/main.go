package main

import "fmt"

type Strategy interface {
	ExecuteStrategy(int, int) int
}

type Strategy1 struct{}

func (s *Strategy1) ExecuteStrategy(a int, b int) int {
	return a + b
}

type Strategy2 struct{}

func (s *Strategy2) ExecuteStrategy(a int, b int) int {
	return a - b
}

type Context struct {
	Strategy Strategy
}

func (c *Context) SetStrategy(strategy Strategy) *Context {
	return &Context{
		Strategy: strategy,
	}
}

func (c *Context) ExecuteStrategy(a, b int) int {
	return c.Strategy.ExecuteStrategy(a, b)
}

func main() {
	context := &Context{}
	str := context.SetStrategy(&Strategy1{})
	fmt.Println(str.ExecuteStrategy(1, 2))
	str = context.SetStrategy(&Strategy2{})
	fmt.Println(str.ExecuteStrategy(1, 2))
}
