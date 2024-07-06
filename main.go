package main

import (
	"container/heap"
	"fmt"
	"time"
)

type Order struct {
	ID       string
	Time     time.Time
	Stock    string
	Side     string
	Price    float64
	Quantity int
	Priority int
}

type OrderPriorityQueue []*Order

func (pq OrderPriorityQueue) Len() int { return len(pq) }

func (pq OrderPriorityQueue) Less(i, j int) bool {
	if pq[i].Price != pq[j].Price {
		if pq[i].Side == "buy" {
			return pq[i].Price > pq[j].Price
		}
		return pq[i].Price < pq[j].Price
	}
	return pq[i].Time.Before(pq[j].Time)
}

func (pq OrderPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Priority = i
	pq[j].Priority = j
}

func (pq *OrderPriorityQueue) Push(x interface{}) {
	n := len(*pq)
	order := x.(*Order)
	order.Priority = n
	*pq = append(*pq, order)
}

func (pq *OrderPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	order := old[n-1]
	order.Priority = -1
	old[n-1] = nil
	*pq = old[0 : n-1]
	return order
}

type OrderBook struct {
	buyOrders  OrderPriorityQueue
	sellOrders OrderPriorityQueue
}

func (ob *OrderBook) addOrder(order *Order) {
	if order.Side == "buy" {
		heap.Push(&ob.buyOrders, order)
	} else {
		heap.Push(&ob.sellOrders, order)
	}
	ob.matchOrders()
}

func (ob *OrderBook) matchOrders() {
	for ob.buyOrders.Len() > 0 && ob.sellOrders.Len() > 0 {
		buyOrder := ob.buyOrders[0]
		sellOrder := ob.sellOrders[0]

		if buyOrder.Price > sellOrder.Price {
			tradeQuantity := min(buyOrder.Quantity, sellOrder.Quantity)
			tradePrice := sellOrder.Price

			fmt.Printf("%s %f %d %s\n", buyOrder.ID, tradePrice, tradeQuantity, sellOrder.ID)

			buyOrder.Quantity -= tradeQuantity
			sellOrder.Quantity -= tradeQuantity

			if buyOrder.Quantity == 0 {
				heap.Pop(&ob.buyOrders)
			}
			if sellOrder.Quantity == 0 {
				heap.Pop(&ob.sellOrders)
			}
		} else {
			break
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Exchange maintains order books for multiple stocks
type Exchange struct {
	orderBooks map[string]*OrderBook
}


func NewExchange() *Exchange {
	return &Exchange{
		orderBooks: make(map[string]*OrderBook),
	}
}

func (e *Exchange) addOrder(order *Order) {
	if _, exists := e.orderBooks[order.Stock]; !exists {
		e.orderBooks[order.Stock] = &OrderBook{
			buyOrders:  make(OrderPriorityQueue, 0),
			sellOrders: make(OrderPriorityQueue, 0),
		}
		heap.Init(&e.orderBooks[order.Stock].buyOrders)
		heap.Init(&e.orderBooks[order.Stock].sellOrders)
	}
	e.orderBooks[order.Stock].addOrder(order)
}