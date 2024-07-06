package search

import "github.com/rishu/design/library-management/book"

type Search interface {
	SearchByTitle(title string) []*book.Book
	SearchByAuthor(author string) []*book.Book
	SearchByCategory(category string) []*book.Book
	SearchByPublicationDate(date string) []*book.Book
}
