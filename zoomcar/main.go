package main

import (
	"fmt"
	"time"
)

type Booking struct {
	StartDate time.Time
	EndDate   time.Time
}

func (b *Booking) Overlaps(start, end time.Time) bool {
	return (start.Before(b.EndDate) && end.After(b.StartDate)) ||
		(start.Equal(b.StartDate) || end.Equal(b.EndDate))
}

type Vehicle struct {
	id       string
	model    string
	price    float64
	location string
	bookings []Booking // List of bookings for this vehicle
}

func (v *Vehicle) IsAvailableForDates(start, end time.Time) bool {
	for _, b := range v.bookings {
		if b.Overlaps(start, end) {
			return false // If any booking overlaps, the vehicle is not available
		}
	}
	return true // No overlapping bookings, so the vehicle is available
}

func (v *Vehicle) BookForDates(start, end time.Time) error {
	if !v.IsAvailableForDates(start, end) {
		return fmt.Errorf("vehicle not available for the selected dates")
	}
	v.bookings = append(v.bookings, Booking{StartDate: start, EndDate: end})
	return nil
}

type Receipt struct {
	VehicleID  string    `json:"vehicle_id"`
	Model      string    `json:"model"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
}

type ReceiptGenerator interface {
	GenerateReceipt(receipt Receipt) string
}

type TextReceipt struct{}

// GenerateReceipt Method for TextReceipt
func (tr *TextReceipt) GenerateReceipt(receipt Receipt) string {
	return fmt.Sprintf("Receipt:\nVehicle ID: %s\nModel: %s\nStart Date: %s\nEnd Date: %s\nTotal Price: %.2f",
		receipt.VehicleID, receipt.Model, receipt.StartDate.Format("2006-01-02"),
		receipt.EndDate.Format("2006-01-02"), receipt.TotalPrice)
}

type BookingManager struct {
	ReceiptStrategy ReceiptGenerator // Use strategy pattern to select receipt format
}

// Method to book a vehicle and generate a receipt
func (bm *BookingManager) BookVehicle(v *Vehicle, start, end time.Time) (string, error) {
	if err := v.BookForDates(start, end); err != nil {
		return "", err
	}
	// Calculate total price for the booking (e.g., price per day * number of days)
	days := end.Sub(start).Hours() / 24
	totalPrice := v.price * days

	// Create a receipt
	receipt := Receipt{
		VehicleID:  v.id,
		Model:      v.model,
		StartDate:  start,
		EndDate:    end,
		TotalPrice: totalPrice,
	}

	// Generate receipt using the selected strategy
	return bm.ReceiptStrategy.GenerateReceipt(receipt), nil
}
