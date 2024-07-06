package main

import (
	"errors"
	"fmt"
	"time"
)

type IParkingLot interface {
	AddSpace(space IParkingSpace) error
	RemoveSpace(spaceId string) error
	FindSpace() (IParkingSpace, error)
	ParkVehicle(vehicle IVehicle) (ITicket, error)
	ExitVehicle(ticketId string) error
}

type IParkingSpace interface {
	GetId() string
	IsAvailable() bool
	Allocate(vehicle IVehicle) error
	Vacate() error
}

type IVehicle interface {
	GetId() string
	GetType() string
}

type ITicket interface {
	GetId() string
	GetIssueTime() time.Time
	GetVehicleId() string
	GetSpaceId() string
}

type IPaymentProcessor interface {
	CalculateFee(ticket ITicket) float64
	ProcessPayment(ticket ITicket, amount float64) error
}

type IGate interface {
	Open()
	Exit()
}

type INotificationService interface {
	NotifySpaceAvailable(spaceID string)
	NotifyPaymentDue(ticketID string)
}

type ParkingLot struct {
	Spaces          map[string]IParkingSpace
	Tickets         map[string]ITicket
	paymentProc     IPaymentProcessor
	notificationSvc INotificationService
}

func (p *ParkingLot) AddSpace(space IParkingSpace) error {
	p.Spaces[space.GetId()] = space
	return nil
}

func (p *ParkingLot) RemoveSpace(spaceId string) error {
	delete(p.Spaces, spaceId)
	return nil
}

func (p *ParkingLot) FindSpace() (IParkingSpace, error) {
	for _, space := range p.Spaces {
		if space.IsAvailable() {
			return space, nil
		}
	}
	return nil, errors.New("not found")
}

func (p *ParkingLot) ParkVehicle(vehicle IVehicle) (ITicket, error) {
	space, _ := p.FindSpace()
	err := space.Allocate(vehicle)
	if err != nil {
		return nil, err
	}
	ticket := NewTicket(vehicle.GetId(), space.GetId())
	p.Tickets[ticket.GetId()] = ticket
	return ticket, nil
}

func (p *ParkingLot) ExitVehicle(ticketId string) error {
	ticket, exists := p.Tickets[ticketId]
	if !exists {
		return errors.New("ticket not found")
	}
	space := p.Spaces[ticket.GetSpaceId()]
	err := space.Vacate()
	if err != nil {
		return err
	}
	fee := p.paymentProc.CalculateFee(ticket)
	err = p.paymentProc.ProcessPayment(ticket, fee)
	if err != nil {
		return err
	}
	delete(p.Tickets, ticketId)
	// p.notificationSvc.NotifySpaceAvailable(ticket.GetSpaceId())
	return nil
}

func NewParkingLot(paymentProc IPaymentProcessor, notificationSvc INotificationService) *ParkingLot {
	return &ParkingLot{
		Spaces:          make(map[string]IParkingSpace),
		Tickets:         make(map[string]ITicket),
		paymentProc:     paymentProc,
		notificationSvc: notificationSvc,
	}
}

type ParkingSpace struct {
	Id        string
	Available bool
	Vehicle   IVehicle
}

func NewParkingSpace(id string) *ParkingSpace {
	return &ParkingSpace{
		Id:        id,
		Available: true,
	}
}

func (p *ParkingSpace) GetId() string {
	return p.Id
}

func (p *ParkingSpace) IsAvailable() bool {
	return p.Available
}

func (p *ParkingSpace) Allocate(vehicle IVehicle) error {
	if !p.Available {
		return errors.New("space not available")
	}
	p.Vehicle = vehicle
	p.Available = false
	return nil
}

func (p *ParkingSpace) Vacate() error {
	if p.Available {
		return errors.New("space is already vacant")
	}
	p.Vehicle = nil
	p.Available = true
	return nil
}

var _ ITicket = &Ticket{}

type Ticket struct {
	id        string
	issueTime time.Time
	vehicleID string
	spaceID   string
}

func NewTicket(vehicleID string, spaceID string) *Ticket {
	return &Ticket{
		id:        fmt.Sprintf("%s-%s", vehicleID, time.Now()),
		issueTime: time.Now(),
		vehicleID: vehicleID,
		spaceID:   spaceID,
	}
}

func (t *Ticket) GetId() string {
	return t.id
}

func (t *Ticket) GetIssueTime() time.Time {
	return t.issueTime
}

func (t *Ticket) GetVehicleId() string {
	return t.vehicleID
}

func (t *Ticket) GetSpaceId() string {
	return t.spaceID
}

type PaymentProcessor struct{}

func (pp *PaymentProcessor) CalculateFee(ticket ITicket) float64 {
	duration := time.Since(ticket.GetIssueTime())
	hours := duration.Hours()
	return hours * 10.0 // Assume $10 per hour
}

func (pp *PaymentProcessor) ProcessPayment(ticket ITicket, amount float64) error {
	fmt.Printf("Processed payment of $%.2f for ticket %s\n", amount, ticket.GetId())
	return nil
}

type Vehicle struct {
	id    string
	vType string
}

func NewVehicle(id, vType string) *Vehicle {
	return &Vehicle{id: id, vType: vType}
}

func (v *Vehicle) GetId() string {
	return v.id
}

func (v *Vehicle) GetType() string {
	return v.vType
}

func main() {
	paymentProcessor := &PaymentProcessor{}
	parkingLot := NewParkingLot(paymentProcessor, nil)
	parkingLot.AddSpace(NewParkingSpace("1"))
	parkingLot.AddSpace(NewParkingSpace("2"))

	car := NewVehicle("CAR123", "car")
	motorcycle := NewVehicle("MOTO123", "motorcycle")

	ticket1, _ := parkingLot.ParkVehicle(car)
	ticket2, _ := parkingLot.ParkVehicle(motorcycle)

	parkingLot.ExitVehicle(ticket2.GetId())
	parkingLot.ExitVehicle(ticket1.GetId())
}
