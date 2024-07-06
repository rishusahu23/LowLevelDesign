package main

import (
	"errors"
	"fmt"
	"time"
)

type IParkingLot interface {
	AddSpace(space *ParkingSpace) error
	RemoveSpace(spaceId string) error
	FindSpace() (*ParkingSpace, error)
	ParkVehicle(vehicle *Vehicle) (*Ticket, error)
	ExitVehicle(ticketId string) error
}

type IParkingSpace interface {
	GetId() string
	IsAvailable() bool
	Allocate(vehicle *Vehicle) error
	Vacate() error
}

type IVehicle interface {
	GetId() string
	GetType() string
}

type ITicket interface {
	GetId() string
	GetIssueTime() time.Time
	GetSpaceId() string
	GetVehicleId() string
}

type IPayment interface {
	CalculateFee()
	ProcessPayment()
}

type INotification interface {
	NotifySpaceAvailable()
	NotifyPaymentDue()
}

type Ticket struct {
	Id        string
	SpaceId   string
	VehicleId string
	IssueTime time.Time
}

func NewTicket(spaceId, vehicleId string) *Ticket {
	return &Ticket{
		Id:        fmt.Sprintf("%v-%v", vehicleId, time.Now().UnixNano()),
		SpaceId:   spaceId,
		VehicleId: vehicleId,
		IssueTime: time.Now(),
	}
}

func (t *Ticket) GetId() string {
	return t.Id
}

func (t *Ticket) GetIssueTime() time.Time {
	return t.IssueTime
}

func (t *Ticket) GetVehicleId() string {
	return t.VehicleId
}

func (t *Ticket) GetSpaceId() string {
	return t.SpaceId
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

type ParkingSpace struct {
	Id        string
	Available bool
	Vehicle   *Vehicle
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

func (p *ParkingSpace) Allocate(vehicle *Vehicle) error {
	if !p.IsAvailable() {
		return errors.New("not available")
	}
	p.Vehicle = vehicle
	p.Available = false
	return nil
}

func (p *ParkingSpace) Vacate() error {
	if p.IsAvailable() {
		return errors.New("space is empty")
	}
	p.Vehicle = nil
	p.Available = false
	return nil
}

type ParkingLot struct {
	Spaces  map[string]*ParkingSpace
	Tickets map[string]*Ticket
}

func NewParkingLot() *ParkingLot {
	return &ParkingLot{
		Spaces:  make(map[string]*ParkingSpace),
		Tickets: make(map[string]*Ticket),
	}
}

func (p *ParkingLot) AddSpace(space *ParkingSpace) error {
	p.Spaces[space.GetId()] = space
	return nil
}

func (p *ParkingLot) RemoveSpace(spaceId string) error {
	delete(p.Spaces, spaceId)
	return nil
}

func (p *ParkingLot) FindSpace() (*ParkingSpace, error) {
	for _, sp := range p.Spaces {
		if sp.IsAvailable() {
			return sp, nil
		}
	}
	return nil, errors.New("not found")
}

func (p *ParkingLot) ParkVehicle(vehicle *Vehicle) (*Ticket, error) {
	sp, _ := p.FindSpace()
	err := sp.Allocate(vehicle)
	if err != nil {
		return nil, err
	}
	ticket := NewTicket(sp.GetId(), vehicle.GetId())
	p.Tickets[ticket.GetId()] = ticket
	return ticket, nil
}

func (p *ParkingLot) ExitVehicle(ticketId string) error {
	//
	return nil
}
