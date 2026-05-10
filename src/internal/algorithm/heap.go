package algorithm

import "container/heap"

type pqItem struct {
	node  *Node
	index int
	order int
}

type innerHeap []*pqItem

func (h innerHeap) Len() int {
	return len(h)
}

func (h innerHeap) Less(i, j int) bool {
	if h[i].node.F == h[j].node.F {
		return h[i].order < h[j].order
	}
	return h[i].node.F < h[j].node.F
}

func (h innerHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *innerHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*pqItem)
	item.index = n
	*h = append(*h, item)
}

func (h *innerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return item
}

type PriorityQueue struct {
	h     innerHeap
	order int
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(&pq.h)
	return pq
}

func (p *PriorityQueue) Push(n *Node) {
	heap.Push(&p.h, &pqItem{node: n, order: p.order})
	p.order++
}

func (p *PriorityQueue) Pop() *Node {
	return heap.Pop(&p.h).(*pqItem).node
}

func (p *PriorityQueue) Len() int {
	return p.h.Len()
}
