package main

import "errors"

type Product struct {
	Id    string
	Price float64
	Name  string
	Stock int
}

type Inventory struct {
	Products map[string]*Product
}

func (i *Inventory) AddProduct(p *Product) {
	i.Products[p.Id] = p
}

func (i *Inventory) GetProduct(id string) (*Product, error) {
	product, ok := i.Products[id]
	if !ok {
		return nil, errors.New("product not found")
	}
	return product, nil
}
