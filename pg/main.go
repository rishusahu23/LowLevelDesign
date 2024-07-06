package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

type PaymentHandler interface {
	ProcessPayment(details string) error
}

type NetBankingDetails struct {
	Username string
	Password string
}

type CardDetails struct {
	CardNumber string
	Expiry     string
	CVV        string
}

type UPIDetails struct {
	VPA string
}

type NetBankingHandler struct {
}

type CardHandler struct {
}

type UPIHandler struct {
}

func randomResponse() error {
	if rand.Float32() < 0.9 {
		return nil
	}
	return errors.New("payment failed")
}

func (n *NetBankingHandler) ProcessPayment(details string) error {
	var nbd NetBankingDetails
	err := json.Unmarshal([]byte(details), &nbd)
	if err != nil {
		return fmt.Errorf("invalid net banking details: %w", err)
	}
	fmt.Println("processing net banking payment for user: ", nbd.Username)
	return randomResponse()
}

func (n *CardHandler) ProcessPayment(details string) error {
	var cardDetails CardDetails
	err := json.Unmarshal([]byte(details), &cardDetails)
	if err != nil {
		return fmt.Errorf("invalid card details: %w", err)
	}
	fmt.Println("Processing Card payment for card number:", cardDetails.CardNumber)
	return randomResponse()
}

func (n *UPIHandler) ProcessPayment(details string) error {
	var upiDetails UPIDetails
	err := json.Unmarshal([]byte(details), &upiDetails)
	if err != nil {
		return fmt.Errorf("invalid UPI details: %w", err)
	}
	fmt.Println("Processing UPI payment for VPA:", upiDetails.VPA)
	return randomResponse()
}

type PaymentRequest struct {
	Mode     string
	ClientId string
	Details  string
}

type Bank struct {
	Name         string
	Handlers     map[string]PaymentHandler
	Traffic      int
	TrafficCount int
}

type Client struct {
	Id    string
	Modes map[string]bool
}

type PaymentProcessor struct {
	Clients      map[string]*Client
	Banks        map[string]*Bank
	BankPriority map[string][]string
}

func NewPaymentProcessor() *PaymentProcessor {
	return &PaymentProcessor{
		Clients: make(map[string]*Client),
		Banks: map[string]*Bank{
			"HDFC": {
				Name:         "HDFC",
				Handlers:     map[string]PaymentHandler{"netbanking": &NetBankingHandler{}, "card": &CardHandler{}, "upi": &UPIHandler{}},
				Traffic:      0,
				TrafficCount: 0,
			},
			"ICICI": {
				Name:         "HDFC",
				Handlers:     map[string]PaymentHandler{"netbanking": &NetBankingHandler{}, "card": &CardHandler{}, "upi": &UPIHandler{}},
				Traffic:      0,
				TrafficCount: 0,
			},
			"SBI": {
				Name:         "HDFC",
				Handlers:     map[string]PaymentHandler{"netbanking": &NetBankingHandler{}, "card": &CardHandler{}, "upi": &UPIHandler{}},
				Traffic:      0,
				TrafficCount: 0,
			},
		},
		BankPriority: make(map[string][]string),
	}
}

func (p *PaymentProcessor) AddClient(client string) {
	p.Clients[client] = &Client{
		Id:    client,
		Modes: make(map[string]bool),
	}
}

func (p *PaymentProcessor) RemoveClient(id string) {
	delete(p.Clients, id)
}

func (p *PaymentProcessor) HasClient(id string) bool {
	_, ok := p.Clients[id]
	return ok
}

func (p *PaymentProcessor) ListSupportedModes(id string) []string {
	var modes []string
	if client, ok := p.Clients[id]; ok {
		for mode := range client.Modes {
			modes = append(modes, mode)
		}
	} else {
		for mode := range p.Banks[id].Handlers {
			modes = append(modes, mode)
		}
	}
	return modes
}

func (p *PaymentProcessor) AddSupportedModes(mode, id string) {
	if client, ok := p.Clients[id]; ok {
		client.Modes[mode] = true
	}
}

func (p *PaymentProcessor) RemovePaymode(mode, clientID string) {
	if client, exists := p.Clients[clientID]; exists {
		delete(client.Modes, mode)
	}
}

func (p *PaymentProcessor) ShowDistribution() map[string]int {
	return nil
}

func (p *PaymentProcessor) MakePayment(req *PaymentRequest) error {
	client, ok := p.Clients[req.ClientId]
	if !ok {
		return fmt.Errorf("client %s does not exist", req.ClientId)
	}
	if _, ok := client.Modes[req.Mode]; !ok {
		return fmt.Errorf("payment mode %s not supported for client %s", req.Mode, req.ClientId)
	}
	bankName, err := p.routeToBank(req.Mode)
	if err != nil {
		return err
	}
	bank, bankExists := p.Banks[bankName]
	if !bankExists {
		return fmt.Errorf("bank %s not found", bankName)
	}
	handler, modeSupported := bank.Handlers[req.Mode]
	if !modeSupported {
		return fmt.Errorf("payment mode %s not supported by bank %s", req.Mode, bankName)
	}
	err = handler.ProcessPayment(req.Details)
	if err == nil {
		bank.TrafficCount++
	}
	return err
}

func (p *PaymentProcessor) routeToBank(mode string) (string, error) {
	banks := p.BankPriority[mode]
	if len(banks) == 0 {
		return "", fmt.Errorf("no banks configured for mode %s", mode)
	}
	bankName := banks[0]
	banks = append(banks[1:], bankName)
	p.BankPriority[mode] = banks
	return bankName, nil
}
