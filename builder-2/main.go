package main

import "fmt"

type House struct {
	Rooms  int
	Window int
	Sofa   int
}

func (h *House) String() string {
	return fmt.Sprintf("%v,%v,%v", h.Rooms, h.Window, h.Sofa)
}

type HouseBuilder interface {
	SetRooms(int) HouseBuilder
	SetWindow(int) HouseBuilder
	SetSofa(int) HouseBuilder
	Build() *House
}

type FirstHouseBuilder struct {
	house House
}

func (f *FirstHouseBuilder) SetRooms(i int) HouseBuilder {
	f.house.Rooms = i
	return f
}

func (f *FirstHouseBuilder) SetWindow(i int) HouseBuilder {
	f.house.Window = i
	return f
}

func (f *FirstHouseBuilder) SetSofa(i int) HouseBuilder {
	f.house.Sofa = i
	return f
}

func (f *FirstHouseBuilder) Build() *House {
	return &f.house
}

var _ HouseBuilder = &FirstHouseBuilder{}

type Director struct {
	builder HouseBuilder
}

func NewDirector(builder HouseBuilder) *Director {
	return &Director{
		builder: builder,
	}
}

func (d *Director) construct() *House {
	return d.builder.SetSofa(3).SetWindow(2).SetRooms(5).Build()
}

func main() {
	d := NewDirector(&FirstHouseBuilder{
		house: House{},
	})
	fmt.Println(d.construct().String())
}
