package main

import (
	"errors"
	"fmt"
)

// User represents a user in the system
type User struct {
	ID     string
	Name   string
	Amount float64
}

// Group represents a group of users
type Group struct {
	ID    string
	Name  string
	Users []User
}

// Split defines the amount that a user owes for an expense
type Split struct {
	User   User
	Amount float64
}

// Expense represents an expense
type Expense struct {
	Amount float64
	PaidBy User
	Splits []Split
}

// ExpenseSplitter is an interface for various splitting strategies
type ExpenseSplitter interface {
	SplitExpense(expense Expense) ([]Split, error)
}

// EqualSplit splits the expense equally among users
type EqualSplit struct{}

func (e *EqualSplit) SplitExpense(expense Expense) ([]Split, error) {
	totalUsers := len(expense.Splits)
	if totalUsers == 0 {
		return nil, errors.New("no users to split the expense")
	}

	splitAmount := expense.Amount / float64(totalUsers)
	for i := range expense.Splits {
		expense.Splits[i].Amount = splitAmount
	}
	return expense.Splits, nil
}

// PercentageSplit splits the expense based on percentage
type PercentageSplit struct {
	Percentages []float64
}

func (p *PercentageSplit) SplitExpense(expense Expense) ([]Split, error) {
	if len(p.Percentages) != len(expense.Splits) {
		return nil, errors.New("percentages count mismatch with users")
	}

	for i, percentage := range p.Percentages {
		expense.Splits[i].Amount = expense.Amount * percentage / 100
	}
	return expense.Splits, nil
}

// SplitFactory provides the appropriate splitting strategy
type SplitFactory struct{}

func (sf *SplitFactory) GetSplitter(strategy string) ExpenseSplitter {
	switch strategy {
	case "equal":
		return &EqualSplit{}
	case "percentage":
		return &PercentageSplit{}
	default:
		return nil
	}
}

// GroupService manages the group and its users
type GroupService struct {
	Groups []Group
}

// AddGroup adds a new group to the service
func (gs *GroupService) AddGroup(name string, users []User) Group {
	group := Group{
		ID:    fmt.Sprintf("group-%d", len(gs.Groups)+1),
		Name:  name,
		Users: users,
	}
	gs.Groups = append(gs.Groups, group)
	return group
}

// GetGroup retrieves a group by its ID
func (gs *GroupService) GetGroup(id string) (*Group, error) {
	for _, group := range gs.Groups {
		if group.ID == id {
			return &group, nil
		}
	}
	return nil, errors.New("group not found")
}

// GroupExpenseService handles expense management within a group
type GroupExpenseService struct {
	Expenses map[string][]Expense // map of groupID to list of expenses
}

// AddExpense adds an expense to the group and splits it according to the strategy
func (ges *GroupExpenseService) AddExpense(group Group, amount float64, paidBy User, splitter ExpenseSplitter) error {
	// Prepare splits
	splits := make([]Split, len(group.Users))
	for i, user := range group.Users {
		splits[i] = Split{User: user}
	}

	expense := Expense{
		Amount: amount,
		PaidBy: paidBy,
		Splits: splits,
	}

	// Use the splitter strategy to calculate split amounts
	splitResult, err := splitter.SplitExpense(expense)
	if err != nil {
		return err
	}
	expense.Splits = splitResult

	// Add expense to group's expense list
	if ges.Expenses == nil {
		ges.Expenses = make(map[string][]Expense)
	}
	ges.Expenses[group.ID] = append(ges.Expenses[group.ID], expense)

	// Update balances
	for i, split := range splitResult {
		group.Users[i].Amount -= split.Amount
		if split.User.ID == paidBy.ID {
			group.Users[i].Amount += expense.Amount
		}
	}
	return nil
}

// Main Function
func main() {
	// Example users
	user1 := User{ID: "1", Name: "Alice", Amount: 0}
	user2 := User{ID: "2", Name: "Bob", Amount: 0}
	user3 := User{ID: "3", Name: "Charlie", Amount: 0}

	// Group service and adding users to a group
	groupService := &GroupService{}
	group := groupService.AddGroup("Friends", []User{user1, user2, user3})

	// Create an expense service for the group
	expenseService := &GroupExpenseService{}

	// Add an equal split expense
	splitFactory := &SplitFactory{}
	equalSplitter := splitFactory.GetSplitter("equal")

	err := expenseService.AddExpense(group, 120, user1, equalSplitter)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Display updated balances for the group
	for _, user := range group.Users {
		fmt.Printf("%s owes %.2f\n", user.Name, user.Amount)
	}

	// Adding another expense with percentage split
	percentageSplitter := &PercentageSplit{
		Percentages: []float64{50, 30, 20},
	}

	err = expenseService.AddExpense(group, 200, user2, percentageSplitter)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Display updated balances for the group
	fmt.Println("Updated Balances After Percentage Split:")
	for _, user := range group.Users {
		fmt.Printf("%s owes %.2f\n", user.Name, user.Amount)
	}
}
