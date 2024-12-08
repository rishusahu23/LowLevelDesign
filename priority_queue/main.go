package main

import (
	"container/heap"
	"fmt"
)

// Item represents an item in the priority queue.
type Item struct {
	value    string // The value of the item; can be any type.
	priority int    // The priority of the item in the queue.
}

// PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

// Len returns the number of items in the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares priorities between two items (for a min-heap, lower priority comes first).
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

// Swap swaps two items in the queue.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push adds an item to the priority queue.
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the highest priority (lowest number).
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func main() {
	// Initialize an empty priority queue and add some items.
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Adding items to the queue
	heap.Push(&pq, &Item{value: "task 1", priority: 3})
	heap.Push(&pq, &Item{value: "task 2", priority: 1})
	heap.Push(&pq, &Item{value: "task 3", priority: 2})

	// Pop items by priority
	fmt.Println("Items by priority:")
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%s (priority: %d)\n", item.value, item.priority)
	}
}

func mainn() {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &Item{
		value:    "fd",
		priority: 0,
	})

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Println(item.value)
	}
}