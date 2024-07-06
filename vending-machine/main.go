package main

import "fmt"

// State represents the interface for different states of the vending machine.
type State interface {
	InsertCoin(vm *VendingMachine)
	SelectItem(vm *VendingMachine)
	DispenseItem(vm *VendingMachine)
}

// ReadyState represents the state when the machine is ready.
type ReadyState struct{}

func (s *ReadyState) InsertCoin(vm *VendingMachine) {
	fmt.Println("Coin inserted.")
	vm.SetState(&InsertCoinState{})
}

func (s *ReadyState) SelectItem(vm *VendingMachine) {
	fmt.Println("Please insert a coin first.")
}

func (s *ReadyState) DispenseItem(vm *VendingMachine) {
	fmt.Println("Please select an item first.")
}

// InsertCoinState represents the state when a coin is inserted.
type InsertCoinState struct{}

func (s *InsertCoinState) InsertCoin(vm *VendingMachine) {
	fmt.Println("Coin already inserted.")
}

func (s *InsertCoinState) SelectItem(vm *VendingMachine) {
	fmt.Println("Item selected.")
	vm.SetState(&SelectItemState{})
}

func (s *InsertCoinState) DispenseItem(vm *VendingMachine) {
	fmt.Println("Please select an item first.")
}

// SelectItemState represents the state when an item is selected.
type SelectItemState struct{}

func (s *SelectItemState) InsertCoin(vm *VendingMachine) {
	fmt.Println("Coin already inserted.")
}

func (s *SelectItemState) SelectItem(vm *VendingMachine) {
	fmt.Println("Item already selected.")
}

func (s *SelectItemState) DispenseItem(vm *VendingMachine) {
	if vm.Count > 0 {
		fmt.Println("Item dispensed.")
		vm.SetState(&DispenseItemState{})
		vm.Count--
		if vm.Count == 0 {
			vm.SetState(&SoldOutState{})
		}
	} else {
		fmt.Println("Item out of stock.")
		vm.SetState(&SoldOutState{})
	}
}

// DispenseItemState represents the state when an item is being dispensed.
type DispenseItemState struct{}

func (s *DispenseItemState) InsertCoin(vm *VendingMachine) {
	fmt.Println("Item already being dispensed.")
}

func (s *DispenseItemState) SelectItem(vm *VendingMachine) {
	fmt.Println("Item already being dispensed.")
}

func (s *DispenseItemState) DispenseItem(vm *VendingMachine) {
	fmt.Println("Item dispensed. Enjoy your purchase!")
	vm.SetState(&ReadyState{})
}

// SoldOutState represents the state when the machine is sold out.
type SoldOutState struct{}

func (s *SoldOutState) InsertCoin(vm *VendingMachine) {
	fmt.Println("Machine is sold out.")
}

func (s *SoldOutState) SelectItem(vm *VendingMachine) {
	fmt.Println("Machine is sold out.")
}

func (s *SoldOutState) DispenseItem(vm *VendingMachine) {
	fmt.Println("Machine is sold out.")
}

// VendingMachine represents the context of the vending machine.
type VendingMachine struct {
	State State
	Count int // Number of items in the machine
}

// SetState sets the state of the vending machine.
func (vm *VendingMachine) SetState(state State) {
	vm.State = state
}

func main() {
	vm := &VendingMachine{
		State: &ReadyState{},
		Count: 2, // Number of items initially
	}

	vm.State.InsertCoin(vm)   // Output: Coin inserted.
	vm.State.InsertCoin(vm)   // Output: Coin already inserted.
	vm.State.SelectItem(vm)   // Output: Item selected.
	vm.State.SelectItem(vm)   // Output: Item already selected.
	vm.State.DispenseItem(vm) // Output: Item dispensed. Enjoy your purchase!
	vm.State.DispenseItem(vm) // Output: Please select an item first.
}
