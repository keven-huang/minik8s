package queue

import "sync"

type ConcurrentQueue struct {
	mutex sync.Mutex
	queue []interface{}
}

func NewConcurrentQueue() *ConcurrentQueue {
	return &ConcurrentQueue{
		queue: make([]interface{}, 0),
	}
}

func (q *ConcurrentQueue) Push(item interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.queue = append(q.queue, item)
}

func (q *ConcurrentQueue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.queue) == 0 {
		return nil
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item
}

func (q *ConcurrentQueue) IsEmpty() bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.queue) == 0
}
