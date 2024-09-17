package main

import (
	"fmt"
	"sync"
)

// User represents a user of the book reader system.
type User struct {
	UserID         string
	Name           string
	Email          string
	ReadingHistory []Book
	Bookmarks      map[string]Bookmark
	mu             sync.Mutex
}

// AddBookmark adds a bookmark for a specific book and page number.
func (u *User) AddBookmark(book Book, pageNumber int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	bookmark := Bookmark{
		BookmarkID: fmt.Sprintf("%s-%d", book.BookID, pageNumber),
		User:       *u,
		Book:       book,
		PageNumber: pageNumber,
	}
	u.Bookmarks[bookmark.BookmarkID] = bookmark
}

// RemoveBookmark removes a bookmark.
func (u *User) RemoveBookmark(book Book, pageNumber int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	bookmarkID := fmt.Sprintf("%s-%d", book.BookID, pageNumber)
	delete(u.Bookmarks, bookmarkID)
}

// ViewBook simulates viewing a book.
func (u *User) ViewBook(book Book) {
	fmt.Printf("%s is reading %s by %s\n", u.Name, book.Title, book.Author)
}

// SearchBooks searches for books by title.
func (u *User) SearchBooks(library Library, query string) []Book {
	return library.SearchBooks(query)
}

// GetReadingHistory returns reading history.
func (u *User) GetReadingHistory() []Book {
	return u.ReadingHistory
}

// Book represents a book in the system.
type Book struct {
	BookID     string
	Title      string
	Author     string
	Content    []Page
	TotalPages int
}

// GetPage returns the content of a specific page.
func (b *Book) GetPage(pageNumber int) Page {
	return b.Content[pageNumber-1]
}

// Page represents a page in a book.
type Page struct {
	PageNumber int
	Text       string
	Notes      []string
	Highlights []string
	mu         sync.Mutex
}

// AddHighlight adds a highlight to a page.
func (p *Page) AddHighlight(text string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Highlights = append(p.Highlights, text)
}

// AddNote adds a note to a page.
func (p *Page) AddNote(note string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Notes = append(p.Notes, note)
}

// RemoveHighlight removes a highlight from a page.
func (p *Page) RemoveHighlight(text string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, h := range p.Highlights {
		if h == text {
			p.Highlights = append(p.Highlights[:i], p.Highlights[i+1:]...)
			break
		}
	}
}

// RemoveNote removes a note from a page.
func (p *Page) RemoveNote() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.Notes) > 0 {
		p.Notes = p.Notes[:len(p.Notes)-1]
	}
}

// Bookmark represents a bookmark added by a user.
type Bookmark struct {
	BookmarkID string
	User       User
	Book       Book
	PageNumber int
}

// Library represents a collection of books.
type Library struct {
	Books map[string]Book
	mu    sync.RWMutex
}

// SearchBooks searches for books by title.
func (l *Library) SearchBooks(query string) []Book {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var results []Book
	for _, book := range l.Books {
		if contains(book.Title, query) {
			results = append(results, book)
		}
	}
	return results
}

// GetBookByID retrieves a book by its ID.
func (l *Library) GetBookByID(bookID string) (Book, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	book, found := l.Books[bookID]
	return book, found
}

// contains checks if the query string is contained in the target string.
func contains(target, query string) bool {
	return true // Simplified for demonstration. You can use strings.Contains or other methods for actual implementation.
}

// BookSource represents a source from where books are fetched.
type BookSource interface {
	SearchBooks(query string) []Book
	FetchBookByID(bookID string) (Book, bool)
}

// Aggregator aggregates multiple book sources.
type Aggregator struct {
	Sources []BookSource
}

// SearchAcrossSources searches for books across all sources.
func (a *Aggregator) SearchAcrossSources(query string) []Book {
	var results []Book
	for _, source := range a.Sources {
		results = append(results, source.SearchBooks(query)...)
	}
	return results
}

// FetchBookFromSource fetches a book from a specific source.
func (a *Aggregator) FetchBookFromSource(sourceName, bookID string) (Book, bool) {
	for _, source := range a.Sources {
		if book, found := source.FetchBookByID(bookID); found {
			return book, true
		}
	}
	return Book{}, false
}

// Main function to demonstrate the use of the classes.
func main() {
	// Initialize library and books
	library := Library{Books: make(map[string]Book)}
	book1 := Book{
		BookID:     "1",
		Title:      "The Great Gatsby",
		Author:     "F. Scott Fitzgerald",
		Content:    []Page{{PageNumber: 1, Text: "Page 1 Content"}, {PageNumber: 2, Text: "Page 2 Content"}},
		TotalPages: 2,
	}
	library.Books[book1.BookID] = book1

	// Initialize user
	user := User{
		UserID:         "user1",
		Name:           "John Doe",
		Email:          "john.doe@example.com",
		Bookmarks:      make(map[string]Bookmark),
		ReadingHistory: []Book{book1},
	}

	// User actions
	user.ViewBook(book1)
	user.AddBookmark(book1, 1)
	fmt.Println("Bookmarks:", user.Bookmarks)
	user.RemoveBookmark(book1, 1)
	fmt.Println("Bookmarks after removal:", user.Bookmarks)
}
