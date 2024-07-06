package main

type Context struct {
	Strategy DiceStrategy
}

func NewContext(strategy DiceStrategy) *Context {
	return &Context{
		Strategy: strategy,
	}
}

func (c *Context) SetStrategy(strategy DiceStrategy) {
	c.Strategy = strategy
}

func (c *Context) GetDiceValue() int {
	return c.Strategy.GetDiceValue()
}
