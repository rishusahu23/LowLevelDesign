package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type VehicleType int

const (
	Car VehicleType = iota
	Bike
	Truck
)

type Vehicle interface {
	GetNumberPlate() string
	GetType() VehicleType
	GetEntryTime() time.Time
	SetEntryTime(time.Time)
}

type VehicleImpl struct {
	numberPlate string
	entryTime   time.Time
	vType       VehicleType
}

func NewVehicleImpl(numberPlate string, vType VehicleType) *VehicleImpl {
	return &VehicleImpl{numberPlate: numberPlate, vType: vType}
}

func (v *VehicleImpl) GetNumberPlate() string {
	return v.numberPlate
}

func (v *VehicleImpl) GetType() VehicleType {
	return v.vType
}

func (v *VehicleImpl) GetEntryTime() time.Time {
	return v.entryTime
}

func (v *VehicleImpl) SetEntryTime(t time.Time) {
	v.entryTime = t
}

type ParkingSpace interface {
	GetId() string
	GetLevel() int
	GetType() VehicleType
	IsOccupied() bool
	OccupySpace()
	FreeSpace()
}

type ParkingSpaceImpl struct {
	id       string
	level    int
	vType    VehicleType
	occupied bool
}

func (p *ParkingSpaceImpl) GetId() string {
	return p.id
}

func (p *ParkingSpaceImpl) GetLevel() int {
	return p.level
}

func (p *ParkingSpaceImpl) GetType() VehicleType {
	return p.vType
}

func (p *ParkingSpaceImpl) IsOccupied() bool {
	return p.occupied
}

func (p *ParkingSpaceImpl) OccupySpace() {
	p.occupied = true
}

func (p *ParkingSpaceImpl) FreeSpace() {
	p.occupied = false
}

type PricingStrategy interface {
	CalculateCharges(vehicle Vehicle) float64
}

type HourlyPricing struct {
	rates map[VehicleType]float64
}

func (h *HourlyPricing) CalculateCharges(vehicle Vehicle) float64 {
	duration := time.Now().Sub(vehicle.GetEntryTime()).Hours()
	return duration * h.rates[vehicle.GetType()]
}

type ParkingLot interface {
	Enter(vehicle Vehicle) (string, int, error)
	Exit(vehicle Vehicle, spaceId string) (float64, error)
	DisplayAvailability()
}

type ParkingLotImpl struct {
	levels       map[int][]ParkingSpace
	mutex        sync.Mutex
	pricing      PricingStrategy
	spacePerType map[VehicleType]int
}

func (p *ParkingLotImpl) Enter(vehicle Vehicle) (string, int, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for level, spaces := range p.levels {
		for _, space := range spaces {
			if space.GetType() == vehicle.GetType() && !space.IsOccupied() {
				space.OccupySpace()
				vehicle.SetEntryTime(time.Now())
				p.spacePerType[vehicle.GetType()]--
				return space.GetId(), level, nil
			}
		}
	}
	return "", 0, errors.New("no available space")
}

func (p *ParkingLotImpl) Exit(vehicle Vehicle, spaceId string) (float64, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, spaces := range p.levels {
		for _, space := range spaces {
			if space.GetId() == spaceId && space.IsOccupied() && space.GetType() == vehicle.GetType() {
				space.FreeSpace()
				p.spacePerType[vehicle.GetType()]++
				return p.pricing.CalculateCharges(vehicle), nil
			}
		}
	}
	return 0, errors.New("invalid space id or vehicle")
}

func (p *ParkingLotImpl) DisplayAvailability() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	fmt.Println("Available Parking Spaces:")
	for vType, count := range p.spacePerType {
		fmt.Printf("%d: %d\n", vType, count)
	}
}

var instance *ParkingLotImpl
var once sync.Once

func GetParkingLot(rates map[VehicleType]float64, levels, spacePerLevel int) {
	once.Do(func() {
		instance = &ParkingLotImpl{
			levels: make(map[int][]ParkingSpace),
			pricing: &HourlyPricing{
				rates: rates,
			},
			spacePerType: make(map[VehicleType]int),
		}
		id := 1
		for i := 1; i <= levels; i++ {
			for j := 0; j < spacePerLevel; j++ {
				var vType VehicleType
				switch j % 3 {
				case 0:
					vType = Bike
				case 1:
					vType = Car
				case 2:
					vType = Truck
				}
				instance.levels[i] = append(instance.levels[i], &ParkingSpaceImpl{
					id:    fmt.Sprintf("%v", id),
					level: i,
					vType: vType,
				})
				instance.spacePerType[vType]++
				id++
			}
		}
	})
}
