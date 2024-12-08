package main

import (
	"fmt"
	"sort"
)

type Command interface {
	Execute()
}

type Elevator struct {
	CurrentFloor int
	TotalFloors  int
	Moving       bool
}

func NewElevator(total, current int) *Elevator {
	return &Elevator{
		CurrentFloor: current,
		TotalFloors:  total,
		Moving:       false,
	}
}

type MoveUpCommand struct {
	ele *Elevator
}

func (m *MoveUpCommand) Execute() {
	if m.ele.CurrentFloor < m.ele.TotalFloors {
		m.ele.CurrentFloor++
		fmt.Printf("Elevator moving up to floor %d\n", m.ele.CurrentFloor)
		m.ele.Moving = true
	}
}

type MoveDownCommand struct {
	ele *Elevator
}

func (m *MoveDownCommand) Execute() {
	if m.ele.CurrentFloor > 1 {
		m.ele.CurrentFloor--
		fmt.Printf("elevator moving down to %d\n", m.ele.CurrentFloor)
		m.ele.Moving = true
	}
}

type StopAndOpenDoorCommand struct {
	ele *Elevator
}

func (s *StopAndOpenDoorCommand) Execute() {
	s.ele.Moving = false
	fmt.Printf("Elevator stopped at floor %d\n", s.ele.CurrentFloor)
}

type CloseDoorCommand struct {
}

func (c *CloseDoorCommand) Execute() {
	fmt.Println("Doors are closing")
}

type CommandProcessor struct {
	ele       *Elevator
	upQueue   []int
	downQueue []int
	movingUp  bool
}

func NewCommandProcessor(ele *Elevator) *CommandProcessor {
	return &CommandProcessor{
		ele:       ele,
		upQueue:   make([]int, 0),
		downQueue: make([]int, 0),
		movingUp:  true,
	}
}

func (cp *CommandProcessor) AddRequest(floor int) {
	if floor > cp.ele.CurrentFloor {
		cp.upQueue = append(cp.upQueue, floor)
	} else {
		cp.downQueue = append(cp.downQueue, floor)
	}
	cp.upQueue = deduplicateAndSort(cp.upQueue, true)
	cp.downQueue = deduplicateAndSort(cp.downQueue, false)
}

func deduplicateAndSort(queue []int, ascending bool) []int {
	set := make(map[int]struct{})
	for _, floor := range queue {
		set[floor] = struct{}{}
	}
	result := make([]int, 0)
	for floor := range set {
		result = append(result, floor)
	}
	if ascending {
		sort.Ints(result)
	} else {
		sort.Sort(sort.Reverse(sort.IntSlice(result)))
	}
	return result
}

func (cp *CommandProcessor) ProcessCommands() {
	for len(cp.upQueue) > 0 || len(cp.downQueue) > 0 {
		if cp.movingUp {
			if len(cp.upQueue) > 0 {
				cp.ExecuteCommandsSequence(&cp.upQueue, true)
			} else {
				cp.movingUp = false
			}
		} else {
			if len(cp.downQueue) > 0 {
				cp.ExecuteCommandsSequence(&cp.downQueue, false)
			} else {
				cp.movingUp = true
			}
		}
	}
}

func (cp *CommandProcessor) ExecuteCommandsSequence(queue *[]int, ascending bool) {
	nextFloor := (*queue)[0]
	*queue = (*queue)[1:]

	for cp.ele.CurrentFloor != nextFloor {
		var cmd Command
		if ascending {
			cmd = &MoveUpCommand{ele: cp.ele}
		} else {
			cmd = &MoveDownCommand{ele: cp.ele}
		}
		cmd.Execute()
	}
	var cmd Command
	cmd = &StopAndOpenDoorCommand{ele: cp.ele}
	cmd.Execute()
	cmd = &CloseDoorCommand{}
	cmd.Execute()
}

func main() {
	ele := NewElevator(10, 2)
	commandProcessor := NewCommandProcessor(ele)
	commandProcessor.AddRequest(5)
	commandProcessor.AddRequest(8)
	commandProcessor.AddRequest(1)
	commandProcessor.AddRequest(3)
	commandProcessor.ProcessCommands()
}
