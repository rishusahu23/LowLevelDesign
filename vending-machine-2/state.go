package main

type State interface {
	SelectProduct(ctx *Context, productId string) error
	InsertMoney(ctx *Context, amount float64) error
	DispenseProduct(ctx *Context) error
	Cancel(ctx *Context) error
	Refund(ctx *Context) error
}
