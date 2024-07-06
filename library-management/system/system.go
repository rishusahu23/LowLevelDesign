package system

import (
	"fmt"
	"github.com/rishu/design/library-management/book"
	"github.com/rishu/design/library-management/member"
	"github.com/rishu/design/library-management/notification"
	"github.com/rishu/design/library-management/search"
	"sync"
)

type LibrarySystem struct {
	Books         []*book.Book
	Members       []*member.Member
	Notifications []*notification.Notification
	Search        search.Search
	mu            sync.Mutex
}

var instance *LibrarySystem
var once sync.Once

func GetLibrarySystem() *LibrarySystem {
	once.Do(func() {
		instance = &LibrarySystem{}
	})
	return instance
}

func (ls *LibrarySystem) CheckoutBook(memberId, bookItemId string) {

}

func (ls *LibrarySystem) ReserveBook(memberId, bookItemId string) {

}

func (ls *LibrarySystem) ReturnBook(memberId, bookItemId string) {

}

func (ls *LibrarySystem) AddNotification(memberId, message string) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.Notifications = append(ls.Notifications, &notification.Notification{
		MemberId: memberId,
		Message:  message,
	})
}

func (ls *LibrarySystem) ShowNotification(memberId string) {
	for _, n := range ls.Notifications {
		if n.MemberId == memberId {
			fmt.Println(n.Message)
		}
	}
}
