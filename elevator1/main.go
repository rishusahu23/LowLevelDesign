package main

import "fmt"

type Elevator interface {
	MoveToFloor(floor int)
	OpenDoor()
	CloseDoor()
	GetCurrentFloor() int
	IsMoving() bool
}

type ElevatorImpl struct {
	CurrentFloor int
	isMoving     bool
	Door         Door
}

func (e *ElevatorImpl) MoveToFloor(floor int) {
	if e.Door.IsOpen() {
		e.Door.Close()
	}
	e.isMoving = true
	fmt.Printf("moving from floor %v, to %v", e.CurrentFloor, floor)
	e.CurrentFloor = floor
	e.isMoving = false
	e.Door.Open()
}

func (e *ElevatorImpl) OpenDoor() {
	e.Door.Open()
}

func (e *ElevatorImpl) CloseDoor() {
	e.Door.Close()
}

func (e *ElevatorImpl) GetCurrentFloor() int {
	return e.CurrentFloor
}

func (e *ElevatorImpl) IsMoving() bool {
	return e.isMoving
}

func NewElevatorImpl() *ElevatorImpl {
	return &ElevatorImpl{
		CurrentFloor: 0,
		isMoving:     false,
		Door:         NewDoorImpl(),
	}
}

type Door interface {
	Open()
	Close()
	IsOpen() bool
}

type DoorImpl struct {
	open bool
}

func NewDoorImpl() Door {
	return &DoorImpl{
		open: false,
	}
}

func (d *DoorImpl) Open() {
	d.open = true
}

func (d *DoorImpl) Close() {
	d.open = false
}

func (d *DoorImpl) IsOpen() bool {
	return d.open
}

type Request interface {
	GetFloor() int
	IsInternal() bool
}

type RequestImpl struct {
	Floor    int
	Internal bool
}

func (r *RequestImpl) GetFloor() int {
	return r.Floor
}

func (r *RequestImpl) IsInternal() bool {
	return r.Internal
}

func NewRequestImpl(floor int, internal bool) *RequestImpl {
	return &RequestImpl{Floor: floor, Internal: internal}
}

type Scheduler interface {
	AddRequest(request Request)
	GetNextFloor(currentFloor int) int
}

type SchedulerImpl struct {
	Requests []Request
}

func (s *SchedulerImpl) AddRequest(request Request) {
	s.Requests = append(s.Requests, request)
}

func (s *SchedulerImpl) GetNextFloor(currentFloor int) int {
	if len(s.Requests) == 0 {
		return -1
	}
	nextFloor := s.Requests[0]
	s.Requests = s.Requests[1:]
	return nextFloor.GetFloor()
}

func NewSchedulerImpl() *SchedulerImpl {
	return &SchedulerImpl{Requests: make([]Request, 0)}
}

type ControlSystem interface {
	RequestElevator(floor int, internal bool)
	Run()
}

type ControlSystemImpl struct {
	elevator  Elevator
	scheduler Scheduler
}

func (c *ControlSystemImpl) RequestElevator(floor int, internal bool) {
	req := NewRequestImpl(floor, internal)
	c.scheduler.AddRequest(req)
}

func (c *ControlSystemImpl) Run() {
	for {
		if !c.elevator.IsMoving() {
			nextFloor := c.scheduler.GetNextFloor(c.elevator.GetCurrentFloor())
			if nextFloor != -1 {
				c.elevator.MoveToFloor(nextFloor)
			} else {
				fmt.Println("no more requests")
				break
			}
		}
	}
}

func NewControlSystemImpl() *ControlSystemImpl {
	return &ControlSystemImpl{
		elevator:  NewElevatorImpl(),
		scheduler: NewSchedulerImpl(),
	}
}

func main() {
	cs := NewControlSystemImpl()

	// Simulate external requests
	cs.RequestElevator(3, false)
	cs.RequestElevator(1, false)

	// Simulate internal requests
	cs.RequestElevator(5, true)
	cs.RequestElevator(2, true)

	// Run the control system
	cs.Run()
}
