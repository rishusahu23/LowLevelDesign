package main

import "fmt"

type Product struct {
	part1 string
	part2 int
}

type Builder interface {
	SetProduct1(string2 string)
	SetProduct2(int2 int)
	GetProduct() *Product
}

type ConcreteBuilder struct {
	Product *Product
}

func NewConcreteBuilder(product *Product) *ConcreteBuilder {
	return &ConcreteBuilder{
		Product: product,
	}
}

func (c *ConcreteBuilder) SetProduct1(part1 string) {
	c.Product.part1 = part1
}

func (c *ConcreteBuilder) SetProduct2(part2 int) {
	c.Product.part2 = part2
}

func (c *ConcreteBuilder) GetProduct() *Product {
	return c.Product
}

type Director struct {
	builder Builder
}

func NewDirector(builder Builder) *Director {
	return &Director{
		builder: builder,
	}
}

func (d *Director) Construct(string2 string, int2 int) *Product {
	d.builder.SetProduct1(string2)
	d.builder.SetProduct2(int2)
	return d.builder.GetProduct()
}

func main() {
	p1 := NewConcreteBuilder(&Product{})
	d1 := NewDirector(p1)
	d1.Construct("rishu", 23)
	fmt.Println(d1)
}
