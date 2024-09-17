package main

import "fmt"

func main() {
	inventory := &Inventory{
		Products: map[string]*Product{
			"A1": {Id: "A1", Name: "Soda", Price: 1.25, Stock: 10},
			"B2": {Id: "B2", Name: "Chips", Price: 1.00, Stock: 5},
		},
	}
	context := &Context{
		state:     &IdleState{},
		inventory: inventory,
		balance:   0,
	}

	// Simulate interactions
	err := context.SelectProduct("A1")
	if err != nil {
		fmt.Println(err)
	}
	err = context.InsertMoney(1.00)
	if err != nil {
		fmt.Println(err)
	}
	err = context.InsertMoney(0.25)
	if err != nil {
		fmt.Println(err)
	}
	err = context.DispenseProduct()
	if err != nil {
		fmt.Println(err)
	}

	// Simulate cancellation
	err = context.Cancel()
	if err != nil {
		fmt.Println(err)
	}

	// Simulate refund
	err = context.Refund()
	if err != nil {
		fmt.Println(err)
	}
}
