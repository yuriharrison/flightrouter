package flightsdb

import (
	"container/heap"
)

// FindCheapestRoute implements Dijkstra Algorithm to find
// the cheapest route for a given origin and destination
func FindCheapestRoute(origin, destination *Airport) []*Flight {
	routeTrace := make(map[*Airport]*Flight)
	visited := make(map[*Airport]bool)
	accPriceTable := make(map[*Airport]float32)
	queue := &Queue{&QueueItem{data: origin, value: 0}}
	accPriceTable[origin] = 0
	heap.Init(queue)

	for queue.Len() > 0 {
		if _, ok := visited[destination]; ok {
			break
		}

		item := heap.Pop(queue).(*QueueItem)
		ap := item.data
		for _, flight := range ap.flights {
			dest := flight.dest
			accPrice := accPriceTable[ap] + flight.Price
			if oldAccPrice, ok := accPriceTable[dest]; !ok || accPrice < oldAccPrice {
				accPriceTable[dest] = accPrice
				routeTrace[dest] = flight
				heap.Push(queue, &QueueItem{data: dest, value: accPrice})
			}
		}
		visited[ap] = true
	}
	if _, ok := visited[destination]; !ok {
		return nil
	}
	return reverseRoute(destination, routeTrace)
}

func reverseRoute(destination *Airport, routeTrace map[*Airport]*Flight) []*Flight {
	route := []*Flight{}
	for next := routeTrace[destination]; next != nil; next = routeTrace[next.orig] {
		route = append(route, next)
	}
	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}
	return route
}
