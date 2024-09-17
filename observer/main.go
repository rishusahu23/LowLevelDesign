package main

import "fmt"

type Observer interface {
	Update(temp float64)
}

type Subject interface {
	RegisterObserver(observer Observer)
	RemoveObserver(observer Observer)
	NotifyAll()
}

type WeatherStation struct {
	observers []Observer
	temp      float64
}

func (w *WeatherStation) RegisterObserver(observer Observer) {
	w.observers = append(w.observers, observer)
}

func (w *WeatherStation) RemoveObserver(observer Observer) {
	res := make([]Observer, 0)
	for _, obs := range w.observers {
		if obs == observer {
			continue
		}
		res = append(res, obs)
	}
	w.observers = res
}

func (w *WeatherStation) NotifyAll() {

}

func (w *WeatherStation) ChangeTemperature(temp float64) {
	fmt.Println("temperature changes")
	for _, obs := range w.observers {
		obs.Update(temp)
	}
}

type Android struct {
}

func (a *Android) Update(temp float64) {
	fmt.Println("temp in android update to, ", temp)
}

type IOS struct {
}

func (a *IOS) Update(temp float64) {
	fmt.Println("temp in IOS update to, ", temp)
}

func main() {
	and := &Android{}
	ios := &IOS{}
	ws := &WeatherStation{
		observers: make([]Observer, 0),
	}
	ws.RegisterObserver(and)
	ws.RegisterObserver(ios)
	ws.ChangeTemperature(34)
}
