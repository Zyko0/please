package queue

import "container/heap"

type Item struct {
	index    int
	priority uint64
	object   any
}

type queue []*Item

func (pq queue) Len() int {
	return len(pq)
}

func (pq queue) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq queue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *queue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *queue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

type TickQueue[T any] struct {
	pq *queue
}

func NewTickQueue[T any]() *TickQueue[T] {
	tq := &TickQueue[T]{
		pq: &queue{},
	}
	return tq
}

func (tq *TickQueue[T]) Push(object *T, priority uint64) {
	heap.Push(tq.pq, &Item{
		index:    -1,
		priority: priority,
		object:   object,
	})
}

func (tq *TickQueue[T]) Pop() *T {
	if len(*tq.pq) == 0 {
		return nil
	}
	return heap.Pop(tq.pq).(*Item).object.(*T)
}

func (tq *TickQueue[T]) Len() int {
	return len(*tq.pq)
}
