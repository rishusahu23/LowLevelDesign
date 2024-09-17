package main

import (
	"fmt"
	"time"
)

// UserType represents different types of users in the system.
type UserType string

const (
	Customer        UserType = "Customer"
	DeliveryPartner UserType = "DeliveryPartner"
	RestaurantOwner UserType = "RestaurantOwner"
)

// User interface defines common behavior for all users.
type User interface {
	GetID() int
	GetName() string
	GetUserType() UserType
}

// UserImpl struct represents a generic user in the system.
type UserImpl struct {
	ID   int
	Name string
	Type UserType
}

func (u *UserImpl) GetID() int {
	return u.ID
}

func (u *UserImpl) GetName() string {
	return u.Name
}

func (u *UserImpl) GetUserType() UserType {
	return u.Type
}

// Restaurant represents a restaurant in the system.
type Restaurant struct {
	ID       int
	Name     string
	Location string
	Menu     []MenuItem
}

// MenuItem represents an item in the restaurant's menu.
type MenuItem struct {
	ID    int
	Name  string
	Price float64
}

// Order represents an order in the system.
type Order struct {
	ID         int
	CustomerID int
	Restaurant *Restaurant
	Items      []MenuItem
	TotalPrice float64
	Status     string
	CreatedAt  time.Time
}

// Delivery represents a delivery in the system.
type Delivery struct {
	ID              int
	OrderID         int
	DeliveryPartner User
	Status          string
	AssignedAt      time.Time
}

// UserFactory creates users based on user type.
type UserFactory struct{}

func (f *UserFactory) CreateUser(id int, name string, userType UserType) User {
	return &UserImpl{
		ID:   id,
		Name: name,
		Type: userType,
	}
}

// RestaurantService handles operations related to restaurants.
type RestaurantService struct {
	restaurants []*Restaurant
}

func (s *RestaurantService) RegisterRestaurant(r *Restaurant) {
	s.restaurants = append(s.restaurants, r)
	fmt.Println("Restaurant registered:", r.Name)
}

func (s *RestaurantService) AddMenuItem(restaurantID int, item MenuItem) {
	for _, r := range s.restaurants {
		if r.ID == restaurantID {
			r.Menu = append(r.Menu, item)
			fmt.Println("Menu item added:", item.Name, "to restaurant:", r.Name)
		}
	}
}

// OrderService handles operations related to orders.
type OrderService struct {
	orders []*Order
}

func (s *OrderService) PlaceOrder(order *Order) {
	order.Status = "Placed"
	order.CreatedAt = time.Now()
	s.orders = append(s.orders, order)
	fmt.Println("Order placed:", order.ID, "for customer:", order.CustomerID)
}

// DeliveryService handles operations related to deliveries.
type DeliveryService struct {
	deliveries []*Delivery
}

func (s *DeliveryService) AssignDelivery(orderID int, partner User) {
	delivery := &Delivery{
		ID:              len(s.deliveries) + 1,
		OrderID:         orderID,
		DeliveryPartner: partner,
		Status:          "Assigned",
		AssignedAt:      time.Now(),
	}
	s.deliveries = append(s.deliveries, delivery)
	fmt.Println("Delivery assigned:", delivery.ID, "to partner:", partner.GetName())
}

// Observer defines an interface for notifying users.
type Observer interface {
	Notify(event string)
}

// CustomerNotifier implements Observer for notifying customers.
type CustomerNotifier struct {
	customerID int
}

func (n *CustomerNotifier) Notify(event string) {
	fmt.Printf("Notifying customer %d: %s\n", n.customerID, event)
}

// NotificationService manages observers and sends notifications.
type NotificationService struct {
	observers map[int]Observer
}

func (s *NotificationService) RegisterObserver(userID int, observer Observer) {
	if s.observers == nil {
		s.observers = make(map[int]Observer)
	}
	s.observers[userID] = observer
}

func (s *NotificationService) NotifyObservers(event string) {
	for _, observer := range s.observers {
		observer.Notify(event)
	}
}

func main() {
	// Initialize services
	userFactory := &UserFactory{}
	restaurantService := &RestaurantService{}
	orderService := &OrderService{}
	deliveryService := &DeliveryService{}
	notificationService := &NotificationService{}

	// Create users
	customer := userFactory.CreateUser(1, "John Doe", Customer)
	partner := userFactory.CreateUser(2, "Jane Smith", DeliveryPartner)
	_ = userFactory.CreateUser(3, "Restaurant Owner", RestaurantOwner)

	// Register a restaurant
	restaurant := &Restaurant{
		ID:       1,
		Name:     "Pizza Place",
		Location: "123 Main St",
	}
	restaurantService.RegisterRestaurant(restaurant)

	// Add menu items
	restaurantService.AddMenuItem(restaurant.ID, MenuItem{ID: 1, Name: "Margherita Pizza", Price: 8.99})
	restaurantService.AddMenuItem(restaurant.ID, MenuItem{ID: 2, Name: "Pepperoni Pizza", Price: 9.99})

	// Place an order
	order := &Order{
		ID:         1,
		CustomerID: customer.GetID(),
		Restaurant: restaurant,
		Items:      restaurant.Menu,
		TotalPrice: 18.98,
	}
	orderService.PlaceOrder(order)

	// Assign delivery
	deliveryService.AssignDelivery(order.ID, partner)

	// Register observers for notifications
	customerNotifier := &CustomerNotifier{customerID: customer.GetID()}
	notificationService.RegisterObserver(customer.GetID(), customerNotifier)

	// Notify customer about order status
	notificationService.NotifyObservers("Your order has been placed and is being prepared.")
}
