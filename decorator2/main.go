package main

import "fmt"

type Tea interface {
	GetDescription() string
	GetCost() float64
}

type MasalaChai struct {
}

func (m *MasalaChai) GetDescription() string {
	return "masala chai"
}

func (m *MasalaChai) GetCost() float64 {
	return 1
}

type GreenTea struct {
}

func (g *GreenTea) GetDescription() string {
	return "green tea"
}

func (g *GreenTea) GetCost() float64 {
	return 2
}

type CondimentDecorator struct {
	tea Tea
}

type Milk struct {
	CondimentDecorator
}

func (g *Milk) GetDescription() string {
	return g.tea.GetDescription() + " milk"
}

func (g *Milk) GetCost() float64 {
	return g.tea.GetCost() + 1
}

type Ginger struct {
	CondimentDecorator
}

func (g *Ginger) GetDescription() string {
	return g.tea.GetDescription() + " ginger"
}

func (g *Ginger) GetCost() float64 {
	return g.tea.GetCost() + 1
}

func main() {
	masalaChai := &MasalaChai{}
	fmt.Println(masalaChai.GetDescription(), masalaChai.GetCost())
	ginger := &Ginger{
		CondimentDecorator: CondimentDecorator{
			tea: masalaChai,
		},
	}
	milk := &Milk{
		CondimentDecorator: CondimentDecorator{
			tea: ginger,
		},
	}
	fmt.Println(milk.GetDescription(), milk.GetCost())
}
