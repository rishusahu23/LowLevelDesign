package main

import (
	"fmt"
	"time"
)

type State interface {
	TransitionLight(light *TrafficLight)
}

type RedState struct {
}

type YellowState struct {
}

func (y YellowState) TransitionLight(light *TrafficLight) {
	fmt.Println("yellow light is on")
	light.changeState(&GreenState{})
	light.startTime(2 * time.Second)
}

type GreenState struct {
}

func (g GreenState) TransitionLight(light *TrafficLight) {
	fmt.Println("Green light is on")
	light.changeState(&RedState{})
	light.startTime(2 * time.Second)
}

var _ State = &RedState{}

func (r RedState) TransitionLight(light *TrafficLight) {
	fmt.Println("red light is on")
	light.changeState(&YellowState{})
	light.startTime(2 * time.Second)
}

type TrafficLight struct {
	state State
	timer *time.Timer
}

func (t *TrafficLight) changeState(state State) {
	if t.timer != nil {
		t.timer.Stop()
	}
	t.state = state
	fmt.Println("Traffic light state:", t.state)
}

func (t *TrafficLight) startTime(duration time.Duration) {
	t.timer = time.NewTimer(duration)
	go func() {
		<-t.timer.C
		t.state.TransitionLight(t)
	}()
}

func main() {
	tr := &TrafficLight{
		state: &RedState{},
		timer: nil,
	}
	tr.startTime(2 * time.Second)
}
