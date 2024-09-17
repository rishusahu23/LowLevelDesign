package main

import "fmt"

type PaymentStrategy interface {
	Pay(amount float64)
}

type Card struct {
	cardNumber int
}

func (c *Card) Pay(amount float64) {
	fmt.Println("pay via card")
}

type PayPal struct {
}

func (p *PayPal) Pay(amount float64) {
	fmt.Println("pay via paypal")
}

type context struct {
	strategy PaymentStrategy
}

func (c *context) SetStrategy(strategy PaymentStrategy) {
	c.strategy = strategy
}

func (c *context) Pay(float642 float64) {
	c.strategy.Pay(float642)
}

func main() {
	card := &Card{}
	ctx := &context{}
	ctx.SetStrategy(card)
	ctx.Pay(7)
	ctx.SetStrategy(&PayPal{})
	ctx.Pay(78)
}
