package main

import "fmt"

type Account struct {
	balance  float64
	userName string
}

func (a *Account) Balance() float64 {
	return a.balance
}

func (a *Account) SetBalance(balance float64) {
	a.balance = balance
}

func (a *Account) UserName() string {
	return a.userName
}

func (a *Account) SetUserName(userName string) {
	a.userName = userName
}

func NewAccount(balance float64, userName string) *Account {
	return &Account{balance: balance, userName: userName}
}

type Card struct {
	account *Account
}

func NewCard(account *Account) *Card {
	return &Card{account: account}
}

func (c *Card) GetAccount() *Account {
	return c.account
}

type StateEnum int

const (
	WELCOME StateEnum = iota
	INSERT_CARD
	TRANSACTION
	REMOVE_CARD
	THANK_YOU
)

type StateHandler interface {
	HandleFlow()
}

type WelcomeState struct {
}

func (w *WelcomeState) HandleFlow() {
	fmt.Println("Hello! Welcome to XYZ atm")
}

type InsertCardState struct{}

func (s *InsertCardState) HandleFlow() {
	fmt.Println("Insert your card")
}

type TransactionState struct {
	card *Card
}

func NewTransactionState(card *Card) *TransactionState {
	return &TransactionState{card: card}
}

func (s *TransactionState) HandleFlow() {
	fmt.Println("Press 1 to withdraw or 2 to deposit")
	var keyPress int
	fmt.Scanln(&keyPress)
	fmt.Println("Enter amount")
	var amount float64
	fmt.Scanln(&amount)
	accountBalance := s.card.GetAccount().Balance()
	if keyPress == 1 {
		if accountBalance < amount {
			fmt.Println("Insufficient funds!")
		} else {
			s.card.GetAccount().SetBalance(accountBalance - amount)
		}
	} else if keyPress == 2 {
		s.card.GetAccount().SetBalance(accountBalance + amount)
	}
	fmt.Println("Your current balance:", s.card.GetAccount().Balance())
}

type RemoveCardState struct{}

func (s *RemoveCardState) HandleFlow() {
	fmt.Println("Please remove your card")
}

type EndState struct{}

func (s *EndState) HandleFlow() {
	fmt.Println("Thank you!")
}

type AtmContext struct {
	cache        map[StateEnum]StateHandler
	allStates    []StateEnum
	currentState int
	totalStates  int
	state        StateHandler
	card         *Card
}

func NewAtmContext(card *Card) *AtmContext {
	return &AtmContext{
		cache:        make(map[StateEnum]StateHandler),
		allStates:    []StateEnum{WELCOME, INSERT_CARD, TRANSACTION, REMOVE_CARD, THANK_YOU},
		currentState: -1,
		totalStates:  5,
		card:         card,
	}
}

func (c *AtmContext) SetState(state StateEnum, stateHandler StateHandler) {
	c.cache[state] = stateHandler
}
