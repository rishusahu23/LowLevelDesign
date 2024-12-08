package main

import "fmt"

type Book interface {
	GetID() int
	GetTitle() string
	GetAuthor() string
}

type PhysicalBook struct {
	ID     int
	Title  string
	Author string
}

func (b PhysicalBook) GetID() int {
	return b.ID
}

func (b PhysicalBook) GetTitle() string {
	return b.Title
}

func (b PhysicalBook) GetAuthor() string {
	return b.Author
}

type User interface {
	GetID() int
	GetName() string
	GetBooks() []Book
	AddBook(b Book)
	AddFriend(friend User)
	GetFriends() []User
}

type RegularUser struct {
	ID      int
	Name    string
	Books   []Book
	Friends []User
}

func (u *RegularUser) GetID() int {
	return u.ID
}

// GetName returns the name of the user.
func (u *RegularUser) GetName() string {
	return u.Name
}

// GetBooks returns the list of books owned by the user.
func (u *RegularUser) GetBooks() []Book {
	return u.Books
}

// AddBook adds a book to the user's collection.
func (u *RegularUser) AddBook(book Book) {
	u.Books = append(u.Books, book)
}

// AddFriend adds a friend to the user's friend list.
func (u *RegularUser) AddFriend(friend User) {
	u.Friends = append(u.Friends, friend)
}

// GetFriends returns the list of friends.
func (u *RegularUser) GetFriends() []User {
	return u.Friends
}

type UserManager interface {
	AddBook(userId int, book Book)
	GetBooks(userId int) []Book
	AddFriend(userId int, friend User)
	GetFriends(userId int) []User
}

type BookManager interface {
	GetTopBooks(userID int) []Book
}

type FriendBookService interface {
	GetTopBooksFromFriends(userID int) []Book
	GetTopBooksFromNetwork(userID int) []Book
}

type ConcreteUserManager struct {
	users map[int]User
}

func NewUserManager() *ConcreteUserManager {
	return &ConcreteUserManager{users: make(map[int]User)}
}

func (um *ConcreteUserManager) AddBook(userID int, book Book) {
	user, exists := um.users[userID]
	if !exists {
		fmt.Println("User not found")
		return
	}
	user.AddBook(book)
}

func (um *ConcreteUserManager) GetBooks(userID int) []Book {
	user, exists := um.users[userID]
	if !exists {
		fmt.Println("User not found")
		return nil
	}
	return user.GetBooks()
}

func (um *ConcreteUserManager) AddFriend(userID int, friend User) {
	user, exists := um.users[userID]
	if !exists {
		fmt.Println("User not found")
		return
	}
	user.AddFriend(friend)
	um.users[friend.GetID()] = friend
}

func (um *ConcreteUserManager) GetFriends(userID int) []User {
	user, exists := um.users[userID]
	if !exists {
		fmt.Println("User not found")
		return nil
	}
	return user.GetFriends()
}

type ConcreteFriendBookService struct {
	userManager *ConcreteUserManager
}

func NewFriendBookService(um *ConcreteUserManager) *ConcreteFriendBookService {
	return &ConcreteFriendBookService{userManager: um}
}

// GetTopBooksFromFriends returns top books that friends of the user have read.
func (fbs *ConcreteFriendBookService) GetTopBooksFromFriends(userID int) []Book {
	user := fbs.userManager.users[userID]
	bookCount := make(map[string]int)

	for _, friend := range user.GetFriends() {
		for _, book := range friend.GetBooks() {
			bookCount[book.GetTitle()]++
		}
	}

	return getTopBooksByCount(bookCount)
}

func getTopBooksByCount(count map[string]int) []Book {
	return nil
}
