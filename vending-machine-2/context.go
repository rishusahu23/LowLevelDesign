package main

type Context struct {
	state           State
	inventory       *Inventory
	balance         float64
	selectedProduct *Product
}

func (ctx *Context) SetState(state State) {
	ctx.state = state
}

func (ctx *Context) SelectProduct(id string) error {
	return ctx.state.SelectProduct(ctx, id)
}

func (ctx *Context) InsertMoney(amount float64) error {
	return ctx.state.InsertMoney(ctx, amount)
}

func (ctx *Context) DispenseProduct() error {
	return ctx.state.DispenseProduct(ctx)
}

func (ctx *Context) Cancel() error {
	return ctx.state.Cancel(ctx)
}

func (ctx *Context) Refund() error {
	return ctx.state.Refund(ctx)
}
