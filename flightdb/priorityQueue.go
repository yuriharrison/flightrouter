package flightsdb

// QueueItem item struct
type QueueItem struct {
	data  *Airport
	value float32
	index int
}

// Queue queue to be used with container/heap
type Queue []*QueueItem

func (q Queue) Len() int {
	return len(q)
}

func (q Queue) Less(i, j int) bool {
	return q[i].value < q[j].value
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

// Push use this method with container/heap
func (q *Queue) Push(x interface{}) {
	n := len(*q)
	item := x.(*QueueItem)
	item.index = n
	*q = append(*q, item)
}

// Pop use this method with container/heap
func (q *Queue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*q = old[0 : n-1]
	return item
}
