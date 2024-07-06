package main

import (
	"fmt"
	"time"
)

type Direction int

const (
	UP Direction = iota
	DOWN
)

func (d Direction) String() string {
	switch d {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	default:
		return "UNKNOWN"
	}
}

// DirectionBehavior interface
type DirectionBehavior interface {
	ProcessRequest(e *Elevator, request Request)
	AddPendingJobsToCurrentJobs(e *Elevator)
}

// UpDirection behavior
type UpDirection struct{}

func (d *UpDirection) ProcessRequest(e *Elevator, request Request) {
	fmt.Println("Processing up request")

	if e.currentFloor < request.externalRequest.sourceFloor {
		for i := e.currentFloor; i <= request.externalRequest.sourceFloor; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("We have reached floor --", i)
			e.currentFloor = i
		}
	}

	fmt.Println("Reached Source Floor--opening door")

	for i := e.currentFloor + 1; i <= request.internalRequest.destinationFloor; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("We have reached floor --", i)
		e.currentFloor = i
		if e.checkIfNewJobCanBeProcessed(request) {
			break
		}
	}
}

func (d *UpDirection) AddPendingJobsToCurrentJobs(e *Elevator) {
	if len(e.upPendingJobs) > 0 {
		fmt.Println("Pick a pending up job and execute it by putting in current jobs")
		e.currentJobs = append(e.currentJobs, e.upPendingJobs...)
		e.upPendingJobs = nil
	} else {
		e.currentState = &IdleState{}
		fmt.Println("The elevator is in Idle state")
	}
}

// DownDirection behavior
type DownDirection struct{}

func (d *DownDirection) ProcessRequest(e *Elevator, request Request) {
	fmt.Println("Processing down request")

	if e.currentFloor < request.externalRequest.sourceFloor {
		for i := e.currentFloor; i <= request.externalRequest.sourceFloor; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("We have reached floor --", i)
			e.currentFloor = i
		}
	}

	fmt.Println("Reached Source Floor--opening door")

	for i := e.currentFloor - 1; i >= request.internalRequest.destinationFloor; i-- {
		time.Sleep(1 * time.Second)
		fmt.Println("We have reached floor --", i)
		e.currentFloor = i
		if e.checkIfNewJobCanBeProcessed(request) {
			break
		}
	}
}

func (d *DownDirection) AddPendingJobsToCurrentJobs(e *Elevator) {
	if len(e.downPendingJobs) > 0 {
		fmt.Println("Pick a pending down job and execute it by putting in current jobs")
		e.currentJobs = append(e.currentJobs, e.downPendingJobs...)
		e.downPendingJobs = nil
	} else {
		e.currentState = &IdleState{}
		fmt.Println("The elevator is in Idle state")
	}
}

// ElevatorState interface
type ElevatorState interface {
	StartElevator(e *Elevator)
	AddJob(e *Elevator, request Request)
}

// IdleState
type IdleState struct{}

func (s *IdleState) StartElevator(e *Elevator) {
	fmt.Println("The Elevator is idle.")
}

func (s *IdleState) AddJob(e *Elevator, request Request) {
	e.currentState = &MovingState{}
	e.currentDirection = getDirectionBehavior(request.externalRequest.directionToGo)
	e.currentJobs = append(e.currentJobs, request)
	e.currentState.StartElevator(e)
}

// MovingState
type MovingState struct{}

func (s *MovingState) StartElevator(e *Elevator) {
	fmt.Println("The Elevator has started functioning")
	for {
		if e.checkIfJob() {
			e.currentDirection.ProcessRequest(e, e.currentJobs[0])
			e.currentJobs = e.currentJobs[1:]
			if len(e.currentJobs) == 0 {
				e.currentDirection.AddPendingJobsToCurrentJobs(e)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (s *MovingState) AddJob(e *Elevator, request Request) {
	switch e.currentDirection.(type) {
	case UpDirection:
		if request.externalRequest.directionToGo != Up {
			e.addToPendingJobs(request)
		} else {
			if request.internalRequest.destinationFloor < e.currentFloor {
				e.addToPendingJobs(request)
			} else {
				e.currentJobs = append(e.currentJobs, request)
			}
		}
	case DownDirection:
		if request.externalRequest.directionToGo != Down {
			e.addToPendingJobs(request)
		} else {
			if request.internalRequest.destinationFloor > e.currentFloor {
				e.addToPendingJobs(request)
			} else {
				e.currentJobs = append(e.currentJobs, request)
			}
		}
	default:
		// Handle error, unknown direction
	}
}



// Elevator struct
type Elevator struct {
	currentState     ElevatorState
	currentDirection DirectionBehavior
	currentFloor     int

	currentJobs     []Request
	upPendingJobs   []Request
	downPendingJobs []Request
}

func NewElevator() *Elevator {
	return &Elevator{
		currentState: &IdleState{},
		currentFloor: 0,
	}
}

// Request structs
type ExternalRequest struct {
	directionToGo Direction
	sourceFloor   int
}

type InternalRequest struct {
	destinationFloor int
}

type Request struct {
	internalRequest InternalRequest
	externalRequest ExternalRequest
}

// Elevator methods
func (e *Elevator) checkIfJob() bool {
	return len(e.currentJobs) > 0
}

func (e *Elevator) checkIfNewJobCanBeProcessed(request Request) bool {
	if e.checkIfJob() {
		if e.currentDirection == &UpDirection{} {
			lastRequest := e.currentJobs[len(e.currentJobs)-1]
			if lastRequest.internalRequest.destinationFloor < request.internalRequest.destinationFloor {
				e.currentJobs = append(e.currentJobs, request)
				return true
			}
		} else if e.currentDirection == &DownDirection{} {
			firstRequest := e.currentJobs[0]
			if firstRequest.internalRequest.destinationFloor > request.internalRequest.destinationFloor {
				e.currentJobs = append(e.currentJobs, request)
				return true
			}
		}
	}
	return false
}

func (e *Elevator) addToPendingJobs(request Request) {
	if request.externalRequest.directionToGo == UP {
		fmt.Println("Add to pending up jobs")
		e.upPendingJobs = append(e.upPendingJobs, request)
	} else {
		fmt.Println("Add to pending down jobs")
		e.downPendingJobs = append(e.downPendingJobs, request)
	}
}

func (e *Elevator) addJob(request Request) {
	e.currentState.AddJob(e, request)
}

func getDirectionBehavior(direction Direction) DirectionBehavior {
	if direction == UP {
		return &UpDirection{}
	} else {
		return &DownDirection{}
	}
}

func main() {
	elevator := NewElevator()

	go elevator.currentState.StartElevator(elevator)

	time.Sleep(3 * time.Second)

	er := ExternalRequest{directionToGo: UP, sourceFloor: 0}
	ir := InternalRequest{destinationFloor: 5}
	request1 := Request{internalRequest: ir, externalRequest: er}

	go func() {
		time.Sleep(200 * time.Millisecond)
		elevator.addJob(request1)
	}()

	time.Sleep(10 * time.Second) // Let the elevator run for a bit before the program exits
}
