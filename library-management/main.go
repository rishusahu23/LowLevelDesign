package main

import "github.com/rishu/design/library-management/system"

func main() {
	sys := system.GetLibrarySystem()
	memberId := "member1"
	bookItemId := "bookItem1"

	sys.CheckoutBook(memberId, bookItemId)
	sys.ShowNotification(memberId)
}
