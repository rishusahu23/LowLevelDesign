package main

import (
	"errors"
	"fmt"
)

var (
	_ State = &IdleState{}
	_ State = &SelectionState{}
	_ State = &DispenseState{}
	_ State = &PaymentState{}
)

type IdleState struct {
}

func (i *IdleState) SelectProduct(ctx *Context, productId string) error {
	product, err := ctx.inventory.GetProduct(productId)
	if err != nil {
		return err
	}
	if product.Stock == 0 {
		return errors.New("product out of stock")
	}
	ctx.selectedProduct = product
	ctx.SetState(&SelectionState{})
	return nil
}

func (i *IdleState) InsertMoney(ctx *Context, amount float64) error {
	return errors.New("no product selected")
}

func (i *IdleState) DispenseProduct(ctx *Context) error {
	return errors.New("no product selected")
}

func (i *IdleState) Cancel(ctx *Context) error {
	return errors.New("no transaction to cancel")
}

func (i *IdleState) Refund(ctx *Context) error {
	return errors.New("no transaction to refund")
}

type SelectionState struct {
}

func (s *SelectionState) SelectProduct(ctx *Context, productId string) error {
	return fmt.Errorf("product already selected")
}

func (s *SelectionState) InsertMoney(ctx *Context, amount float64) error {
	ctx.balance += amount
	if ctx.balance >= ctx.selectedProduct.Price {
		ctx.SetState(&DispenseState{})
	}
	return nil
}

func (s *SelectionState) DispenseProduct(ctx *Context) error {
	return fmt.Errorf("payment not completed")
}

func (s *SelectionState) Cancel(ctx *Context) error {
	ctx.selectedProduct = nil
	ctx.balance = 0
	ctx.SetState(&IdleState{})
	return nil
}

func (s *SelectionState) Refund(ctx *Context) error {
	if ctx.balance == 0 {
		return errors.New("no money inserted to refund")
	}
	refundAmount := ctx.balance
	ctx.balance = 0
	ctx.SetState(&IdleState{})
	fmt.Println("refunded amount: ", refundAmount)
	return nil
}

type PaymentState struct {
}

func (p *PaymentState) SelectProduct(ctx *Context, productId string) error {
	return fmt.Errorf("payment in progress")
}

func (p *PaymentState) InsertMoney(ctx *Context, amount float64) error {
	ctx.balance += amount
	if ctx.balance >= ctx.selectedProduct.Price {
		ctx.SetState(&DispenseState{})
	}
	return nil
}

func (p *PaymentState) DispenseProduct(ctx *Context) error {
	//TODO implement me
	panic("implement me")
}

func (p *PaymentState) Cancel(ctx *Context) error {
	ctx.selectedProduct = nil
	ctx.balance = 0
	ctx.SetState(&IdleState{})
	return nil
}

func (p *PaymentState) Refund(ctx *Context) error {
	//TODO implement me
	panic("implement me")
}

type DispenseState struct {
}

func (d *DispenseState) SelectProduct(ctx *Context, productId string) error {
	return fmt.Errorf("dispensing in progress")
}

func (d *DispenseState) InsertMoney(ctx *Context, amount float64) error {
	return fmt.Errorf("dispensing in progress")
}

func (d *DispenseState) DispenseProduct(ctx *Context) error {
	product := ctx.selectedProduct
	product.Stock -= 1
	ctx.balance -= product.Price
	ctx.selectedProduct = nil
	ctx.SetState(&IdleState{})
	if ctx.balance > 0 {
		ctx.SetState(&SelectionState{})
	}
	ctx.inventory.AddProduct(product)
	return nil
}

func (d *DispenseState) Cancel(ctx *Context) error {
	return fmt.Errorf("cannot cancel during dispensing")
}

func (d *DispenseState) Refund(ctx *Context) error {
	return fmt.Errorf("cannot refund during dispensing")
}
