package main

import "fmt"

type Character interface {
	Introduce()
}

type Elf struct {
}

func (e *Elf) Introduce() {
	fmt.Println("elf")
}

type Orc struct {
}

func (e *Orc) Introduce() {
	fmt.Println("Orc")
}

type Weapon interface {
	Hit()
}

type Bow struct {
}

func (b *Bow) Hit() {
	fmt.Println("bow")
}

type Axe struct {
}

func (b *Axe) Hit() {
	fmt.Println("axe")
}

type CharacterFactory interface {
	CreateCharacter() Character
	CreateWeapon() Weapon
}

type ElfFactory struct {
}

func (e *ElfFactory) CreateCharacter() Character {
	return &Elf{}
}

func (e *ElfFactory) CreateWeapon() Weapon {
	return &Bow{}
}

type OrcFactory struct {
}

func (e *OrcFactory) CreateCharacter() Character {
	return &Orc{}
}

func (e *OrcFactory) CreateWeapon() Weapon {
	return &Axe{}
}

func main() {
	elfFact := &ElfFactory{}
	elfFact.CreateWeapon().Hit()
	elfFact.CreateCharacter().Introduce()
}
