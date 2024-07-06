package book

type Book struct {
	Id              string
	Title           string
	Author          string
	Category        string
	PublicationDate string
	RackNumber      string
	Copies          []*Item
}

type Item struct {
	Book          *Book
	Barcode       string
	IsCheckedOut  bool
	DueDate       string
	CurrentMember string
}
