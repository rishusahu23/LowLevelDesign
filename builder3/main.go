package main

import "fmt"

type House struct {
	foundation string
	structure string
	roof string
	interior string
}

type HouseBuilder interface {
	SetFoundation()
	SetStructure()
	SetRoof()
	SetInterior()
	GetHouse() House
}

type ConcreteHouseBuilder struct {
	house House
}

func (c ConcreteHouseBuilder) SetFoundation() {
	c.house.foundation = "found"
}

func (c ConcreteHouseBuilder) SetStructure() {
	c.house.structure = "struct"
}

func (c ConcreteHouseBuilder) SetRoof() {
	c.house.roof = "roof"
}

func (c ConcreteHouseBuilder) SetInterior() {
	c.house.interior = "int"
}

func (c ConcreteHouseBuilder) GetHouse() House {
	return c.house
}

func NewConcreteHouseBuilder() *ConcreteHouseBuilder {
	return &ConcreteHouseBuilder{}
}

type Director struct {
	builder HouseBuilder
}

func NewDirector(builder HouseBuilder) *Director {
	return &Director{
		builder: builder,
	}
}

func (d *Director) Construct() {
	d.builder.SetInterior()
	d.builder.SetFoundation()
	d.builder.SetStructure()
	d.builder.SetRoof()
	d.builder.GetHouse()
}

func main() {
	builder := NewConcreteHouseBuilder()
	director := NewDirector(builder)
	director.Construct()
	house := builder.GetHouse()
	fmt.Println(house)
}