package main

import (
	"fmt"
	"sort"
	"sync"
)

// Book represents a book with a title and an author.
type Book struct {
	Title  string
	Author string
}

// User represents a user with an ID, name, and a list of books.
type User struct {
	ID      int
	Name    string
	Books   []Book
	Friends []*User
}

// Database is a singleton that stores users and their relationships.
type Database struct {
	users map[int]*User
	mu    sync.Mutex
}

var instance *Database
var once sync.Once

// GetDatabaseInstance returns the singleton instance of the database.
func GetDatabaseInstance() *Database {
	once.Do(func() {
		instance = &Database{
			users: make(map[int]*User),
		}
	})
	return instance
}

// UserService defines the methods to get books information.
type UserService interface {
	GetUserBooks(userID int) []Book
	GetTopKFriendsBooks(userID int, k int) []Book
	GetTopKNetworkBooks(userID int, k int) []Book
}

// TopKStrategy defines the strategy for selecting the top K books.
type TopKStrategy interface {
	GetTopKBooks(bookCount map[Book]int, k int) []Book
}

// SimpleTopKStrategy implements TopKStrategy using a simple sorting algorithm.
type SimpleTopKStrategy struct{}

// GetTopKBooks returns the top K books based on their count using a simple sorting algorithm.
func (s *SimpleTopKStrategy) GetTopKBooks(bookCount map[Book]int, k int) []Book {
	type bookCountPair struct {
		book  Book
		count int
	}
	var bookCountPairs []bookCountPair
	for book, count := range bookCount {
		bookCountPairs = append(bookCountPairs, bookCountPair{book, count})
	}

	sort.Slice(bookCountPairs, func(i, j int) bool {
		return bookCountPairs[i].count > bookCountPairs[j].count
	})

	var topKBooks []Book
	for i := 0; i < k && i < len(bookCountPairs); i++ {
		topKBooks = append(topKBooks, bookCountPairs[i].book)
	}

	return topKBooks
}

// GoodreadsService implements UserService.
type GoodreadsService struct {
	db       *Database
	topKAlgo TopKStrategy
}

// NewGoodreadsService creates a new GoodreadsService.
func NewGoodreadsService(db *Database, strategy TopKStrategy) *GoodreadsService {
	return &GoodreadsService{db: db, topKAlgo: strategy}
}

// GetUserBooks returns the list of books a user has read.
func (s *GoodreadsService) GetUserBooks(userID int) []Book {
	user := s.db.users[userID]
	if user == nil {
		return nil
	}
	return user.Books
}

// GetTopKFriendsBooks returns the top K books that friends have read.
func (s *GoodreadsService) GetTopKFriendsBooks(userID int, k int) []Book {
	user := s.db.users[userID]
	if user == nil {
		return nil
	}

	bookCount := make(map[Book]int)
	for _, friend := range user.Friends {
		for _, book := range friend.Books {
			bookCount[book]++
		}
	}

	return s.topKAlgo.GetTopKBooks(bookCount, k)
}

// GetTopKNetworkBooks returns the top K books that the network has read.
func (s *GoodreadsService) GetTopKNetworkBooks(userID int, k int) []Book {
	visited := make(map[int]bool)
	bookCount := make(map[Book]int)
	queue := []*User{s.db.users[userID]}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == nil || visited[current.ID] {
			continue
		}
		visited[current.ID] = true
		for _, book := range current.Books {
			bookCount[book]++
		}
		for _, friend := range current.Friends {
			if !visited[friend.ID] {
				queue = append(queue, friend)
			}
		}
	}

	return s.topKAlgo.GetTopKBooks(bookCount, k)
}

func main() {
	db := GetDatabaseInstance()

	// Create users and books
	book1 := Book{Title: "1984", Author: "George Orwell"}
	book2 := Book{Title: "To Kill a Mockingbird", Author: "Harper Lee"}
	book3 := Book{Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"}

	user1 := &User{ID: 1, Name: "Alice", Books: []Book{book1, book2}}
	user2 := &User{ID: 2, Name: "Bob", Books: []Book{book2, book3}}
	user3 := &User{ID: 3, Name: "Charlie", Books: []Book{book1, book3}}
	user4 := &User{ID: 4, Name: "Diana", Books: []Book{book2}}

	// Create friendships
	user1.Friends = []*User{user2, user3}
	user2.Friends = []*User{user1, user4}
	user3.Friends = []*User{user1}
	user4.Friends = []*User{user2}

	// Add users to the database
	db.users[user1.ID] = user1
	db.users[user2.ID] = user2
	db.users[user3.ID] = user3
	db.users[user4.ID] = user4

	topKStrategy := &SimpleTopKStrategy{}
	service := NewGoodreadsService(db, topKStrategy)

	fmt.Println("User 1 Books:", service.GetUserBooks(1))
	fmt.Println("Top 2 Friends Books of User 1:", service.GetTopKFriendsBooks(1, 2))
	fmt.Println("Top 2 Network Books of User 1:", service.GetTopKNetworkBooks(1, 2))
}
