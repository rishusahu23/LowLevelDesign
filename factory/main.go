package main

import "fmt"

type Character interface {
	PerformAbility() string
}

type Warrior struct {
}

func (w Warrior) PerformAbility() string {
	return "i attack"
}

var _ Character = &Warrior{}

type King struct {
}

func (w King) PerformAbility() string {
	return "i take descision"
}

var _ Character = &King{}

type CharacterFactory struct {
}

func (c *CharacterFactory) GetCharacter(chType string) Character {
	switch chType {
	case "king":
		return &King{}
	default:
		return &Warrior{}
	}
}

func main() {
	chFact := CharacterFactory{}
	char := chFact.GetCharacter("king")
	fmt.Println(char.PerformAbility())
	char = chFact.GetCharacter("kking")
	fmt.Println(char.PerformAbility())

}